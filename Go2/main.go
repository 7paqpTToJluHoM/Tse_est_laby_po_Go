package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func ReadLine(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
func ReadInt(prompt string) int {
	for {
		input := ReadLine(prompt)
		number, err := strconv.Atoi(input)
		if err == nil {
			return number
		}
		fmt.Println("Ошибка! Введите целое число.")
	}
}
func ReadFloat(prompt string) float64 {
	for {
		input := ReadLine(prompt)
		number, err := strconv.ParseFloat(input, 64)
		if err == nil {
			return number
		}
		fmt.Println("Ошибка! Введите число.")
	}
}
func EvenOdd() {
	fmt.Println("------------Задание 1: четность числа------------")
	number := ReadInt("Введите целое число: ")
	if number%2 == 0 {
		fmt.Println("Число четное")
	} else {
		fmt.Println("Число нечетное")
	}
}
func NumberSign(number int) string {
	if number > 0 {
		return "Positive"
	}
	if number < 0 {
		return "Negative"
	}
	return "Zero"
}
func SignTask() {
	fmt.Println("------------Задание 2: позитив, негатив и зеро------------")
	number := ReadInt("Введите целое число: ")
	fmt.Println("Число - ", NumberSign(number))
}
func NumbersFromOneToTen() {
	fmt.Println("------------Задание 3: for------------")
	for i := 1; i <= 10; i++ {
		fmt.Println(i)
	}
}
func StringLength(str string) int { return len([]rune(str)) }
func StringTask() {
	fmt.Println("------------Задание 4: длина строки------------")
	fmt.Print("Введите строку: ")
	str := ReadLine("Введите строку: ")
	fmt.Println("Длина строки:", StringLength(str))
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 { return r.Width * r.Height }
func RectangleTask() {
	fmt.Println("------------Задание 5: площадь прямоугольника------------")
	width := ReadFloat("Введите ширину прямоугольника: ")
	height := ReadFloat("Введите высоту прямоугольника: ")
	rect := Rectangle{
		Width:  width,
		Height: height,
	}
	fmt.Println("Площадь прямоугольника:", rect.Area())
}
func Average(a, b int) float64 { return float64(a+b) / 2 }
func AverageTask() {
	fmt.Println("------------Задание 6: среднее значение двух чисел------------")
	a := ReadInt("Введите первое число: ")
	b := ReadInt("Введите второе число: ")
	fmt.Println("Среднее значение:", Average(a, b))
}
func main() {
	EvenOdd()
	SignTask()
	NumbersFromOneToTen()
	StringTask()
	RectangleTask()
	AverageTask()
}
