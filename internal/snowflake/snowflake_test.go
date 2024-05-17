package snowflake

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

func TestFirstSnowflakeSequenceIsZero(t *testing.T) {
	g, _ := NewSnowflakeGenerator(1)
	id, _ := g.NextID()

	const mask int64 = 0b1111_1111_1111
	sequence := id.Id() & mask
	if sequence != 0 {
		t.Errorf("First snowflake sequence is not zero, got %d", sequence)
	}
}

func TestSequenceIsZeroForNewTimestamp(t *testing.T) {
	g, _ := NewSnowflakeGenerator(1)
	id1, _ := g.NextID()
	time.Sleep(1 * time.Second)
	id2, _ := g.NextID()

	const mask int64 = 0b1111_1111_1111
	sequence1 := id1.Id() & mask
	sequence2 := id2.Id() & mask
	if sequence1 != 0 {
		t.Errorf("First snowflake sequence is not zero, got %d", sequence1)
	}
	if sequence2 != 0 {
		t.Errorf("Second snowflake sequence is not zero, got %d", sequence2)
	}
}

func TestSecondSnowflakeSequenceIsOne(t *testing.T) {
	g, _ := NewSnowflakeGenerator(1)

	var wg sync.WaitGroup
	ch := make(chan Identifier, 2)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			id, _ := g.NextID()
			ch <- id
		}()
	}
	wg.Wait()

	id1 := <-ch
	id2 := <-ch

	if id1.Id()>>22 != id2.Id()>>22 {
		fmt.Println("Timestamps do not match")
		fmt.Printf("Id1 time: %d\n", id1.Id()>>22)
		fmt.Printf("Id2 time: %d\n", id2.Id()>>22)
		t.SkipNow()
	}

	const mask int64 = 0b1111_1111_1111
	sequence1 := id1.Id() & mask
	sequence2 := id2.Id() & mask
	if math.Abs(float64(sequence1-sequence2)) != 1 {
		fmt.Printf("Id1: %d\n", id1.Id())
		fmt.Printf("Id2: %d\n", id2.Id())
		t.Errorf("Sequences are not 1 apart, got %d and %d", sequence1, sequence2)
	}
}

func TestNodeIdIsCorrect(t *testing.T) {
	g, _ := NewSnowflakeGenerator(1)
	id, _ := g.NextID()

	const mask int64 = 0b11_1111_1111
	nodeId := (id.Id() >> 12) & mask
	if nodeId != 1 {
		t.Errorf("Node ID is not 1, got %d", nodeId)
	}
}

func TestDuplicates(t *testing.T) {
	g, _ := NewSnowflakeGenerator(1)

	const count = 10000
	var wg sync.WaitGroup
	ch := make(chan Identifier, count)
	wg.Add(count)
	defer close(ch)
	// Concurrently count goroutines for snowFlake ID generation
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id, _ := g.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	m := make(map[int64]int)

	for i := 0; i < count; i++ {
		id := <-ch
		// If there is a key with id in the map, it means that the generated snowflake ID is duplicated
		_, ok := m[id.Id()]
		if ok {
			t.Errorf("repeat on index %d snowflake %s\n", i, id)
			return
		}
		// store id as key in map
		m[id.Id()] = i
	}
}

func TestSequenceReset(t *testing.T) {
	g, _ := NewSnowflakeGenerator(1)

	var wg sync.WaitGroup
	ch := make(chan Identifier, 2)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			id, _ := g.NextID()
			ch <- id
		}()
	}
	wg.Wait()

	id1 := <-ch
	id2 := <-ch

	if id1.Id()>>22 != id2.Id()>>22 {
		fmt.Println("Timestamps do not match")
		fmt.Printf("Id1 time: %d\n", id1.Id()>>22)
		fmt.Printf("Id2 time: %d\n", id2.Id()>>22)
		t.SkipNow()
	}

	const mask int64 = 0b1111_1111_1111
	time.Sleep(1 * time.Second)
	id, _ := g.NextID()
	sequence := id.Id() & mask
	if sequence != 0 {
		t.Errorf("Sequence is not zero, got %d", sequence)
	}
}
