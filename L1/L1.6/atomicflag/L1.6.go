package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// остановка через атомарный флаг
func stopAtomicFlag() {
	var wg sync.WaitGroup
	var stop atomic.Bool
	wg.Add(1)
	go worker(&stop, &wg)
	time.Sleep(2 * time.Second)
	stop.Store(true)
	wg.Wait()
}

func worker(stop *atomic.Bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for !stop.Load() {

		fmt.Println("Горутина работает...")
		time.Sleep(400 * time.Millisecond)
	}
	fmt.Println("Остановлен атомарным флагом")
}

func main() {
	stopAtomicFlag()
}
