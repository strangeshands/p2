package main

import (
	"fmt"
	"sync"
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

	var err error

	maxInstances, tanks, healers, dps, minTime, maxTime, err = LoadConfig("config.txt")

	if err != nil {
		fmt.Printf(">> Error loading config: %v\n", err)
		return
	}

	// Ask for inputs
	fmt.Print("\n ========= CONFIGURATION ========= \n\n")

	fmt.Printf(">> Maximum instances: %d\n", maxInstances)
	fmt.Printf(">> Tanks: %d\n", tanks)
	fmt.Printf(">> Healers: %d\n", healers)
	fmt.Printf(">> DPS: %d\n", dps)
	fmt.Printf(">> Minimum instance time: %d seconds\n", minTime)
	fmt.Printf(">> Maximum instance time: %d seconds\n", maxTime)

	fmt.Print("\n ================================= \n\n")

	if !ValidateInput(maxInstances, tanks, healers, dps, minTime, maxTime) {
		fmt.Print(">> Invalid inputs. Please check.\n\n")
		return
	}

	if !ValidateParty(tanks, healers, dps) {
		fmt.Print("\n>> Roles are not sufficient to make a party.\n\n")
		return
	}

	// Calculate parties
	numParty := CalculateParties(tanks, healers, dps)
	fmt.Printf(">> Possible Parties: %d\n", numParty)

	// Calculate remaining
	remT, remH, remD := CalculateRemaining(tanks, healers, dps, numParty)
	fmt.Printf(">> Remainings tank, healer, and DPS: %d, %d, %d\n", remT, remH, remD)

	fmt.Print("\n ================================= \n\n")

	fmt.Printf(">> Creating parties...\n")

	// Assign party
	for i := 1; i <= numParty; i++ {
		p := Party{
			id: i,
		}
		parties = append(parties, p)
	}
	PrintParties(parties)

	fmt.Printf("\n>> Creating dungeons...\n")

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
	PrintDungeons(dungeons)

	fmt.Print("\n ================================= \n")

	fmt.Print("\n ========== !! START !! ========== \n")

	// Workers
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxInstances)
	var mutex sync.Mutex

	wg.Add(numParty)
	for i := 0; i < numParty; i++ {
		semaphore <- struct{}{}
		go ProcessParty(parties[i], dungeons, semaphore, &mutex, minTime, maxTime, &wg)
	}
	wg.Wait()

	fmt.Print("\n ================================= \n")

	PrintDungeonSummary(dungeons)
}
