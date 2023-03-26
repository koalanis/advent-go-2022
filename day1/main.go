package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var year = 2022
var day = 1
var path = os.Getenv("ADVENT_HOME")

func main() {

	dataPath := fmt.Sprintf("%s/%d/data/day%d/data.txt", path, year, day)
	fmt.Println(dataPath)

	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var calorieList []int

	scanner := bufio.NewScanner(file)

	acc := 0
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)

		if len(text) > 0 {
			i, err := strconv.Atoi(text)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			acc += i
		} else {
			calorieList = append(calorieList, acc)
			acc = 0
		}
	}

	sort.Slice(calorieList, func(i, j int) bool {
		return calorieList[i] > calorieList[j]
	})

	maxCalories := calorieList[0]

	biggestThree := calorieList[0] + calorieList[1] + calorieList[2]

	fmt.Println("solution p1 =", maxCalories)
	fmt.Println("solution p2 =", biggestThree)

}
