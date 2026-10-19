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

	// Ask for inputs
	fmt.Print("\n ========= CONFIGURATION ========= \n\n")

	fmt.Print(">> Number of maximum instances: ")
	fmt.Scan(&maxInstances)

	fmt.Print(">> Number of tanks: ")
	fmt.Scan(&tanks)

	fmt.Print(">> Number of healers: ")
	fmt.Scan(&healers)

	fmt.Print(">> Number of DPS: ")
	fmt.Scan(&dps)

	fmt.Print(">> Minimum time of instance: ")
	fmt.Scan(&minTime)

	fmt.Print(">> Maximum time of instance: ")
	fmt.Scan(&maxTime)

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

	fmt.Print("\n ========== !! START !! ========== \n\n")

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

	PrintDungeonSummary(dungeons)
}
