// Package snowflake provides a very simple way to generate unique snowflake
// IDs. It is based on Twitter's snowflake algorithm. The snowflake ID is a 63
// bit integer, composed of 41 bits for time, 10 bits for node id, and 12 bits
// for a sequence number.
package snowflake

import (
	"sync"
	"time"
)

const (
	epoch        int64  = 1704067200000
	nodeBits     uint64 = 10
	sequenceBits uint64 = 12
	nodeMax      uint16 = -1 ^ (-1 << nodeBits)
	nodeMask     uint64 = uint64(nodeMax) << sequenceBits
	sequenceMask uint16 = -1 ^ (-1 << sequenceBits)
	timeShift    uint64 = nodeBits + sequenceBits
	nodeShift    uint64 = sequenceBits
)

// Snowflake is a distributed unique ID.
//
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
type Snowflake struct {
	timestamp uint64
	nodeId    uint16
	sequence  uint16
}

func (s Snowflake) Id() Snowflake {
	return s
}

func ParseId(id uint64) Snowflake {
	return Snowflake{
		timestamp: id >> timeShift,
		nodeId:    uint16((id & nodeMask) >> sequenceBits),
		sequence:  uint16(id & uint64(sequenceMask)),
	}
}

func (s Snowflake) ToInt() uint64 {
	return uint64((s.timestamp << (timeShift)) |
		((uint64(s.nodeId) << nodeShift) & nodeMask) |
		(uint64(s.sequence) & uint64(sequenceMask)))
}

type SnowflakeGenerator struct {
	mu        sync.Mutex
	lastStamp uint64
	sequence  uint16
	nodeId    uint16
}

func NewSnowflakeGenerator(nodeId uint16) *SnowflakeGenerator {
	if nodeId > nodeMax {
		panic("node id must be less than 10 bits long")
	}
	return &SnowflakeGenerator{
		nodeId: nodeId,
	}
}

// NextID generates a new snowflake ID.
func (s *SnowflakeGenerator) NextID() Snowflake {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.nextID()
}

// Get the current time in milliseconds from the Unix epoch
func (s *SnowflakeGenerator) getMilliSeconds() int64 {
	return time.Now().UnixMilli()
}

func (s *SnowflakeGenerator) nextID() Snowflake {
	now := uint64(s.getMilliSeconds() - epoch)
	if now < s.lastStamp {
		panic("time is moving backwards, current time is before the last time this method was called")
	}

	if now == s.lastStamp {
		s.sequence = (s.sequence + 1) & sequenceMask

		if s.sequence == 0 {
			for now <= s.lastStamp {
				now = uint64(s.getMilliSeconds() - epoch)
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastStamp = now
	return Snowflake{
		timestamp: now,
		nodeId:    s.nodeId,
		sequence:  s.sequence,
	}
}
