//go:build integration
// +build integration

package tests

import (
	"Avito/internal/config"
	"Avito/tests/posgresql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	db *posgresql.TDB
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Загрузка переменных среды
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "TestAvito")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "1111")

	// Преобразование порта в int
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("Error converting port to int:", err)
	}

	// Создание объекта конфигурации базы данных
	dbConfig := config.StorageConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	// Создание подключения к базе данных
	db = posgresql.NewFromEnv(dbConfig)
}
