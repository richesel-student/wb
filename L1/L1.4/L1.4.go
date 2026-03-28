package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func workers(id int, wg *sync.WaitGroup, ch chan int) {

	for w := range ch {
		fmt.Printf("worker %d: %d\n", id, w)
	}

	defer wg.Done()

}

func myerror(err error, countWorker *int) {
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if *countWorker <= 0 {
		fmt.Print(errors.New("Количество воркеров не должно быть меньше или равно 0"))
		return
	}

}

func main() {

	var wg sync.WaitGroup
	var workerPools int
	_, err := fmt.Scanln(&workerPools)

	ch := make(chan int, 100)
	myerror(err, &workerPools)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for i := 1; i <= workerPools; i++ {
		wg.Add(1)
		go workers(i, &wg, ch)
	}

	go func() {
		defer close(ch)
		for i := 1; i <= 15; i++ {
			select {
			case <-sigs:
				fmt.Printf("\nВвод команды Ctrl+C")
				return
			case ch <- i:
				time.Sleep(1 * time.Second)
			}
		}
	}()
	wg.Wait()

	if workerPools != 0 {
		fmt.Println("\nГорутины завершили выполнение")

	}

}
