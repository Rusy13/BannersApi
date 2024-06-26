package handlers

import (
	errs "Avito/internal/config/errors"
	initial "Avito/internal/infrastructure/kafka/initialization"
	"Avito/internal/service"
	"Avito/internal/storage/repository"
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Server1 struct {
	Serv service.BannerServ
}

type BannerResponse struct {
	BannerID  int64       `json:"banner_id"`
	TagIDs    []int64     `json:"tag_ids"`
	FeatureID int64       `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (s *Server1) GetUserBanner(w http.ResponseWriter, req *http.Request) {
	tagID, err := strconv.ParseInt(req.URL.Query().Get("tag_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	featureID, err := strconv.ParseInt(req.URL.Query().Get("feature_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid feature ID", http.StatusBadRequest)
		return
	}

	useLastRevision, err := strconv.ParseBool(req.URL.Query().Get("use_last_revision"))
	if err != nil {
		useLastRevision = false
	}
	var banner *repository.Banner
	cache := memcache.New("localhost:11211")
	cacheKey := "banner_" + strconv.Itoa(int(tagID)) + "_" + strconv.Itoa(int(featureID))
	if useLastRevision == false {
		item, err := cache.Get(cacheKey)
		if err == nil {
			err := json.Unmarshal(item.Value, &banner)
			log.Println("Используем данные из кеша")
			if err != nil {
				log.Println("Ошибка Unmarshal кэша")
			}
		}
	}

	banner, err = s.Serv.GetUserBanner(req.Context(), tagID, featureID, useLastRevision)
	if err != nil {
		log.Println("Error fetching user banner:", err) // Добавленная строка
		if errors.Is(err, errs.ErrBannerNotFound) {
			http.Error(w, "Баннер для не найден", http.StatusNotFound)
		} else {
			errorResponse := struct {
				Error string `json:"error"`
			}{
				Error: "Internal Server Error",
			}
			responseBody, err := json.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(responseBody)
		}
		return
	}

	jsonBytes, err := json.Marshal(banner)
	if err != nil {
		log.Println("Ошибка записи в кэш")
	}
	cache.Set(&memcache.Item{Key: cacheKey, Value: jsonBytes, Expiration: 5 * 60}) // Время жизни кеша 5 минут
	log.Println("Баннер сохранен в кэш:", cacheKey)

	if !banner.IsActive {
		adminToken := req.Header.Get("token")
		if adminToken != "admin_token" {
			http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banner)
}

func (s *Server1) GetBanners(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	featureID, _ := strconv.ParseInt(req.URL.Query().Get("feature_id"), 10, 64)
	tagID, _ := strconv.ParseInt(req.URL.Query().Get("tag_id"), 10, 64)
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))

	banners, err := s.Serv.GetBanners(req.Context(), featureID, tagID, limit, offset)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

func (s *Server1) CreateBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var banner repository.Banner
	if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	bannerID, err := s.Serv.CreateBanner(req.Context(), &banner)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"banner_id": bannerID})
}

// ---------------------------
func (s *Server1) UpdateBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)
	bannerID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var banner repository.Banner
	if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Получение текущей версии баннера
	currentBanner, err := s.Serv.GetBanner(req.Context(), bannerID)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера: Failed to get current banner", http.StatusInternalServerError)
		return
	}

	// Проверка, что количество версий баннера не превышает максимальное значение
	versions, err := s.Serv.GetBannerVersionsCount(req.Context(), bannerID)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера: Failed to get banner versions count", http.StatusInternalServerError)
		return
	}
	if versions >= 3 {
		// Удаление самой старой версии баннера
		err := s.Serv.DeleteOldestBannerVersion(req.Context(), bannerID)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера: Failed to delete oldest banner version", http.StatusInternalServerError)
			return
		}
	}

	// Запись предыдущей версии баннера в таблицу banner_versions
	err = s.Serv.CreateBannerVersion(req.Context(), currentBanner)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера: Failed to create banner version", http.StatusInternalServerError)
		return
	}

	// Обновление баннера
	err = s.Serv.UpdateBanner(req.Context(), bannerID, &banner)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера: Failed to update banner", http.StatusInternalServerError)
		return
	}

	// Обновление связанных тегов
	err = s.Serv.UpdateFeatureTags(req.Context(), bannerID, banner.FeatureID, banner.TagIDs)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера: Failed to update feature tags", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server1) DeleteBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)
	bannerID, _ := strconv.ParseInt(vars["id"], 10, 64)

	err := s.Serv.DeleteBanner(req.Context(), bannerID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBannerNotFound):
			http.Error(w, "Баннер для не найден", http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server1) DeleteByFeatureTagIDHandler(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	var brokers = []string{
		"127.0.0.1:9091",
		"127.0.0.1:9092",
	}
	initial.ProducerExample(brokers, req.URL.Path, req.Method)
}

func (s *Server1) GetVersionHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	banners, err := s.Serv.GetVersionHandler(req.Context(), id)
	if err != nil {
		log.Println("Err", err)
		switch {
		case errors.Is(err, errs.ErrBannerNotFound):
			http.Error(w, "Баннер не найден", http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	// Отправка баннеров в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

func (s *Server1) ApplyVersionHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	version, _ := strconv.ParseInt(vars["version_number"], 10, 64)

	err := s.Serv.ApplyVersionHandler(req.Context(), id, version)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBannerNotFound):
			http.Error(w, "Баннер не найден", http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
