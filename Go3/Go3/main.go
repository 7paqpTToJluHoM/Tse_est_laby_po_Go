package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"Go3/mathutils"
	"Go3/stringutils"
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
func FactorialTask() {
	fmt.Println("------------Задание 1: факториал------------")
	number := ReadInt("Введите число: ")
	if number < 0 {
		fmt.Println("Факториал отрицательного числа не определен.")
		return
	}
	fmt.Println("Факториал:", mathutils.Faktorial_vychisli_mne_bistro(number))
}
func StringReverseTask() {
	fmt.Println("------------Задание 2: переворот строки------------")

	str := ReadLine("Введите строку: ")
	fmt.Println("Перевернутая строка:", stringutils.Perevorot_stroki(str))
}
func ArrayTask() {
	fmt.Println("------------Задание 3: массив из 5 чисел------------")
	arr := [5]int{10, 20, 30, 40, 50}
	fmt.Println("Массив:")
	for _, value := range arr {
		fmt.Println(value)
	}
}
func SliceTask() {
	fmt.Println("------------Задание 4: срез------------")
	arr := [5]int{10, 20, 30, 40, 50}
	slice := arr[:]
	fmt.Println("Исходный срез:", slice)
	slice = append(slice, 60)
	fmt.Println("После добавления:", slice)
	index := 2
	slice = append(slice[:index], slice[index+1:]...)
	fmt.Println("После удаления элемента:", slice)
}
func LongestStringTask() {
	fmt.Println("------------Задание 5: самая длинная строка------------")
	stringsArray := []string{
		"Написать программу",
		"которая создает срез из строк",
		"и находит самую длинную",
		"строку",
	}
	longest := stringsArray[0]
	for _, str := range stringsArray {
		if len([]rune(str)) > len([]rune(longest)) {
			longest = str
		}
	}
	fmt.Println("Срез строк:", stringsArray)
	fmt.Println("Самая длинная строка:", longest)
	fmt.Println("Ее длина:", len([]rune(longest)))
}
func main() {
	FactorialTask()
	StringReverseTask()
	ArrayTask()
	SliceTask()
	LongestStringTask()
}
