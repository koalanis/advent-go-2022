package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var year = 2022
var day = 2
var path = os.Getenv("ADVENT_HOME")

const (
	_ = iota
	ROCK
	PAPER
	SCISSORS
)

const (
	_ = iota
	LOSS
	DRAW
	WIN
)

var rpsMap = map[string]int{
	"A": ROCK,
	"B": PAPER,
	"C": SCISSORS,
	"X": ROCK,
	"Y": PAPER,
	"Z": SCISSORS,
}

var rpsValueMap = map[int]int{
	ROCK:     1,
	PAPER:    2,
	SCISSORS: 3,
}

var rpsScoreValueMap = map[int]int{
	LOSS: 0,
	DRAW: 3,
	WIN:  6,
}

func doRockPaperScissors(cpu int, player int) int {
	if cpu == player {
		return DRAW
	} else {
		if (cpu == ROCK && player == SCISSORS) || (cpu == PAPER && player == ROCK) || (cpu == SCISSORS && player == PAPER) {
			return LOSS
		}
		return WIN
	}
}

func rockPaperScissorsRound(cpu string, player string) int {
	cpuEnum := rpsMap[cpu]
	playerEnum := rpsMap[player]
	status := doRockPaperScissors(cpuEnum, playerEnum)
	return rpsScoreValueMap[status] + rpsValueMap[playerEnum]
}

// -----------------------------------------------------------------------------------------------------------------------------

var scoreValueMap = map[string]int{
	"X": LOSS,
	"Y": DRAW,
	"Z": WIN,
}

func doRockPaperScissorsV2(cpu int, player int) int {
	if player == DRAW {
		return cpu
	} else {
		if player == LOSS {
			if cpu == ROCK {
				return SCISSORS
			} else if cpu == PAPER {
				return ROCK
			} else {
				return PAPER
			}
		} else {
			if cpu == ROCK {
				return PAPER
			} else if cpu == PAPER {
				return SCISSORS
			} else {
				return ROCK
			}
		}
	}
}

func rockPaperScissorsRoundV2(cpu string, player string) int {
	cpuEnum := rpsMap[cpu]
	scoreEnum := scoreValueMap[player]
	playerEnum := doRockPaperScissorsV2(cpuEnum, scoreEnum)
	return rpsScoreValueMap[scoreEnum] + rpsValueMap[playerEnum]
}

func main() {

	dataPath := fmt.Sprintf("%s/%d/data/day%d/data.txt", path, year, day)
	fmt.Println(dataPath)

	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// fmt.Printf("Hello %s\n", rpsMap["A"])

	score := 0
	scoreTwo := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tokens := strings.Split(line, " ")
		score += rockPaperScissorsRound(tokens[0], tokens[1])
		scoreTwo += rockPaperScissorsRoundV2(tokens[0], tokens[1])
	}

	fmt.Printf("p1. Score=%d\n", score)
	fmt.Printf("p2. Score=%d\n", scoreTwo)

}
