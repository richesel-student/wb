package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

// Вывод будет результатом конкурентного слияния двух последовательностей:
// 1, 3, 5, 7 и 2, 4, 6, 8.
// Порядок значений внутри каждого канала сохраняется,
// но порядок между значениями из разных каналов недетерминирован.
//
// Данный конвейер с использованием select отрабатывает в конкурентном режиме.
// В каждой итерации select выбирает либо канал a, либо канал b,
// в зависимости от того, какой из них готов в данный момент.
//
// При выполнении case v, ok := <-a происходит чтение из канала a.
// Если ok == true, считанное значение передаётся в выходной канал c.
// Если ok == false, значит канал a закрыт, и он исключается из select
// путём присваивания a = nil (аналогично для канала b).
//
// Когда оба входных канала a и b завершены (a == nil и b == nil),
// выходной канал c закрывается, цикл завершается и горутина выходит.
