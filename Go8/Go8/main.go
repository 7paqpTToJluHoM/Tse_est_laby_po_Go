package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type API struct {
	db *sql.DB
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type UsersResponse struct {
	Users []User `json:"users"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}

func initDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=jrl390fz dbname=lab8 sslmode=disable"
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INTEGER NOT NULL CHECK (age >= 0 AND age <= 150)
	);
	`

	if _, err := db.Exec(query); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}
func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Ошибка отправки JSON:", err)
	}
}
func validateUser(user User) string {
	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return "Имя не может быть пустым"
	}
	if len([]rune(user.Name)) > 100 {
		return "Имя не должно превышать 100 символов"
	}
	if user.Age < 1 || user.Age > 100 {
		return "Возраст должен быть от 1 до 100"
	}
	return ""
}
func decodeUser(w http.ResponseWriter, r *http.Request) (User, error) {
	var user User
	decoder := json.NewDecoder(
		http.MaxBytesReader(w, r.Body, 1<<20),
	)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		return user, err
	}
	var extra interface{}
	err := decoder.Decode(&extra)
	if !errors.Is(err, io.EOF) {
		if err == nil {
			return user, fmt.Errorf("ожидался только один JSON-объект")
		}
		return user, err
	}
	return user, nil
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf(
			"[%s] %s%s — %v\n",
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			time.Since(start),
		)
	})
}
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "users" {
		switch r.Method {
		case http.MethodGet:
			api.getUsers(w, r)

		case http.MethodPost:
			api.createUser(w, r)

		default:
			w.Header().Set("Allow", "GET, POST")
			sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
		return
	}
	if strings.HasPrefix(path, "users/") {
		idStr := strings.TrimPrefix(path, "users/")
		if idStr == "" || strings.Contains(idStr, "/") {
			sendError(w, "Некорректный маршрут", http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			sendError(w, "Некорректный ID пользователя", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			api.getUser(w, r, id)

		case http.MethodPut:
			api.updateUser(w, r, id)

		case http.MethodDelete:
			api.deleteUser(w, r, id)

		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
		return
	}
	sendError(w, "Маршрут не найден", http.StatusNotFound)
}
func (api *API) getUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := 1
	limit := 10
	var err error
	if query.Get("page") != "" {
		page, err = strconv.Atoi(query.Get("page"))
		if err != nil || page < 1 {
			sendError(w, "Некорректный номер страницы", http.StatusBadRequest)
			return
		}
	}
	if query.Get("limit") != "" {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil || limit < 1 || limit > 100 {
			sendError(w, "Лимит должен быть от 1 до 100", http.StatusBadRequest)
			return
		}
	}
	if page-1 > int(^uint(0)>>1)/limit {
		sendError(w, "Слишком большой номер страницы", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(query.Get("name"))
	var minAge, maxAge *int
	if query.Get("min_age") != "" {
		value, err := strconv.Atoi(query.Get("min_age"))
		if err != nil || value < 1 || value > 100 {
			sendError(w, "Некорректный минимальный возраст", http.StatusBadRequest)
			return
		}
		minAge = &value
	}
	if query.Get("max_age") != "" {
		value, err := strconv.Atoi(query.Get("max_age"))
		if err != nil || value < 1 || value > 100 {
			sendError(w, "Некорректный максимальный возраст", http.StatusBadRequest)
			return
		}
		maxAge = &value
	}
	if minAge != nil && maxAge != nil && *minAge > *maxAge {
		sendError(w, "Минимальный возраст больше максимального", http.StatusBadRequest)
		return
	}
	where := []string{"1=1"}
	args := []interface{}{}
	if name != "" {
		args = append(args, "%"+name+"%")
		where = append(where, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if minAge != nil {
		args = append(args, *minAge)
		where = append(where, fmt.Sprintf("age >= $%d", len(args)))
	}
	if maxAge != nil {
		args = append(args, *maxAge)
		where = append(where, fmt.Sprintf("age <= $%d", len(args)))
	}
	condition := strings.Join(where, " AND ")
	var total int
	countQuery := "SELECT COUNT(*) FROM users WHERE " + condition
	err = api.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		sendError(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}
	offset := (page - 1) * limit
	selectQuery := `
		SELECT id, name, age
		FROM users
		WHERE ` + condition + fmt.Sprintf(`
		ORDER BY id
		LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)
	selectArgs := append(args, limit, offset)
	rows, err := api.db.Query(selectQuery, selectArgs...)
	if err != nil {
		sendError(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			sendError(w, "Ошибка чтения пользователей", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		sendError(w, "Ошибка чтения пользователей", http.StatusInternalServerError)
		return
	}
	sendJSON(w, UsersResponse{
		Users: users,
		Page:  page,
		Limit: limit,
		Total: total,
	}, http.StatusOK)
}
func (api *API) getUser(w http.ResponseWriter, r *http.Request, id int) {
	var user User
	err := api.db.QueryRow(
		"SELECT id, name, age FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Age)
	if err == sql.ErrNoRows {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		sendError(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}
	sendJSON(w, user, http.StatusOK)
}
func (api *API) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUser(w, r)
	if err != nil {
		sendError(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if message := validateUser(user); message != "" {
		sendError(w, message, http.StatusBadRequest)
		return
	}
	user.Name = strings.TrimSpace(user.Name)
	err = api.db.QueryRow(
		"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id",
		user.Name,
		user.Age,
	).Scan(&user.ID)
	if err != nil {
		sendError(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}
	sendJSON(w, user, http.StatusCreated)
}
func (api *API) updateUser(w http.ResponseWriter, r *http.Request, id int) {
	user, err := decodeUser(w, r)
	if err != nil {
		sendError(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if message := validateUser(user); message != "" {
		sendError(w, message, http.StatusBadRequest)
		return
	}
	user.Name = strings.TrimSpace(user.Name)
	result, err := api.db.Exec(
		"UPDATE users SET name = $1, age = $2 WHERE id = $3",
		user.Name,
		user.Age,
		id,
	)
	if err != nil {
		sendError(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		sendError(w, "Ошибка проверки обновления", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	user.ID = id
	sendJSON(w, user, http.StatusOK)
}
func (api *API) deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	result, err := api.db.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		sendError(w, "Ошибка удаления пользователя", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		sendError(w, "Ошибка проверки удаления", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	api := &API{db: db}
	handler := loggingMiddleware(api)
	fmt.Println("REST API сервер запущен на :8082")
	fmt.Println("База данных: PostgreSQL (lab8)")
	fmt.Println("Маршруты:")
	fmt.Println("GET    /users")
	fmt.Println("GET    /users/{id}")
	fmt.Println("POST   /users")
	fmt.Println("PUT    /users/{id}")
	fmt.Println("DELETE /users/{id}")
	if err := http.ListenAndServe(":8082", handler); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
