package main

import (
	"fmt"
	"sync"
)

type coutStruct struct {
	value int
	mu    sync.Mutex
}

func (c *coutStruct) Inc(n int) {
	c.mu.Lock()
	c.value += n
	c.mu.Unlock()

}

func main() {

	var c coutStruct

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc(1)

		}()

	}
	wg.Wait()
	fmt.Print(c.value)

}
