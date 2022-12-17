// Modify: https://github.com/asyncins/mist/blob/master/mist.go

package mid

import (
	"crypto/rand"
	"math/big"
	"sync"
)

const saltBit = uint(8)                  // 随机因子二进制位数
const saltShift = uint(8)                // 随机因子移位数
const increasShift = saltBit + saltShift // 自增数移位数

type Mist struct {
	sync.Mutex       // 互斥锁
	increas    int64 // 自增数
	saltA      int64 // 随机因子一
	saltB      int64 // 随机因子二
}

func (c *Mist) Generate() int64 {
	c.Lock()
	c.increas++
	// 获取随机因子数值 ｜ 使用真随机函数提高性能
	randA, _ := rand.Int(rand.Reader, big.NewInt(255))
	c.saltA = randA.Int64()
	randB, _ := rand.Int(rand.Reader, big.NewInt(255))
	c.saltB = randB.Int64()
	// 通过位运算实现自动占位
	mist := (c.increas << increasShift) | (c.saltA << saltShift) | c.saltB
	c.Unlock()
	return mist
}

func UniqueId() int64 {
	mist := &Mist{increas: 1}
	return mist.Generate()
}
