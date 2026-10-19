package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func ProcessParty(
	party Party,
	dungeons []Dungeon,
	semaphore chan struct{},
	mutex *sync.Mutex,
	t1 int,
	t2 int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	defer func() {
		<-semaphore
	}()

	mutex.Lock()
	var dungeon *Dungeon
	for i := range dungeons {
		if dungeons[i].status == "empty" {
			dungeon = &dungeons[i]
			dungeon.status = "active"
			break
		}
	}

	duration := t1 + rand.IntN(t2-t1+1)

	fmt.Printf(
		"\n>> Dungeon %d: active | Party %d entered | Duration: %d seconds\n",
		dungeon.id,
		party.id,
		duration,
	)

	PrintDungeonStatus(dungeons)

	mutex.Unlock()

	// Simulate dungeon
	time.Sleep(time.Duration(duration) * time.Second)

	// Party finishes
	mutex.Lock()

	dungeon.partiesServed++
	dungeon.totalTime += duration
	dungeon.status = "empty"

	fmt.Printf(
		"\n>> Dungeon %d: empty | Party %d finished\n",
		dungeon.id,
		party.id,
	)

	PrintDungeonStatus(dungeons)

	mutex.Unlock()
}
