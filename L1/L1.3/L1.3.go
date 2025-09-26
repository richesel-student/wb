package main

import (
	"errors"
	"fmt"
	"sync"
)

func work(id int, wg *sync.WaitGroup, ch chan int) {

	for w := range ch {
		fmt.Printf("worker %d: %d\n", id, w)
	}
	wg.Done()

}

func myerror(err error, countWorker *int) {
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if *countWorker == 0 {
		fmt.Print(errors.New("Количество воркеров не должно рано 0"))
	}

}

func main() {

	var wg sync.WaitGroup
	var countWorker int
	_, err := fmt.Scanln(&countWorker)

	ch := make(chan int, 100)
	myerror(err, &countWorker)

	for i := 1; i <= countWorker; i++ {
		wg.Add(1)
		go work(i, &wg, ch)
	}
	for i := 1; i <= countWorker; i++ {
		ch <- i

	}
	close(ch)

	wg.Wait()

	if countWorker != 0 {
		fmt.Println("\nГорутины завершили выполнение")

	}

}
