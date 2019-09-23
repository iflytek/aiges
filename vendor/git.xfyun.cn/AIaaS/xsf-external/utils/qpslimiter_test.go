package utils

import (
	"testing"
)

func TestCheckQps(t *testing.T) {
	limiter, limiter_err := NewQpsLimiter(100, 1)
	if limiter_err != nil {
		t.Fatal(limiter_err)
	}
	for {
		if limiter.CheckQps() {
			t.Log("beyond qps limit.")
			break
		}
	}
}
