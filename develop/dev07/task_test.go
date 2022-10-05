package main

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(2*time.Second),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	if int(time.Since(start).Seconds()) != 1 {
		t.Error("function is not working correctly")
	}

	start = time.Now()
	<-or(sig(10 * time.Millisecond))
	if int(time.Since(start).Milliseconds()) != 10 {
		t.Error("function is not working correctly")
	}

	start = time.Now()
	<-or()
	if int(time.Since(start).Milliseconds()) != 0 {
		t.Error("function is not working correctly")
	}
}
