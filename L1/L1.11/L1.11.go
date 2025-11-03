package main

import (
	"fmt"
)

type myset struct {
	A   map[int]bool
	B   map[int]bool
	res map[int]bool
}

func (m myset) init() myset {
	m.A = make(map[int]bool)
	m.B = make(map[int]bool)
	m.res = make(map[int]bool)
	for i := 1; i <= 3; i++ {
		m.A[i] = true

	}
	for i := 2; i <= 4; i++ {
		m.B[i] = true
	}
	return m

}
func (m myset) union() myset {
	for key, _ := range m.A {
		if m.A[key] == m.B[key] {
			m.res[key] = true
		}
	}

	return m

}


func main() {
	set := myset{}
	set = set.init()
	set = set.union()
	fmt.Print(set.res)
}
