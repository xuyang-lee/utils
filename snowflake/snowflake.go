package snowflake

import (
	"sync"
	"time"
)

// Snowflake Id生成器结构体
//type Snowflake struct {
//	sync.Mutex         // 同步锁
//	timestamp    int64 // 时间戳
//	workerId     int64 // 工作节点Id
//	datacenterId int64 // 数据中心Id
//	sequence     int64 // 序列号
//}
//
//// NewSnowflake 创建新的Id生成器
//func NewSnowflake(workerId, datacenterId int64) (*Snowflake, error) {
//	if workerId < 0 || workerId > maxWorkerId {
//		return nil, errors.New("worker Id out of range")
//	}
//	if datacenterId < 0 || datacenterId > maxDatacenterId {
//		return nil, errors.New("datacenter Id out of range")
//	}
//
//	return &Snowflake{
//		timestamp:    0,
//		workerId:     workerId,
//		datacenterId: datacenterId,
//		sequence:     0,
//	}, nil
//}
//
//// 生成Id
//func (s *Snowflake) NextId() (int64, error) {
//	s.Lock()
//	defer s.Unlock()
//
//	now := time.Now().UnixMilli()
//	if s.timestamp == now {
//		s.sequence = (s.sequence + 1) & sequenceMask
//		if s.sequence == 0 {
//			for now <= s.timestamp {
//				now = time.Now().UnixMilli()
//			}
//		}
//	} else {
//		s.sequence = 0
//	}
//
//	if now < s.timestamp {
//		return 0, errors.New("clock moved backwards")
//	}
//
//	s.timestamp = now
//	id := ((now - epoch) << timestampLeftShift) |
//		(s.datacenterId << datacenterIdShift) |
//		(s.workerId << workerIdShift) |
//		s.sequence
//
//	return id, nil
//}

// Snowflake 结构体定义。
type Snowflake struct {
	c             *SnowConfig // 配置
	mu            sync.Mutex  // 同步锁
	lastTimestamp int64       // 上次的时间戳
	workerId      int64       // 机器标识
	datacenterId  int64       // 数据中心标识
	sequence      int64       // 序列号

	datacenterIdWorkId int64 // 数据中心标识和机器标识的组合,此部分不变
}

// NewSnowflakeByDefaultConfig create a new Snowflake by default config.
//
// default config: epochTime: 1970-01-01 00:00:00.000 utc, sequenceBits: 12, workerIdBits: 5, datacenterIdBits: 5
func NewSnowflakeByDefaultConfig(workerId, datacenterId int64) (*Snowflake, error) {
	return defaultSnowConfig().NewSnowflake(workerId, datacenterId)
}

// NextId 生成下一个 Id。
func (s *Snowflake) NextId() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli() - s.c.epoch
	if now < s.lastTimestamp {
		return 0, ErrClockMovedBackwards
	}

	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & s.c.sequenceMask
		if s.sequence == 0 { //序列太多了，等到下一毫秒
			for now <= s.lastTimestamp {
				now = time.Now().UnixMilli() - s.c.epoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = now

	id := (now << s.c.timestampLeftShift) | s.datacenterIdWorkId | s.sequence

	return id, nil
}
