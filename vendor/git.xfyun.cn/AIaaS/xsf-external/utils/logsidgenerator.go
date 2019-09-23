package utils

import (
	"strconv"
	"sync/atomic"
	"time"
)

type LogSidGenerator struct {
	count int64
}

func (s *LogSidGenerator) GenerateSid(tag string) (sid string) {
	return "log@" +
		tag +
		strconv.FormatInt(atomic.AddInt64(&s.count, 1), 10) +
		strconv.FormatInt(time.Now().Unix(), 10)
}
