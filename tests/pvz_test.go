//go:build integration
// +build integration

package tests

import (
	"Avito/internal/config"
	"Avito/internal/delivery/handlers"
	dbN "Avito/internal/storage/db"
	dbrepo "Avito/internal/storage/repository/postgresql"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"net/http/httptest"
	"testing"
)

// -------------------------тесты ручек-------------------------------------------------------------------------
func TestBannerHandlers(t *testing.T) {
	// Настройка временной тестовой конфигурации для базы данных
	tempConfig := config.StorageConfig{
		Host:     "localhost",
		Port:     5433, // Порт вашей тестовой базы данных
		Username: "postgres",
		Password: "1111",
		Database: "TestAvito",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создание нового соединения с тестовой базой данных
	tempDatabase, err := dbN.NewDb(ctx, tempConfig)
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	defer tempDatabase.GetPool(ctx).Close()

	// Очистка базы данных перед запуском теста
	if err := cleanDatabase(ctx, tempDatabase.GetPool(ctx)); err != nil {
		t.Fatalf("failed to clean test database: %v", err)
	}

	// Создание объекта репозитория для тестовой базы данных
	bannerRepo := dbrepo.NewBannerRepo(tempDatabase)

	// Создание объекта сервера API с использованием тестового репозитория
	server := handlers.Server1{Serv: bannerRepo}

	// Тест создания баннера
	t.Run("CreateBannerHandler", func(t *testing.T) {
		// Подготовка данных запроса
		requestBody, err := json.Marshal(map[string]interface{}{
			"tag_ids":    []int{61, 71, 81, 1, 3},
			"feature_id": 1,
			"content": map[string]interface{}{
				"title": "New Banner",
				"text":  "Some text",
				"url":   "http://example.com",
			},
			"is_active": true,
		})
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		// Создание HTTP запроса
		req, err := http.NewRequest("POST", "http://localhost:9000/pvz", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", "admin_token") // Установка заголовка с токеном

		// Создание HTTP тестового сервера
		rr := httptest.NewRecorder()

		// Обработка запроса сервером API
		handler := http.HandlerFunc(server.CreateBanner)
		handler.ServeHTTP(rr, req)

		// Проверка кода состояния ответа
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		// Чтение и анализ ответа
		respBody := rr.Body.String()
		fmt.Println("Response:", respBody)
	})

	// Тест получения баннера пользователя
	t.Run("GetUserBannerHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "https://localhost:9000/user_banner?tag_id=1&feature_id=1&use_last_revision=false", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"key": "1"})

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(server.GetUserBanner)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		respBody := rr.Body.String()
		fmt.Println("Response:", respBody)
	})
}

// Очистка базы данных перед запуском теста
func cleanDatabase(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, "DELETE FROM featuretag")
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, "DELETE FROM banners")
	if err != nil {
		return err
	}

	// Другие операции по очистке базы данных

	return nil
}
