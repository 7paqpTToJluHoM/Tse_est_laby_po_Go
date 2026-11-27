package main

import (
	"fmt"
)

// //////////////////////////////////////////////////////////////////////////////Person и вывод информации////////////////////////////////////////
type Person struct {
	name string
	age  int
}

func (p Person) Vyvod_infy() {
	fmt.Printf("\nЕсть человек по имени: %s, его возраст: %d", p.name, p.age)
}
func (p *Person) Birthday() {
	p.age += 1
	fmt.Printf("\nКажется у кого-то день рождения! Проверяйте вывод!")
}

// //////////////////////////////////////////////////////////////////////////////Circle и радиус с площадью////////////////////////////////////////
type Circle struct {
	radius float64
}

func (c Circle) Vychislit_ploshad() float64 {
	return c.radius * c.radius * 3.14
}

// //////////////////////////////////////////////////////////////////////////////Интерфейс Shape и метод Area////////////////////////////////////////
type Rectangle struct {
	Shirina float64
	Vysota  float64
}
type Shape interface {
	Area() float64
}

func (r Rectangle) Area() float64 {
	return r.Shirina * r.Vysota
}
func (c Circle) Area() float64 {
	return c.radius * c.radius * 3.14
}

// //////////////////////////////////////////////////////////////////////////////Срез интерфейсов Shape и вывод через Area////////////////////////////////////////
func printAreas(shapes []Shape) {
	for i, s := range shapes {
		fmt.Printf("\nФигура %d: площадь = %f\n", i+1, s.Area())
	}
}

// //////////////////////////////////////////////////////////////////////////////Интерфейс Stringer и Book////////////////////////////////////////
type Book struct {
	title  string
	author string
	year   int
}
type Stringer interface {
	String() string
}

func (b Book) String() string {
	return fmt.Sprintf("\nКнига \"%s\" (автор: %s, год: %d)", b.title, b.author, b.year)
}

func main() {
	// /////////////////////////////////////////////////////////////////////////////1: Person и вывод информации/////////////////////////////////////////
	fmt.Printf("----------Задание 1: Person и вывод информации о нем----------")
	Chelovek := Person{name: "Иван", age: 25}
	Chelovek.Vyvod_infy()
	// /////////////////////////////////////////////////////////////////////////////2: Birthday/////////////////////////////////////////
	fmt.Printf("\n\n----------Задание 2: Person и Birthday----------")
	Chelovek.Birthday()
	Chelovek.Vyvod_infy()
	Chelovek.Birthday()
	Chelovek.Vyvod_infy()
	// /////////////////////////////////////////////////////////////////////////////3: Circle/////////////////////////////////////////
	fmt.Printf("\n\n----------Задание 3: Circle и площадь----------")
	Krug := Circle{radius: 5}
	fmt.Printf("\nПлощадь круга равна: %f", Krug.Vychislit_ploshad())
	// /////////////////////////////////////////////////////////////////////////////4 и 5: срез интерфейсов и Area для Rectangle с Circle/////////////////////////////////////////
	fmt.Printf("\n\n----------Задание 4 и 5: Срез интерфейсов и Area для Rectangle с Circle----------")
	shapes := []Shape{
		Circle{radius: 3},
		Rectangle{Shirina: 4, Vysota: 5},
		Circle{radius: 1.5},
		Rectangle{Shirina: 10, Vysota: 2},
	}
	printAreas(shapes)
	// /////////////////////////////////////////////////////////////////////////////6: Book/////////////////////////////////////////
	fmt.Printf("\n\n----------Задание 6: Book----------")
	book := Book{title: "Война и Мир", author: "Лев Толстой", year: 1869}
	var s Stringer = book
	fmt.Println(s.String())
	fmt.Println(book)
}
