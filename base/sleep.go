package base

import (
	"math/rand"
	"time"
)

// 随机休眠
func RandSleep(min int64, max int64) {
	num := time.Duration((rand.Int63n(max) + min))
	time.Sleep(num * time.Millisecond)
}
