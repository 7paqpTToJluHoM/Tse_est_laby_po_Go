package main
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)
var wg sync.WaitGroup
func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Сервер слушает :8081")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nПолучен сигнал, останавливаем сервер")
		listener.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil { break }
		wg.Add(1)
		go handleConn(conn)
	}
	wg.Wait()
	fmt.Println("Все соединения завершены. Сервер остановлен")
}
func handleConn(conn net.Conn) {
	defer wg.Done()
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}
	fmt.Printf("Получено: %s", msg)
	conn.Write([]byte("Сообщение получено\n"))
}