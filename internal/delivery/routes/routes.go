package routes

import (
	handlers "Avito/internal/delivery/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const queryParamKey = "key"

func CreateRouter(implemetation handlers.Server1) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user_banner", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetUserBanner(w, req)
		default:
			fmt.Println("error")
		}
	})
	router.HandleFunc("/banner", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetBanners(w, req)
		case http.MethodPost:
			implemetation.CreateBanner(w, req)
		default:
			fmt.Println("error")
		}
	})
	router.HandleFunc(fmt.Sprintf("/banner/{id}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPatch:
			implemetation.UpdateBanner(w, req)
		case http.MethodDelete:
			implemetation.DeleteBanner(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}
