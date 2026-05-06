package or

import (
	"fmt"
	"time"
)

func ExampleOr() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
	)

	fmt.Printf("done after %v\n", time.Since(start).Round(time.Millisecond))
}