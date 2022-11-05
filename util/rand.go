package util

import (
	"math/rand"
	"time"
)

type DefaultInt = int32
type RandPoolTyp = map[DefaultInt]DefaultInt //随机池类型

//RandInt 返回一个(0,total]的随机数
func RandInt(total DefaultInt) DefaultInt {
	rand.Seed(time.Now().UnixNano())
	num := rand.Int31n(total + 1)
	if num == 0 {
		return 1
	}
	return num
}

//RandOne 随机产出一个
//pool 奖池；k-奖品id v-奖品权重
func RandOne(pool RandPoolTyp) (poolId DefaultInt) {
	poolIds := make([]DefaultInt, 0, len(pool))
	weightList := make([]DefaultInt, 0, len(pool))
	weightTotal := DefaultInt(0)
	for id, w := range pool {
		poolIds = append(poolIds, id)
		weightTotal += w
		weightList = append(weightList, weightTotal)
	}
	randInt := RandInt(weightTotal)
	index := 0
	for i, w := range weightList {
		index = i
		if w >= randInt {
			break
		}
	}
	return poolIds[index]
}
