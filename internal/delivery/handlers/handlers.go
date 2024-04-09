package handlers

import (
	"Avito/internal/storage/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	tagID, err := strconv.ParseInt(req.URL.Query().Get("tag_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid tag_id", http.StatusBadRequest)
		return
	}
	featureID, err := strconv.ParseInt(req.URL.Query().Get("feature_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid feature_id", http.StatusBadRequest)
		return
	}
	useLastRevision, err := strconv.ParseBool(req.URL.Query().Get("use_last_revision"))
	if err != nil {
		http.Error(w, "Invalid use_last_revision", http.StatusBadRequest)
		return
	}

	// You need to implement this function to retrieve the banner for the user
	banner, err := s.Repo.GetUserBanner(req.Context(), tagID, featureID, useLastRevision)
	if err != nil {
		http.Error(w, "Эхх", http.StatusBadRequest)

		//switch {
		//case errors.Is(err, repository.ErrUnauthorized):
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//case errors.Is(err, repository.ErrForbidden):
		//	http.Error(w, "Forbidden", http.StatusForbidden)
		//case errors.Is(err, repository.ErrBannerNotFound):
		//	http.Error(w, "Banner not found", http.StatusNotFound)
		//default:
		//	http.Error(w, "Internal server error", http.StatusInternalServerError)
		//}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banner)
}

func (s *Server1) GetBanners(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	featureID, _ := strconv.ParseInt(req.URL.Query().Get("feature_id"), 10, 64)
	tagID, _ := strconv.ParseInt(req.URL.Query().Get("tag_id"), 10, 64)
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))

	// You need to implement this function to retrieve banners with filtering
	banners, err := s.Repo.GetBanners(req.Context(), featureID, tagID, limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

func (s *Server1) CreateBanner(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var banner repository.Banner
	if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// You need to implement this function to create a new banner
	bannerID, err := s.Repo.CreateBanner(req.Context(), &banner)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"banner_id": bannerID})
}

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

	// You need to implement this function to update the banner
	err := s.Repo.UpdateBanner(req.Context(), bannerID, &banner)
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
