package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var year = 2022
var day = 10
var path = os.Getenv("ADVENT_HOME")

func main() {

	dataPath := fmt.Sprintf("%s/%d/data/day%d/data.txt", path, year, day)
	fmt.Println(dataPath)

	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	regX := 1
	cycle := 0

	valMap := make(map[int]int)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "noop") {
			cycle += 1
			valMap[cycle] = regX
		} else if strings.HasPrefix(line, "addx") {
			tokens := strings.Split(line, " ")
			val, _ := strconv.Atoi(tokens[1])
			cycle += 1
			valMap[cycle] = regX
			cycle += 1
			valMap[cycle] = regX
			regX += val
		}
	}

	sumSignalStrength := 0
	cycleAt := []int{20, 60, 100, 140, 180, 220}

	for _, val := range cycleAt {
		sumSignalStrength += val * valMap[val]
	}

	fmt.Printf("sum of signal strengths = %d\n", sumSignalStrength)

	spritePos := 1

	RedColor := "\033[1;31m%s\033[0m"

	for i := 0; i < 240; i++ {
		if (i)%40 == 0 {
			fmt.Printf("\n")
		}
		tick := ((i % 40) + 1)
		spritePos = valMap[i+1]
		if spritePos <= tick && tick <= spritePos+2 {
			fmt.Printf(RedColor, "#")
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Println()

}
