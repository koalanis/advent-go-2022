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
var day = 4
var path = os.Getenv("ADVENT_HOME")

type interval struct {
	left, right int
}

func (i interval) String() string {
	return fmt.Sprintf("(%d, %d)", i.left, i.right)
}

func fullyOverlap(a interval, b interval) bool {
	return (a.left <= b.left && b.right <= a.right) || (b.left <= a.left && a.right <= b.right)
}

func partiallyOverlap(a interval, b interval) bool {
	return (a.left <= b.left && b.left <= a.right) || (b.left <= a.left && a.left <= b.right)
}

func strToInterval(str string) interval {
	tokens := strings.Split(str, "-")
	a, _ := strconv.Atoi(tokens[0])
	b, _ := strconv.Atoi(tokens[1])
	return interval{a, b}
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

	count := 0
	countTwo := 0

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")
		a, b := strToInterval(tokens[0]), strToInterval(tokens[1])
		if fullyOverlap(a, b) {
			count += 1
		}

		if partiallyOverlap(a, b) {
			countTwo += 1
		}
	}

	fmt.Printf("p1. count = %d\n", count)
	fmt.Printf("p2. count = %d\n", countTwo)

}
