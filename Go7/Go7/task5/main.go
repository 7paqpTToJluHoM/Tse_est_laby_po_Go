package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s%s — %v\n", r.Method, r.URL.Path, r.URL.RawQuery, time.Since(start),)
	})
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Привет, мир!")
}
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Получены данные: %+v\n", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Текущее время: %s\n", time.Now().Format(time.RFC3339))
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/data", dataHandler)
	mux.HandleFunc("/time", timeHandler)
	handler := loggingMiddleware(mux)
	fmt.Println("HTTP-сервер слушает :8082 (введите http://localhost:8082/hello, time    )")
	fmt.Println("Маршруты:")
	fmt.Println("  GET  /hello  — приветствие")
	fmt.Println("  POST /data   — приём JSON")
	fmt.Println("  GET  /time   — текущее время")
	fmt.Println("Логирование: метод, URL, время выполнения")
	if err := http.ListenAndServe(":8082", handler); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
