package main

import (
	"fmt"
	"time"
)

func MySleep(t time.Duration) {
	timer := time.NewTimer(t)
	<-timer.C

}

func main() {
	start := time.Now()
	MySleep(3 * time.Second)
	end := time.Since(start)
	fmt.Println(end)

}
