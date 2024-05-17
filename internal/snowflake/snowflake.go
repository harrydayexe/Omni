// Package snowflake provides a very simple way to generate unique snowflake
// IDs. It is based on Twitter's snowflake algorithm. The snowflake ID is a 63
// bit integer, composed of 41 bits for time, 10 bits for node id, and 12 bits
// for a sequence number.
package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch        int64 = 1288834974657
	nodeBits           = 10
	sequenceBits       = 12
	nodeMax      int16 = -1 ^ (-1 << nodeBits)
	nodeMask           = int64(nodeMax) << sequenceBits
	sequenceMask int16 = -1 ^ (-1 << sequenceBits)
	timeShift          = nodeBits + sequenceBits
	nodeShift          = sequenceBits
)

// Snowflake is a distributed unique ID.
type Snowflake struct {
	timestamp int64
	nodeId    int64
	sequence  int64
}

func (s Snowflake) Id() int64 {
	return (s.timestamp << (sequenceBits + nodeBits)) |
		((s.nodeId << sequenceBits) & nodeMask) |
		(s.sequence & int64(sequenceMask))
}

type SnowflakeGenerator struct {
	mu        sync.Mutex
	lastStamp int64
	sequence  int16
	nodeId    int16
}

func NewSnowflakeGenerator(nodeId int16) (*SnowflakeGenerator, error) {
	if nodeId > nodeMax {
		return nil, errors.New("node id must be less than 10 bits lond")
	}
	return &SnowflakeGenerator{
		nodeId: nodeId,
	}, nil
}

// NextID generates a new snowflake ID.
func (s *SnowflakeGenerator) NextID() (Identifier, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.nextID()
}

// Get the current time in milliseconds from the Unix epoch
func (s *SnowflakeGenerator) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (s *SnowflakeGenerator) nextID() (Identifier, error) {
	now := s.getMilliSeconds() - epoch
	if now < s.lastStamp {
		return nil, errors.New("time is moving backwards,waiting until")
	}

	if now == s.lastStamp {
		s.sequence = (s.sequence + 1) & sequenceMask

		if s.sequence == 0 {
			for now <= s.lastStamp {
				now = s.getMilliSeconds() - epoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastStamp = now
	return Snowflake{
		timestamp: now,
		nodeId:    int64(s.nodeId),
		sequence:  int64(s.sequence),
	}, nil
}
