package db

import (
	"sync"
	"time"
)

const (
	sequenceBits  = 10                      // 序列号占用的位数（最多支持 1000 个序列号）
	sequenceMask  = int64(1<<sequenceBits - 1) // 序列号掩码，确保序列号不会超过 999
	timeShift     = sequenceBits            // 时间戳需要左移的位数
	maxSequence   = 999                     // 每毫秒最多生成 1000 个序列号
	startingEpoch = int64(1640995200000)    // 自定义开始时间戳 (2022-01-01 00:00:00 UTC)
)

type IDGenerator struct {
	mu        sync.Mutex
	lastTime  int64 // 上一次生成 ID 时的时间戳（毫秒）
	sequence  int64 // 当前毫秒的序列号
}

// 创建一个新的 ID 生成器
func newIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// GenerateID 生成一个唯一的 uint64 ID
// 暂只考虑时间戳和序列号部分
func (g *IDGenerator) generateID() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()
	if now == g.lastTime {
		// 在同一毫秒内，序列号递增
		g.sequence = (g.sequence + 1) & sequenceMask
		if g.sequence == 0 {
			// 如果序列号达到最大值，则等待到下一毫秒
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 如果是新的一毫秒，重置序列号
		g.sequence = 0
	}
	g.lastTime = now

	// 计算最终的 ID，将时间戳和序列号组合在一起
	timestamp := now - startingEpoch
	id := (timestamp << timeShift) | g.sequence
	return uint64(id)
}
