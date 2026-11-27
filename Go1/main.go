package main

import (
	"fmt"
	"time"
)

func Time() {
	var eto_vremya = time.Now()
	fmt.Println("------------Задание 1: текущее время и дата------------")
	fmt.Println("Полная дата и время:", eto_vremya.Format("02.01.2006 15:04:05"))
	fmt.Println("Полная дата и время в UTC формате:", eto_vremya.UTC().Format("02.01.2006 15:04:05"))
}
func Matesha() {
	var chislo int = 10
	var flot float64 = 10.01
	var stroka string = "Таки что вы от меня хотите? Я всего лишь строка"
	var logika bool = false
	fmt.Println("------------Задание 2: переменные различных типов------------")
	fmt.Printf("Целочисленное число:\n %d\n", chislo)
	fmt.Printf("Дробное число:\n %f\n", flot)
	fmt.Printf("Строка:\n %s\n", stroka)
	fmt.Printf("Ускорение свободного падения равно 10 м/(с*с)?\n %t\n", logika)
}
func Kratkost() {
	god_sssr := 1923
	imperator_tort := "Наполеон"
	odin_delit_na_dva := 0.5
	veroyatnost := false
	fmt.Println("------------Задание 3: краткое объявление чисел------------")
	fmt.Printf("Год образования СССР?\n %d\n", god_sssr)
	fmt.Printf("Император-торт?\n %s\n", imperator_tort)
	fmt.Printf("Сколько один делить на два?\n %f\n", odin_delit_na_dva)
	fmt.Printf("Есть ли у меня шанс сдать этот предмет?\n %t\n", veroyatnost)
}
func Matesha2() {
	a := 10
	b := 2
	fmt.Println("------------Задание 4: арифметика------------")
	fmt.Printf("Есть числа: %d %d\n", a, b)
	fmt.Printf("%d - %d= %d\n", a, b, a-b)
	fmt.Printf("%d + %d= %d\n", a, b, a+b)
	fmt.Printf("%d * %d= %d\n", a, b, a*b)
	fmt.Printf("%d / %d= %d\n", a, b, a/b)
	fmt.Printf("%d %% %d= %d\n", a, b, a%b)
}
func Matesha3() {
	a := 10.2
	b := 2.1
	fmt.Println("------------Задание 5: дробная арифметика------------")
	fmt.Printf("Есть числа: %f %f\n", a, b)
	fmt.Printf("%f - %f= %f\n", a, b, a-b)
	fmt.Printf("%f + %f= %f\n", a, b, a+b)
	fmt.Printf("%f * %f= %f\n", a, b, a*b)
	fmt.Printf("%f / %f= %f\n", a, b, a/b)
}
func Sr_znach() {
	a := 10.2
	b := 2.1
	c := 3.34
	fmt.Println("------------Задание 6: среднее значение трех чисел------------")
	fmt.Printf("Есть числа: %f %f %f\n", a, b, c)
	fmt.Printf("Их среднее значение = %f\n", (a+b+c)/3)
}

func main() {
	Time()
	Matesha()
	Kratkost()
	Matesha2()
	Matesha3()
	Sr_znach()
}
