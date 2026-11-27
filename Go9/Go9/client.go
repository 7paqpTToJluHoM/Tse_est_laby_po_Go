package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const serverURL = "http://localhost:8082"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Client struct {
	httpClient *http.Client
	token      string
}

var scanner = bufio.NewScanner(os.Stdin)

func newClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	if !scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(scanner.Text())
}
func readInt(prompt string) (int, error) {
	input := readLine(prompt)
	return strconv.Atoi(input)
}
func (c *Client) login() {
	token := readLine("Введите токен авторизации: ")
	if token == "" {
		fmt.Println("Токен не может быть пустым.")
		return
	}
	c.token = token
	fmt.Println("Токен сохранён в памяти клиента.")
}
func (c *Client) request(method, path string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования JSON: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(method, serverURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка соединения с сервером: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var serverError ErrorResponse
		if err := json.Unmarshal(responseBody, &serverError); err == nil && serverError.Error != "" {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, serverError.Error)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}
	return responseBody, nil
}
func (c *Client) getUsers() {
	data, err := c.request(http.MethodGet, "/users?page=1&limit=100", nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	var result UsersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("Ошибка обработки JSON:", err)
		return
	}
	printUsers(result.Users)
	fmt.Printf("\nСтраница: %d | Лимит: %d | Всего пользователей: %d\n", result.Page, result.Limit, result.Total)
}
func (c *Client) getUser() {
	id, err := readInt("Введите ID пользователя: ")
	if err != nil || id <= 0 {
		fmt.Println("Некорректный ID.")
		return
	}
	data, err := c.request(http.MethodGet, fmt.Sprintf("/users/%d", id), nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		fmt.Println("Ошибка обработки JSON:", err)
		return
	}
	printUsers([]User{user})
}
func (c *Client) createUser() {
	name := readLine("Введите имя: ")
	age, err := readInt("Введите возраст: ")
	if err != nil {
		fmt.Println("Некорректный возраст.")
		return
	}
	user := User{
		Name: name,
		Age:  age,
	}
	data, err := c.request(http.MethodPost, "/users", user)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	var created User
	if err := json.Unmarshal(data, &created); err != nil {
		fmt.Println("Ошибка обработки JSON:", err)
		return
	}
	fmt.Println("Пользователь успешно добавлен:")
	printUsers([]User{created})
}
func (c *Client) updateUser() {
	id, err := readInt("Введите ID пользователя: ")
	if err != nil || id <= 0 {
		fmt.Println("Некорректный ID.")
		return
	}
	name := readLine("Введите новое имя: ")
	age, err := readInt("Введите новый возраст: ")
	if err != nil {
		fmt.Println("Некорректный возраст.")
		return
	}
	user := User{
		Name: name,
		Age:  age,
	}
	data, err := c.request(http.MethodPut, fmt.Sprintf("/users/%d", id), user)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	var updated User
	if err := json.Unmarshal(data, &updated); err != nil {
		fmt.Println("Ошибка обработки JSON:", err)
		return
	}
	fmt.Println("Пользователь успешно обновлён:")
	printUsers([]User{updated})
}
func (c *Client) deleteUser() {
	id, err := readInt("Введите ID пользователя: ")
	if err != nil || id <= 0 {
		fmt.Println("Некорректный ID.")
		return
	}
	_, err = c.request(http.MethodDelete, fmt.Sprintf("/users/%d", id), nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Пользователь успешно удалён.")
}
func printUsers(users []User) {
	if len(users) == 0 {
		fmt.Println("Пользователи не найдены.")
		return
	}
	fmt.Println()
	fmt.Printf("%-10s %-30s %-10s\n", "ID", "Имя", "Возраст")
	fmt.Println(strings.Repeat("-", 52))
	for _, user := range users {
		fmt.Printf("%-10d %-30s %-10d\n", user.ID, user.Name, user.Age)
	}
	fmt.Println(strings.Repeat("-", 52))
}
func showMenu() {
	fmt.Println("\nМЕНЮ")
	fmt.Println("1. Авторизация")
	fmt.Println("2. Получить список пользователей")
	fmt.Println("3. Получить пользователя по ID")
	fmt.Println("4. Добавить пользователя")
	fmt.Println("5. Обновить пользователя")
	fmt.Println("6. Удалить пользователя")
	fmt.Println("0. Выход")
	fmt.Println(" ")
}
func main() {
	client := newClient()
	fmt.Println("Клиент REST API")
	fmt.Println("Сервер:", serverURL)
	for {
		showMenu()
		choice := readLine("Выберите действие: ")
		switch choice {
		case "1":
			client.login()
		case "2":
			client.getUsers()
		case "3":
			client.getUser()
		case "4":
			client.createUser()
		case "5":
			client.updateUser()
		case "6":
			client.deleteUser()
		case "0":
			fmt.Println("Завершение работы клиента.")
			return
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
