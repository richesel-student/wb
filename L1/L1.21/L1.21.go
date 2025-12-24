package main

import "fmt"

type Human struct {
	name   string
	age    int
	gender bool
}
type monkey struct {
	name   string
	age    int
	gender int
}

type Person interface {
	PersonInfo() string
}

type AdapterMonkey struct {
	m *monkey
}

func (a AdapterMonkey) PersonInfo() string {
	if a.m == nil {
		return ""
	}

	h := Human{
		name:   a.m.name,
		age:    a.m.age,
		gender: a.m.gender != 0,
	}
	return h.PersonInfo() // делегирование через адаптацию
}

func (m monkey) Genderint() string {
	var Gender string

	if m.gender == 1 {
		Gender = "мужского"
	} else {
		Gender = "женского"
	}
	return Gender
}

func (H Human) Genderstr() string {
	var Gender string

	if H.gender {
		Gender = "мужского"
	} else {
		Gender = "женского"
	}
	return Gender
}

func (H Human) AgeHuman() string {
	var agestr string
	ageint := H.age % 10

	switch ageint {
	case 1, 2, 3, 4:
		agestr = "года"
	case 11, 12, 13, 14:
		agestr = "лет"
	case 5, 6, 7, 8, 9, 0:
		agestr = "лет"

	}
	return agestr
}

func (H Human) PersonInfo() string {

	return fmt.Sprintf("Меня зовут %s. Я %s пола. Мне %d %s.", H.name, H.Genderstr(), H.age, H.AgeHuman())

}

func Print(p Person) {
	fmt.Println(p.PersonInfo())
}

func main() {

	monkey := &monkey{"Кинг-конг", 100, 1}
	Print(AdapterMonkey{m: monkey})

}

