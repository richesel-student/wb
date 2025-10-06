package main

import (
	"fmt"
	"sync"
	"time"
)

func workers(id int, wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()
	for w := range ch {
		fmt.Printf("worker %d: %d\n", id, w)
	}
}

func main() {
	var wg sync.WaitGroup
	workerPools := 13
	ch := make(chan int)

	for i := 1; i <= workerPools; i++ {
		wg.Add(1)
		go workers(i, &wg, ch)
	}

	go func() {
		defer close(ch)
		deadline := time.After(3 * time.Second)
		for i := 1; i <= 15; i++ {
			select {
			case <-deadline:
				fmt.Println("Время истекло")
				return
			case ch <- i:

			}
		}

	}()

	wg.Wait()
	if workerPools != 0 {
		fmt.Println("\nГорутины завершили выполнение")
	}
}
