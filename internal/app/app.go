package app

import (
	"Avito/internal/config"
	api "Avito/internal/delivery/handlers"
	routes "Avito/internal/delivery/routes"
	"Avito/internal/storage/db"
	pp "Avito/internal/storage/repository/postgresql"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	securePort   = ":9000"
	insecurePort = ":9001"
)

func Run() {

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to convert PORT to integer:", err)
	}
	config := config.StorageConfig{
		Host:     os.Getenv("HOST"),
		Port:     port,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("PASSWORD"),
		Database: os.Getenv("DBNAME"),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	bannerRepo := pp.NewBannerRepo(database)
	implementation := api.Server1{Repo: bannerRepo}

	go serveSecure(implementation)
	serveInsecure()
}

func serveSecure(implementation api.Server1) {
	secureMux := http.NewServeMux()
	secureMux.Handle("/", routes.CreateRouter(implementation))

	log.Printf("Listening on port %s...\n", securePort)
	if err := http.ListenAndServeTLS(securePort, "internal/app/server.crt", "internal/app/server.key", secureMux); err != nil {
		log.Fatal(err)
	}
}

func serveInsecure() {
	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostParts := strings.Split(req.Host, ":")
		host := hostParts[0]

		// Формируем целевой URL с портом 9000 для HTTPS
		target := "https://" + host + securePort + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, target, http.StatusTemporaryRedirect)
	})

	log.Printf("Listening on port %s...\n", insecurePort)
	if err := http.ListenAndServe(insecurePort, redirectHandler); err != nil {
		log.Fatal(err)
	}
}