package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func factorial(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf("Горутина 1: факториал %d = %d\n", n, result)
}
func randomNumbers(count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		fmt.Printf("Горутина 2: случайное число %d: %d\n", i+1, rand.Intn(100))
		time.Sleep(70 * time.Millisecond)
	}
}
func seriesSum(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0.0
	for i := 1; i <= n; i++ {
		sum += 1.0 / float64(i)
		time.Sleep(40 * time.Millisecond)
	}
	fmt.Printf("Горутина 3: сумма ряда 1/%d = %.4f\n", n, sum)
}
func task1() {
	fmt.Println("\n----------Задание 1: параллельные горутины----------")
	var wg sync.WaitGroup
	startPar := time.Now()
	wg.Add(3)
	go factorial(10, &wg)
	go randomNumbers(5, &wg)
	go seriesSum(10, &wg)
	wg.Wait()
	parTime := time.Since(startPar)
	fmt.Printf("Время выполнения: %v\n", parTime)
}
func generateFibonacci(n int, ch chan<- int) {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a
		a, b = b, a+b
		time.Sleep(30 * time.Millisecond)
	}
	close(ch)
}
func readFibonacci(ch <-chan int, done chan<- bool) {
	fmt.Print("Числа Фибоначчи: ")
	for num := range ch {
		fmt.Printf("%d ", num)
	}
	fmt.Println()
	done <- true
}

func task2() {
	fmt.Println("\n----------Задание 2: каналы и числа Фибоначчи----------")
	ch := make(chan int)
	done := make(chan bool)
	go generateFibonacci(10, ch)
	go readFibonacci(ch, done)
	<-done
}
func task3() {
	fmt.Println("\n----------Задание 3: select и многоканальные операции------------")
	evenCh := make(chan int)
	oddCh := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			n := rand.Intn(100)
			if n%2 == 0 {
				evenCh <- n
			} else {
				oddCh <- n
			}
			time.Sleep(50 * time.Millisecond)
		}
		close(evenCh)
		close(oddCh)
	}()
	evenOpen, oddOpen := true, true
	for evenOpen || oddOpen {
		select {
		case n, ok := <-evenCh:
			if !ok {
				evenOpen = false
				evenCh = nil
				fmt.Println("[select] Канал чётных чисел закрыт")
			} else {
				fmt.Printf("[select] Чётное: %d\n", n)
			}
		case n, ok := <-oddCh:
			if !ok {
				oddOpen = false
				oddCh = nil
				fmt.Println("[select] Канал нечётных чисел закрыт")
			} else {
				fmt.Printf("[select] Нечётное: %d\n", n)
			}
		}
	}
	fmt.Println("select позволяет одновременно ждать данные из нескольких каналов")
}
func task4() {
	fmt.Println("\n----------Задание 4: синхронизация мьютексами----------")
	const goroutines = 10
	const increments = 1000
	var counterNoMutex int
	var wg1 sync.WaitGroup
	wg1.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg1.Done()
			for j := 0; j < increments; j++ {
				counterNoMutex++
			}
		}()
	}
	wg1.Wait()
	fmt.Printf("Без мьютекса: счётчик = %d (ожидалось %d) - возможна потеря значений\n", counterNoMutex, goroutines*increments)
	var counterMutex int
	var mu sync.Mutex
	var wg2 sync.WaitGroup
	wg2.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg2.Done()
			for j := 0; j < increments; j++ {
				mu.Lock()
				counterMutex++
				mu.Unlock()
			}
		}()
	}
	wg2.Wait()
	fmt.Printf("С мьютексом: счётчик = %d (ожидалось %d) — всегда верно\n", counterMutex, goroutines*increments)
}

type CalcRequest struct {
	A, B     float64
	Operator string
	ResultCh chan float64
}

func calculatorServer(requests <-chan CalcRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range requests {
		var result float64
		switch req.Operator {
		case "+":
			result = req.A + req.B
		case "-":
			result = req.A - req.B
		case "*":
			result = req.A * req.B
		case "/":
			if req.B != 0 {
				result = req.A / req.B
			} else {
				fmt.Println("Калькулятор: деление на ноль")
				result = 0
			}
		default:
			fmt.Printf("Калькулятор: неизвестная операция: %s\n", req.Operator)
		}
		req.ResultCh <- result
	}
}
func task5() {
	fmt.Println("\n----------Задание 5: многопоточный калькулятор----------")
	requests := make(chan CalcRequest)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go calculatorServer(requests, &wg)
	}
	type op struct {
		a, b float64
		op   string
	}
	ops := []op{
		{10, 5, "+"},
		{10, 5, "-"},
		{10, 5, "*"},
		{10, 5, "/"},
		{7, 0, "/"},
		{2, 8, "*"},
	}
	for _, o := range ops {
		resultCh := make(chan float64)
		requests <- CalcRequest{A: o.a, B: o.b, Operator: o.op, ResultCh: resultCh}
		result := <-resultCh
		fmt.Printf("[Клиент] %.1f %s %.1f = %.2f\n", o.a, o.op, o.b, result)
	}
	close(requests)
	wg.Wait()
}
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		time.Sleep(50 * time.Millisecond)
		reversed := reverseString(job)
		results <- fmt.Sprintf("[Воркер %d] \"%s\" -> \"%s\"", id, job, reversed)
	}
}
func task6(workerCount int) {
	fmt.Println("\n----------Задание 6: пул воркеров----------")
	inputFile := "input.txt"
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	fmt.Printf("Прочитано строк из %s: %d\n", inputFile, len(lines))
	jobs := make(chan string, len(lines))
	results := make(chan string, len(lines))
	var wg sync.WaitGroup
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}
	for _, line := range lines {
		jobs <- line
	}
	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()
	outFile, err := os.Create("results.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	fmt.Println("Результаты работы воркеров:")
	for res := range results {
		fmt.Println(res)
		fmt.Fprintln(writer, res)
	}
	writer.Flush()
	fmt.Println("Результаты записаны в results.txt")
}
func main() {
	task1()
	task2()
	task3()
	task4()
	task5()
	fmt.Println("\n----------Задание 6: ввод количества воркеров----------")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите кол-во воркеров: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	workerCount := 3
	if input != "" {
		fmt.Sscanf(input, "%d", &workerCount)
	}
	if workerCount < 1 {
		workerCount = 1
	}
	task6(workerCount)
}
