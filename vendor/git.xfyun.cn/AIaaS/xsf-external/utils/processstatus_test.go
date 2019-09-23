package utils

import (
	"testing"
)

func TestProcessStatus(t *testing.T) {
	var p ps
	t.Log(p.GetGoroutineID())
	t.Log(p.GetGoroutines())
	t.Log(p.GetPid())
	t.Log(p.GetUptime())
	t.Log(p.GetUser())
}
