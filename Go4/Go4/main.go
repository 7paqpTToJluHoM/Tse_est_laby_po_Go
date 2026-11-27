package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func averageAge(people map[string]int) float64 {
	if len(people) == 0 {
		return 0
	}

	sum := 0
	for _, age := range people {
		sum += age
	}
	return float64(sum) / float64(len(people))
}
func deletePerson(people map[string]int, name string) {
	delete(people, name)
	fmt.Printf("Запись для %s удалена.\n", name)
}
func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("========1: Карта с именами и возрастами========")
	people := make(map[string]int) // пустая карта
	people["Олег"] = 24
	people["Владислав"] = 24
	people["Алексей"] = 22
	people["Ольга"] = 29
	people["Генрих"] = 27
	people["Фридрих"] = 21
	people["Максим"] = 19
	people["Афанасий"] = 18
	people["Галина"] = 34
	people["Иероним"] = 56
	fmt.Println("Все записи:")
	for name, age := range people {
		fmt.Printf(" \t%s: %d лет\n", name, age)
	}
	////////////////////////////////////////////////////////////////////////////Задание 2////////////////////////////////////////////////////////////////////////////
	fmt.Println("========2: Средний возраст========")
	avg := averageAge(people)
	fmt.Printf("Средний возраст людей в карте: %.1f лет\n", avg)
	////////////////////////////////////////////////////////////////////////////Задание 3////////////////////////////////////////////////////////////////////////////
	fmt.Println("========3: Удаление записи========")
	name := readString("Введите имя для удаления: ")
	deletePerson(people, name)
	fmt.Println("После удаления:")
	for name, age := range people {
		fmt.Printf("\t%s: %d лет\n", name, age)
	}
	////////////////////////////////////////////////////////////////////////////Задание 4////////////////////////////////////////////////////////////////////////////
	fmt.Println("========4: Верхний регистр========")
	str := readString("Введите строку: ")
	karta_dlya_verha := map[rune]rune{
		'a': 'A', 'b': 'B', 'c': 'C', 'd': 'D', 'e': 'E', 'f': 'F',
		'g': 'G', 'h': 'H', 'i': 'I', 'j': 'J', 'k': 'K', 'l': 'L',
		'm': 'M', 'n': 'N', 'o': 'O', 'p': 'P', 'q': 'Q', 'r': 'R',
		's': 'S', 't': 'T', 'u': 'U', 'v': 'V', 'w': 'W', 'x': 'X',
		'y': 'Y', 'z': 'Z',
		'а': 'А', 'б': 'Б', 'в': 'В', 'г': 'Г', 'д': 'Д', 'е': 'Е',
		'ё': 'Ё', 'ж': 'Ж', 'з': 'З', 'и': 'И', 'й': 'Й', 'к': 'К',
		'л': 'Л', 'м': 'М', 'н': 'Н', 'о': 'О', 'п': 'П', 'р': 'Р',
		'с': 'С', 'т': 'Т', 'у': 'У', 'ф': 'Ф', 'х': 'Х', 'ц': 'Ц',
		'ч': 'Ч', 'ш': 'Ш', 'щ': 'Щ', 'ъ': 'Ъ', 'ы': 'Ы', 'ь': 'Ь',
		'э': 'Э', 'ю': 'Ю', 'я': 'Я',
	}
	upper := []rune(str)
	for i := range upper {
		for r := range karta_dlya_verha {
			if upper[i] == r {
				upper[i] = karta_dlya_verha[r]
				break
			}
		}
	}
	fmt.Printf("Результат: %s\n", string(upper))
	////////////////////////////////////////////////////////////////////////////Задание 5////////////////////////////////////////////////////////////////////////////
	fmt.Println("========5: Сумма чисел========")
	var sum int
	fmt.Println("Вводите числа (для завершения введите не число):")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		val, err := strconv.Atoi(input)
		if err != nil {
			break
		}
		sum += val
	}
	fmt.Printf("Сумма введённых чисел: %d\n", sum)
	////////////////////////////////////////////////////////////////////////////Задание 6////////////////////////////////////////////////////////////////////////////
	fmt.Println("========6: Обратный порядок========")
	var numbers []int
	fmt.Println("Вводите числа (для завершения введите не число):")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		val, err := strconv.Atoi(input)
		if err != nil {
			break
		}
		numbers = append(numbers, val)
	}
	if len(numbers) == 0 {
		fmt.Println("Массив пуст.")
	} else {
		fmt.Println("Числа в обратном порядке:")
		for i := len(numbers) - 1; i >= 0; i-- {
			fmt.Printf("\t%d", numbers[i])
		}
		fmt.Println()
	}
}
