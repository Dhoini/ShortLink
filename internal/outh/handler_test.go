package auth

import (
	"Lessons/configs"
	"Lessons/internal/user"
	"Lessons/pkg/db"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func bootstrap() (*AouthHendler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDB,
	})
	handler := AouthHendler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: NewUserService(userRepo),
	}

	return &handler, mock, nil

}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password", "name"})

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}

	data, _ := json.Marshal(&RegisterRequest{
		Email:    "test@gmail.com",
		Password: "1",
		Name:     "test",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)

	handler.Register()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 201 OK, but got %d", w.Code)
		return
	}
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("test@gmail.com", "$2a$10$kZVz97tThenmUzHVYxPG2e8KzhnZVG/8yHsbT4PKxwktbrO8JLsE2")

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}

	data, _ := json.Marshal(&LoginRequest{
		Email:    "test@gmail.com",
		Password: "1",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)

	handler.Login()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200 OK, but got %d", w.Code)
		return
	}
}
