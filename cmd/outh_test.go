package main

import (
	auth "Lessons/internal/outh"
	"Lessons/internal/user"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDb() *gorm.DB {
	// Загружаем переменные окружения из указанного .env файла.
	err := godotenv.Load(".env")
	if err != nil {
		// Если файл не найден или произошла другая ошибка, завершить программу.
		panic(err)
	}

	fmt.Println("Подключение к базе данных успешно!")

	// Подключаемся к базе данных PostgreSQL.
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Включаем логирование запросов.
	})
	if err != nil {
		// Если подключение не удалось, завершить программу.
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@gmail.com",
		Password: "$2a$10$kZVz97tThenmUzHVYxPG2e8KzhnZVG/8yHsbT4PKxwktbrO8JLsE2",
		Name:     "test",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "test@gmail.com").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	//prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@gmail.com",
		Password: "1",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expecting 200, got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatalf("Expecting token, got %s", resData.Token)
	}

	removeData(db)

}

func TestLoginFailed(t *testing.T) {
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@gmail.com",
		Password: "..",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Expecting 401, got %d", res.StatusCode)
	}
	removeData(db)

}
