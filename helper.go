package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadConfig(filename string) (int, int, int, int, int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	defer file.Close()

	var maxInstances int
	var tanks int
	var healers int
	var dps int
	var minTime int
	var maxTime int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		num, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, err
		}

		switch key {
		case "n":
			maxInstances = num
		case "t":
			tanks = num
		case "h":
			healers = num
		case "d":
			dps = num
		case "t1":
			minTime = num
		case "t2":
			maxTime = num
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	return maxInstances, tanks, healers, dps, minTime, maxTime, nil
}

func ValidateInput(n int, t int, h int, d int, t1 int, t2 int) bool {
	if n <= 0 || t < 0 || h < 0 || d < 0 || t1 < 0 || t2 < 0 || t1 > t2 || t2 > MAX_DURATION {
		fmt.Println("--- RULES ---")
		fmt.Println(">> Maximum instances must be greater than 0.")
		fmt.Println(">> Tanks, healers, and DPS must be nonnegative.")
		fmt.Println(">> Minimum and maximum time must be nonnegative.")
		fmt.Println(">> Minimum time must be less than maximum time.")
		fmt.Printf(">> Maximum time must not exceed %d seconds.\n", MAX_DURATION)
		fmt.Println("-------------\n")

		return false
	}

	return true
}

func ValidateParty(t int, h int, d int) bool {
	flag := true

	if t < 1 {
		flag = false
		fmt.Print(">> Not enough tanks.\n")
	}
	if h < 1 {
		flag = false
		fmt.Print(">> Not enough healers.\n")
	}
	if d < 3 {
		flag = false
		fmt.Print(">> Not enough DPS.\n")
	}

	return flag
}

func CalculateParties(t int, h int, d int) int {
	p := 0
	dpsGroup := d / 3

	p = min(t, h, dpsGroup)

	return p
}

func PrintParties(parties []Party, t int, h int, d int) {
	fmt.Println("\n--- Parties ---")

	if len(parties) == 0 {
		fmt.Println("No parties formed.")
		return
	}

	for _, party := range parties {
		fmt.Printf("Party %d -- T: %d | H: %d | DPS: %d\n", party.id, party.tank, party.healer, party.dps)
	}

	fmt.Printf("-- There are %d tank/s, %d healer/s and %d DPS remaining (no party formed).\n", t, h, d)
}

func PrintDungeons(dungeons []Dungeon) {
	fmt.Println("\n--- Dungeon Instances ---")

	for _, dungeon := range dungeons {
		fmt.Printf(
			"Dungeon %d | Status: %s | Parties Served: %d | Total Time: %d seconds\n",
			dungeon.id,
			dungeon.status,
			dungeon.partiesServed,
			dungeon.totalTime,
		)
	}
}

func PrintDungeonStatus(d []Dungeon) {
	fmt.Println("\n--- Current Dungeon Status ---")

	for _, dungeon := range d {
		fmt.Printf("Dungeon %d: %s\n", dungeon.id, dungeon.status)
	}
}

func PrintDungeonSummary(d []Dungeon) {
	fmt.Println("\n--- Dungeon Summary ---")

	for _, dungeon := range d {
		fmt.Printf("Dungeon %d\n", dungeon.id)
		fmt.Printf("Parties Served: %d\n", dungeon.partiesServed)
		fmt.Printf("Total Time Served: %d seconds\n\n", dungeon.totalTime)
	}
}
