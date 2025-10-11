package main

import (
	"fmt"

	"sync"
	"time"
)

// остановка горутин через канал уведомлений
func goroutines(i int, stop chan struct{}, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)

	defer wg.Done()
	for {
		select {
		case <-stop:
			fmt.Printf("\nГорутина №%d остановилась.", i)
			return
		default:
			fmt.Printf("\nГорутина №%d работает.", i)
		}

	}

}

func main() {

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 1; i < 10; i++ {
		wg.Add(1)
		go goroutines(i, stop, &wg)
		if i == 5 {
			close(stop)

		}

	}
	wg.Wait()

}
