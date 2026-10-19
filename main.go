package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Party struct {
	id int
}

type Dungeon struct {
	id            int
	status        string
	partiesServed int
	totalTime     int
}

func validateInput(n int, t int, h int, d int, t1 int, t2 int) bool {
	if n <= 0 || t < 0 || h < 0 || d < 0 || t1 <= 0 || t2 > 15 || t1 > t2 {
		return false
	}

	return true
}

func calculateParties(t int, h int, d int) int {
	p := 0
	dpsGroup := d / 3

	p = min(t, h, dpsGroup)

	return p
}

func calculateRemaining(t int, h int, d int, p int) (int, int, int) {
	remTanks := t - p
	remHeal := h - p
	remDPS := d - (p * 3)

	return remTanks, remHeal, remDPS
}

func processParty(
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
		"Dungeon %d: active | Party %d entered | Duration: %d seconds\n",
		dungeon.id,
		party.id,
		duration,
	)

	printDungeonStatus(dungeons)

	mutex.Unlock()

	// Simulate dungeon
	time.Sleep(time.Duration(duration) * time.Second)

	// Party finishes
	mutex.Lock()

	dungeon.partiesServed++
	dungeon.totalTime += duration
	dungeon.status = "empty"

	fmt.Printf(
		"Dungeon %d: empty | Party %d finished\n",
		dungeon.id,
		party.id,
	)

	printDungeonStatus(dungeons)

	mutex.Unlock()
}

func printDungeonStatus(d []Dungeon) {
	fmt.Println("\n--- Current Dungeon Status ---")

	for _, dungeon := range d {
		fmt.Printf("Dungeon %d: %s\n", dungeon.id, dungeon.status)
	}
}

func printDungeonSummary(d []Dungeon) {
	fmt.Println("\n--- Dungeon Summary ---")

	for _, dungeon := range d {
		fmt.Printf("Dungeon %d\n", dungeon.id)
		fmt.Printf("Parties Served: %d\n", dungeon.partiesServed)
		fmt.Printf("Total Time Served: %d seconds\n\n", dungeon.totalTime)
	}
}

func main() {
	// Variables
	var maxInstances int
	var tanks int
	var healers int
	var dps int
	var minTime int
	var maxTime int

	var parties []Party
	var dungeons []Dungeon

	// Ask for inputs
	fmt.Scan(&maxInstances)
	fmt.Scan(&tanks)
	fmt.Scan(&healers)
	fmt.Scan(&dps)
	fmt.Scan(&minTime)
	fmt.Scan(&maxTime)

	if !validateInput(maxInstances, tanks, healers, dps, minTime, maxTime) {
		fmt.Println("Invalid inputs. Please check.")
		return
	}

	// Print values back
	fmt.Printf("Instances %d\n", maxInstances)
	fmt.Printf("Tanks %d\n", tanks)
	fmt.Printf("Healers %d\n", healers)
	fmt.Printf("DPS %d\n", dps)
	fmt.Printf("Minimum Time %d\n", minTime)
	fmt.Printf("Maximum Time %d\n", maxTime)

	// Calculate parties
	numParty := calculateParties(tanks, healers, dps)
	fmt.Printf("Number of Possible Parties %d\n", numParty)

	// Calculate remaining
	remT, remH, remD := calculateRemaining(tanks, healers, dps, numParty)
	fmt.Printf("Remainings T H DPS %d %d %d\n", remT, remH, remD)

	// Assign prty
	for i := 1; i <= numParty; i++ {
		p := Party{
			id: i,
		}
		parties = append(parties, p)
	}
	fmt.Println(parties)

	// Assign dungeons
	for i := 1; i <= maxInstances; i++ {
		d := Dungeon{
			id:            i,
			status:        "empty",
			partiesServed: 0,
			totalTime:     0,
		}

		dungeons = append(dungeons, d)
	}
	fmt.Println(dungeons)

	// Workers
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxInstances)
	var mutex sync.Mutex

	wg.Add(numParty)
	for i := 0; i < numParty; i++ {
		semaphore <- struct{}{}
		go processParty(parties[i], dungeons, semaphore, &mutex, minTime, maxTime, &wg)
	}
	wg.Wait()

	printDungeonSummary(dungeons)
}
