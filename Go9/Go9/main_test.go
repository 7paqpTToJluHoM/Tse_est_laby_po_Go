package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestAPI(t *testing.T) *API {
	t.Helper()
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=jrl390fz dbname=lab8_test sslmode=disable"
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			age INTEGER NOT NULL CHECK (age >= 1 AND age <= 100)
		)
	`)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE users RESTART IDENTITY")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return &API{db: db}
}
func request(
	api *API,
	method string,
	url string,
	body interface{},
) *httptest.ResponseRecorder {
	var buffer bytes.Buffer
	if body != nil {
		json.NewEncoder(&buffer).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buffer)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	return rec
}
func TestCreateUser(t *testing.T) {
	api := setupTestAPI(t)
	body := User{Name: "Иван", Age: 20}
	rec := request(api, "POST", "/users", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("Ожидался 201, получен %d: %s", rec.Code, rec.Body.String())
	}
}
func TestGetUsers(t *testing.T) {
	api := setupTestAPI(t)
	request(api, "POST", "/users", User{Name: "Иван", Age: 20})
	rec := request(api, "GET", "/users", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("Ожидался 200, получен %d", rec.Code)
	}
	var response UsersResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Total != 1 {
		t.Errorf("Ожидался 1 пользователь, получено %d", response.Total)
	}
}
func TestGetUser(t *testing.T) {
	api := setupTestAPI(t)
	request(api, "POST", "/users", User{Name: "Иван", Age: 20})
	rec := request(api, "GET", "/users/1", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("Ожидался 200, получен %d", rec.Code)
	}
}
func TestUpdateUser(t *testing.T) {
	api := setupTestAPI(t)
	request(api, "POST", "/users", User{Name: "Иван", Age: 20})
	body := User{Name: "Пётр", Age: 25}
	rec := request(api, "PUT", "/users/1", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("Ожидался 200, получен %d", rec.Code)
	}
	var user User
	if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
		t.Fatal(err)
	}
	if user.Name != "Пётр" {
		t.Errorf("Имя не обновилось: %s", user.Name)
	}
}
func TestDeleteUser(t *testing.T) {
	api := setupTestAPI(t)
	request(api, "POST", "/users", User{Name: "Иван", Age: 20})
	rec := request(api, "DELETE", "/users/1", nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("Ожидался 204, получен %d", rec.Code)
	}
}
func TestValidation(t *testing.T) {
	api := setupTestAPI(t)
	body := User{Name: "", Age: 20}
	rec := request(api, "POST", "/users", body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Ожидался 400, получен %d", rec.Code)
	}
}
func TestPagination(t *testing.T) {
	api := setupTestAPI(t)
	request(api, "POST", "/users", User{Name: "Иван", Age: 20})
	request(api, "POST", "/users", User{Name: "Пётр", Age: 25})
	rec := request(api, "GET", "/users?page=1&limit=1", nil)
	var response UsersResponse
	if rec.Code != http.StatusOK {
		t.Fatalf("Ожидался 200, получен %d: %s",
			rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Users) != 1 {
		t.Errorf("Ожидался 1 пользователь")
	}
	if response.Total != 2 {
		t.Errorf("Ожидалось всего 2 пользователя")
	}
}
func TestFilterUsers(t *testing.T) {
	api := setupTestAPI(t)
	request(api, "POST", "/users", User{Name: "Иван", Age: 20})
	request(api, "POST", "/users", User{Name: "Пётр", Age: 25})
	rec := request(api, "GET", "/users?name=Иван", nil)
	var response UsersResponse
	if rec.Code != http.StatusOK {
		t.Fatalf("Ожидался 200, получен %d: %s",
			rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Total != 1 {
		t.Errorf("Ожидался 1 пользователь, получено %d", response.Total)
	}
}
func TestUserNotFound(t *testing.T) {
	api := setupTestAPI(t)
	rec := request(api, "GET", "/users/999", nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("Ожидался 404, получен %d", rec.Code)
	}
}
func TestMethodNotAllowed(t *testing.T) {
	api := setupTestAPI(t)

	rec := request(api, "PATCH", "/users", nil)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался 405, получен %d", rec.Code)
	}
}
