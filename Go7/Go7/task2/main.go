package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)
func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()
	fmt.Print("Введите сообщение: ")
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	conn.Write([]byte(msg))
	resp := make([]byte, 1024)
	n, _ := conn.Read(resp)
	fmt.Printf("Ответ сервера: %s", string(resp[:n]))
}