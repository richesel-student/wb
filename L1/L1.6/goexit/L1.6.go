package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func printer(i int, wg *sync.WaitGroup) {
	defer wg.Done()

	if i > 5 {

		fmt.Printf("\nГорутина №%d не была выполнена.", i)
		runtime.Goexit()
		fmt.Printf("это сообщение никогда не выведется.")
	}
	fmt.Printf("\nГорутина №%d выполняется.", i)
	time.Sleep(100 * time.Microsecond)

}

func goroutines(wg *sync.WaitGroup) {

	for i := 1; i < 10; i++ {

		wg.Add(1)

		go printer(i, wg)

	}

}

func main() {
	var wg sync.WaitGroup
	goroutines(&wg)
	wg.Wait()

}
