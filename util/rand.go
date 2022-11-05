package util

import (
	"math/rand"
	"time"
)

//RandInt 返回一个(0,total]的随机数
func RandInt(total int) int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(total + 1)
	if num == 0 {
		return 1
	}
	return num
}
