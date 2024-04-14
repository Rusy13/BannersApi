package handlers

import (
	errs "Avito/internal/config/errors"
	initial "Avito/internal/infrastructure/kafka/initialization"
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
	Repo repository.BannerRepo
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
	// Извлечение параметров запроса
	tagID, err := strconv.ParseInt(req.URL.Query().Get("tag_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}
	log.Println(tagID)

	featureID, err := strconv.ParseInt(req.URL.Query().Get("feature_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid feature ID", http.StatusBadRequest)
		return
	}
	log.Println(featureID)

	useLastRevision, err := strconv.ParseBool(req.URL.Query().Get("use_last_revision"))
	log.Println(useLastRevision)
	if err != nil {
		useLastRevision = false
	}
	//-----------------------------------------------------------------------------------------
	var banner *repository.Banner
	cache := memcache.New("localhost:11211")
	//проверка кэша
	// Получаем все ключи из кеша------------------------------------
	// Получаем все элементы из кеша

	//конец проверки-------------------------------------------
	cacheKey := "banner_" + strconv.Itoa(int(tagID)) + "_" + strconv.Itoa(int(featureID))
	if useLastRevision == false {
		item, err := cache.Get(cacheKey)
		log.Println("ERR", err)
		if err == nil {
			// Используем данные из кеша
			err := json.Unmarshal(item.Value, &banner)
			log.Println("Используем данные из кеша")
			if err != nil {
				log.Println("Ошибка Unmarshal кэша")
			}
		}
	}

	// Получение баннера пользователя из репозитория
	banner, err = s.Repo.GetUserBanner(req.Context(), tagID, featureID, useLastRevision)
	if err != nil {
		log.Println("Error fetching user banner:", err) // Добавленная строка
		if errors.Is(err, errs.ErrBannerNotFound) {
			http.Error(w, "Banner not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Сохраняем информацию в кеше
	jsonBytes, err := json.Marshal(banner)
	if err != nil {
		log.Println("Ошибка записи в кэш")
	}
	cache.Set(&memcache.Item{Key: cacheKey, Value: jsonBytes, Expiration: 5 * 60}) // Время жизни кеша 5 минут
	log.Println("Баннер сохранен в кэш:", cacheKey)

	// Проверка активности баннера
	if !banner.IsActive {
		// Если баннер неактивен, проверяем наличие токена администратора
		adminToken := req.Header.Get("token")
		if adminToken != "admin_token" {
			http.Error(w, "Banner is not active", http.StatusNotFound)
			return
		}
	}

	// Отправка баннера в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banner)
}

// SSSSSSSSSSSSSSSSSSSSSSSSSSSSIIIIIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUUUUUUUUUUUUUUUUUUUU
func (s *Server1) GetBanners(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parsing request parameters
	featureID, _ := strconv.ParseInt(req.URL.Query().Get("feature_id"), 10, 64)
	tagID, _ := strconv.ParseInt(req.URL.Query().Get("tag_id"), 10, 64)
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))

	// Retrieving banners from the repository
	banners, err := s.Repo.GetBanners(req.Context(), featureID, tagID, limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Returning a successful response with banners in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// SSSSSSSSSSSSSSSSSSSSSSSSSSSSIIIIIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUUUUUUUUUUUUUUUUUUUU
func (s *Server1) CreateBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var banner repository.Banner
	log.Println(banner)
	if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Println(banner)

	// You need to implement this function to create a new banner
	bannerID, err := s.Repo.CreateBanner(req.Context(), &banner)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"banner_id": bannerID})
}

// ---------------------------
func (s *Server1) UpdateBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)
	bannerID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var banner repository.Banner
	if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем баннер
	err := s.Repo.UpdateBanner(req.Context(), bannerID, &banner)
	if err != nil {
		http.Error(w, "Failed to update banner", http.StatusInternalServerError)
		return
	}

	// Теперь обновляем теги
	err = s.Repo.UpdateFeatureTags(req.Context(), bannerID, banner.FeatureID, banner.TagIDs)
	if err != nil {
		http.Error(w, "Failed to update feature tags", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server1) DeleteBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)
	bannerID, _ := strconv.ParseInt(vars["id"], 10, 64)

	// You need to implement this function to delete the banner
	err := s.Repo.DeleteBanner(req.Context(), bannerID)
	if err != nil {
		http.Error(w, "Эхх", http.StatusBadRequest)
		//switch {
		//case errors.Is(err, repository.ErrBannerNotFound):
		//	http.Error(w, "Banner not found", http.StatusNotFound)
		//default:
		//	http.Error(w, "Internal server error", http.StatusInternalServerError)
		//}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server1) DeleteByFeatureTagIDHandler(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var brokers = []string{
		//"kafka1",
		//"kafka2",
		//"kafka3",
		"127.0.0.1:9091",
		"127.0.0.1:9092",
	}

	initial.ProducerExample(brokers, req.URL.Path, req.Method)

}
