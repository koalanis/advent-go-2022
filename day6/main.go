package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var year = 2022
var day = 6
var path = os.Getenv("ADVENT_HOME")

func helper(data string, size int) {

	chars := []rune(data)
	widthMap := map[string]int{}

	for i := 0; i < size; i++ {
		char := string(chars[i])
		if _, ok := widthMap[char]; ok {
			widthMap[char] += 1
		} else {
			widthMap[char] = 1
		}
	}

	fmt.Println(len(widthMap))
	if len(widthMap) == size {
		fmt.Printf("%d\n", size)
		return
	}

	i := size
	for ; i < len(chars); i++ {
		char := string(chars[i])
		toRemove := string(chars[i-size])
		if _, ok := widthMap[toRemove]; ok {
			widthMap[toRemove] -= 1
			if widthMap[toRemove] <= 0 {
				delete(widthMap, toRemove)
			}
		}

		if _, ok := widthMap[char]; ok {
			widthMap[char] += 1
		} else {
			widthMap[char] = 1
		}

		if len(widthMap) == size {
			fmt.Printf("%d\n", i+1)
			return
		}

	}
}

func partOne(data string) {
	helper(data, 4)
}

func partTwo(data string) {
	helper(data, 14)
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

	var data string
	for scanner.Scan() {
		data = scanner.Text()
	}

	partOne(data)
	partTwo(data)
}
