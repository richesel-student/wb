package main

import (
	"fmt"
	"time"
)

// остановка по условию
func stopCondition() {

	go func(limit int) {
		for i := 0; ; i++ {
			if i >= limit {
				fmt.Println("Достигнут лимит")
				return
			}
			fmt.Println("Работает горутина...", i)
			time.Sleep(300 * time.Millisecond)
		}
	}(5)
	time.Sleep(2 * time.Second)
}

func main() {
	stopCondition()

}
