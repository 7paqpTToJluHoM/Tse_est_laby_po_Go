package main
import (
	"bufio"
	"fmt"
	"net"
)
func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Сервер слушает :8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка Accept:", err)
			continue
		}
		handleClient(conn)
	}
}
func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Получено: %s", msg)
	conn.Write([]byte("Сообщение получено\n"))
}