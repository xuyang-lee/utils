package snowflake

import (
	"time"
)

//snowflakeId (64bit) = 0(1bit) + timestamp(41bit) + datacenterId(5bit) + workerId(5bit) + sequence(12bit) = 1 + 41 + 5 + 5 + 12 = 64

const (
	//epoch = int64(1640966400000) // 设置起始时间 (2022-01-01 00:00:00 UTC)
	epoch = int64(0) // 设置起始时间 (1970-01-01 00:00:00 UTC)

	workerIdBits     = uint(5)                                     // 机器标识位数
	datacenterIdBits = uint(5)                                     // 数据中心标识位数
	sequenceBits     = uint(12)                                    // 序列号位数
	maxWorkerId      = int64(-1) ^ (int64(-1) << workerIdBits)     // 机器Id最大值
	maxDatacenterId  = int64(-1) ^ (int64(-1) << datacenterIdBits) // 数据中心Id最大值
	sequenceMask     = int64(-1) ^ (int64(-1) << sequenceBits)     // 序列号最大值|序列号掩码

	workerIdShift      = sequenceBits                                   // 机器Id偏移量
	datacenterIdShift  = sequenceBits + workerIdBits                    // 数据中心Id偏移量
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits // 时间戳偏移量
)

// SnowConfig 结构体定义。
type SnowConfig struct {
	epoch     int64     // 自定义的开始时间
	epochTime time.Time // 自定义的开始时间

	sequenceBits     uint // 序列号位数
	workerIdBits     uint // 机器标识位数
	datacenterIdBits uint // 数据中心标识位数

	maxWorkerId     int64 // 机器Id最大值
	maxDatacenterId int64 // 数据中心Id最大值

	workerIdShift      uint  // 机器Id向左的位移
	datacenterIdShift  uint  // 数据中心Id向左的位移
	timestampLeftShift uint  // 时间戳向左的位移
	sequenceMask       int64 // 序列号掩码
}

func defaultSnowConfig() *SnowConfig {
	return &SnowConfig{
		epoch:     epoch,
		epochTime: time.UnixMilli(epoch),

		sequenceBits:     sequenceBits,
		workerIdBits:     workerIdBits,
		datacenterIdBits: datacenterIdBits,

		maxWorkerId:     maxWorkerId,
		maxDatacenterId: maxDatacenterId,

		workerIdShift:      workerIdShift,
		datacenterIdShift:  datacenterIdShift,
		timestampLeftShift: timestampLeftShift,

		sequenceMask: sequenceMask,
	}
}

// NewSnowConfig create a new Snowflake config。
//
// required : sequenceBits + workerIdBits + datacenterIdBits = 22 && min(sequenceBits,workerIdBits,datacenterIdBits) >0
func NewSnowConfig(epochTime time.Time, datacenterIdBits, workerIdBits, sequenceBits int) *SnowConfig {
	if sequenceBits+workerIdBits+datacenterIdBits != 22 || sequenceBits < 1 || workerIdBits < 1 || datacenterIdBits < 1 {
		panic("illegal snowflake config")
	}

	return &SnowConfig{
		epoch:            epochTime.UnixMilli(),
		epochTime:        epochTime,
		sequenceBits:     uint(sequenceBits),
		workerIdBits:     uint(workerIdBits),
		datacenterIdBits: uint(datacenterIdBits),

		workerIdShift:      uint(sequenceBits),                                   // 机器Id偏移量
		datacenterIdShift:  uint(sequenceBits + workerIdBits),                    // 数据中心Id偏移量
		timestampLeftShift: uint(sequenceBits + workerIdBits + datacenterIdBits), // 时间戳偏移量

		maxWorkerId:     int64(-1) ^ (int64(-1) << workerIdBits),
		maxDatacenterId: int64(-1) ^ (int64(-1) << datacenterIdBits),
		sequenceMask:    int64(-1) ^ (int64(-1) << sequenceBits),
	}
}

// NewSnowflake create a new Snowflake by c
func (c *SnowConfig) NewSnowflake(workerId, datacenterId int64) (*Snowflake, error) {

	if workerId < 0 || workerId > (c.maxWorkerId) {
		return nil, ErrInvalidWorkerId
	}
	if datacenterId < 0 || datacenterId > (c.maxDatacenterId) {
		return nil, ErrInvalidDatacenterId
	}

	return &Snowflake{
		c:                  c,
		lastTimestamp:      -1,
		workerId:           workerId,
		datacenterId:       datacenterId,
		sequence:           0,
		datacenterIdWorkId: (datacenterId << c.datacenterIdShift) | (workerId << c.workerIdShift),
	}, nil

}
