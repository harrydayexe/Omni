package snowflake

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

func TestFirstSnowflakeSequenceIsZero(t *testing.T) {
	g := NewSnowflakeGenerator(1)
	id := g.NextID()

	const mask uint64 = 0b1111_1111_1111
	sequence := id.ToInt() & mask
	if sequence != 0 {
		t.Errorf("First snowflake sequence is not zero, got %d", sequence)
	}
}

func TestSequenceIsZeroForNewTimestamp(t *testing.T) {
	g := NewSnowflakeGenerator(1)
	id1 := g.NextID()
	time.Sleep(1 * time.Second)
	id2 := g.NextID()

	const mask uint64 = 0b1111_1111_1111
	sequence1 := id1.ToInt() & mask
	sequence2 := id2.ToInt() & mask
	if sequence1 != 0 {
		t.Errorf("First snowflake sequence is not zero, got %d", sequence1)
	}
	if sequence2 != 0 {
		t.Errorf("Second snowflake sequence is not zero, got %d", sequence2)
	}
}

func TestSecondSnowflakeSequenceIsOne(t *testing.T) {
	g := NewSnowflakeGenerator(1)

	var wg sync.WaitGroup
	ch := make(chan Snowflake, 2)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			id := g.NextID()
			ch <- id
		}()
	}
	wg.Wait()

	id1 := <-ch
	id2 := <-ch

	if id1.ToInt()>>22 != id2.ToInt()>>22 {
		fmt.Println("Timestamps do not match")
		fmt.Printf("Id1 time: %d\n", id1.ToInt()>>22)
		fmt.Printf("Id2 time: %d\n", id2.ToInt()>>22)
		t.SkipNow()
	}

	const mask uint64 = 0b1111_1111_1111
	sequence1 := id1.ToInt() & mask
	sequence2 := id2.ToInt() & mask
	difference := math.Abs(float64(sequence1) - float64(sequence2))
	if difference != 1 {
		fmt.Printf("Id1: %d\n", id1.ToInt())
		fmt.Printf("Id2: %d\n", id2.ToInt())
		t.Errorf("Sequences are not 1 apart, got %d and %d, with difference %f", sequence1, sequence2, difference)
	}
}

func TestNodeIdIsCorrect(t *testing.T) {
	g := NewSnowflakeGenerator(1)
	id := g.NextID()

	const mask uint64 = 0b11_1111_1111
	nodeId := (id.ToInt() >> 12) & mask
	if nodeId != 1 {
		t.Errorf("Node ID is not 1, got %d", nodeId)
	}
}

func TestDuplicates(t *testing.T) {
	g := NewSnowflakeGenerator(1)

	const count = 10000
	var wg sync.WaitGroup
	ch := make(chan Snowflake, count)
	wg.Add(count)
	defer close(ch)
	// Concurrently count goroutines for snowFlake ID generation
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id := g.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	m := make(map[uint64]int)

	for i := 0; i < count; i++ {
		id := <-ch
		// If there is a key with id in the map, it means that the generated snowflake ID is duplicated
		_, ok := m[id.ToInt()]
		if ok {
			t.Errorf("repeat on index %d snowflake %d\n", i, id.ToInt())
			return
		}
		// store id as key in map
		m[id.ToInt()] = i
	}
}

func TestSequenceReset(t *testing.T) {
	g := NewSnowflakeGenerator(1)

	var wg sync.WaitGroup
	ch := make(chan Snowflake, 2)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			id := g.NextID()
			ch <- id
		}()
	}
	wg.Wait()

	id1 := <-ch
	id2 := <-ch

	if id1.ToInt()>>22 != id2.ToInt()>>22 {
		fmt.Println("Timestamps do not match")
		fmt.Printf("Id1 time: %d\n", id1.ToInt()>>22)
		fmt.Printf("Id2 time: %d\n", id2.ToInt()>>22)
		t.SkipNow()
	}

	const mask uint64 = 0b1111_1111_1111
	time.Sleep(1 * time.Second)
	id := g.NextID()
	sequence := id.ToInt() & mask
	if sequence != 0 {
		t.Errorf("Sequence is not zero, got %d", sequence)
	}
}

func TestParseId(t *testing.T) {
	g := NewSnowflakeGenerator(1)
	id := g.NextID()

	parsedId := ParseId(id.ToInt())
	if parsedId != id {
		t.Errorf("Parsed ID is not equal to original ID, got %v; expected %v", parsedId, id)
	}
}

func TestNodeMaxError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewSnowflakeGenerator(math.MaxUint16)
}

// The aim of this benchmark is to see how quickly we can generate new ids
// This is because there are only so many new IDs available in a single millisecond
func BenchmarkNextId(b *testing.B) {
	g := NewSnowflakeGenerator(1)
	for i := 0; i < b.N; i++ {
		g.NextID()
	}
}
