package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(m map[int]int, i int, wg *sync.WaitGroup, mutex *sync.Mutex) {

	rand.Seed(time.Now().UnixNano())
	j := rand.Intn(100)
	mutex.Lock()
	defer mutex.Unlock()
	m[i] = j
	wg.Done()

}

func printer(m map[int]int) {

	for key, value := range m {
		fmt.Printf("%d -> %d\n", key, value)
	}

}

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	m := make(map[int]int)

	for i := 1; i < 1001; i++ {
		wg.Add(1)
		go worker(m, i, &wg, &mutex)
	}

	wg.Wait()
	printer(m)

}
