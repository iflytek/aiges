package utils

import (
	"fmt"
	"testing"
)

func TestLogSidGenerator_GenerateSid(t *testing.T) {
	inst := LogSidGenerator{}
	fmt.Println(inst.GenerateSid("xxx"))
}
