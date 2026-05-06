package or

import (
	"testing"
	"time"
)


func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func TestOr_Single(t *testing.T) {
	ch := sig(10 * time.Millisecond)

	start := time.Now()
	<-Or(ch)

	if time.Since(start) > 50*time.Millisecond {
		t.Fatal("Or did not return in time for single channel")
	}
}

func TestOr_Multiple(t *testing.T) {
	start := time.Now()

	<-Or(
		sig(2*time.Second),
		sig(500*time.Millisecond),
		sig(100*time.Millisecond), 
	)

	elapsed := time.Since(start)

	if elapsed > 300*time.Millisecond {
		t.Fatalf("Or took too long: %v", elapsed)
	}
}

func TestOr_ZeroChannels(t *testing.T) {
	if Or() != nil {
		t.Fatal("expected nil for zero channels")
	}
}

func TestOr_FirstWins(t *testing.T) {
	fast := sig(50 * time.Millisecond)
	slow := sig(500 * time.Millisecond)

	start := time.Now()
	<-Or(slow, fast)

	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Fatalf("expected fast channel to win, got %v", elapsed)
	}
}