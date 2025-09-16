package main

import "fmt"

type Human struct {
	name   string
	age    int
	gender bool
}
type Action struct {
	Human
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

	return fmt.Sprintf("Mеня зовут %s. Я %s пола. Мне %d %s.", H.name, H.Genderstr(), H.age, H.AgeHuman())

}

func main() {

	Rinat := Human{"Ринат", 32, true}

	clone := Action{
		Human: Rinat,
	}
	fmt.Print(clone.PersonInfo())

}
