package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var year = 2022
var day = 3
var path = os.Getenv("ADVENT_HOME")

func toNumber(item rune) int {
	if 'a' <= item && item <= 'z' {
		return int((item - 'a')) + 1
	} else if 'A' <= item && item <= 'Z' {
		return int(item-'A') + 27
	}
	return -1
}

func part1(str string) int {
	length := len(str)
	compartment := length / 2
	strArr := []rune(str)
	mapA := map[rune]int{}
	mapB := map[rune]int{}

	for i := 0; i < compartment; i++ {
		_, ok := mapA[strArr[i]]
		if ok {
			mapA[strArr[i]] += 1
		} else {
			mapA[strArr[i]] = 1
		}
	}

	for i := compartment; i < len(strArr); i++ {
		_, ok := mapB[strArr[i]]
		if ok {
			mapB[strArr[i]] += 1
		} else {
			mapB[strArr[i]] = 1
		}
	}

	intersection := map[rune]bool{}
	for key, _ := range mapA {
		if _, ok := mapB[key]; ok {
			intersection[key] = true
		}
	}

	val := 0
	for key, _ := range intersection {
		// fmt.Printf("Key: %s\n", string(key))
		val = toNumber(key)
	}

	return val
}

func part2(str1 string, str2 string, str3 string) int {
	mapA := map[rune]int{}
	mapB := map[rune]int{}
	mapC := map[rune]int{}

	strArr := []rune(str1)
	for i := 0; i < len(strArr); i++ {
		_, ok := mapA[strArr[i]]
		if ok {
			mapA[strArr[i]] += 1
		} else {
			mapA[strArr[i]] = 1
		}
	}

	strArr = []rune(str2)
	for i := 0; i < len(strArr); i++ {
		_, ok := mapB[strArr[i]]
		if ok {
			mapB[strArr[i]] += 1
		} else {
			mapB[strArr[i]] = 1
		}
	}

	strArr = []rune(str3)
	for i := 0; i < len(strArr); i++ {
		_, ok := mapC[strArr[i]]
		if ok {
			mapC[strArr[i]] += 1
		} else {
			mapC[strArr[i]] = 1
		}
	}

	intersection := map[rune]bool{}
	for key, _ := range mapA {
		if _, ok := mapB[key]; ok {
			if _, okB := mapC[key]; okB {
				intersection[key] = true
			}
		}
	}

	val := 0
	for key, _ := range intersection {
		// fmt.Printf("Key: %s\n", string(key))
		val = toNumber(key)
	}

	return val
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
	data := []string{}

	sum := 0
	sum2 := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		data = append(data, line)
	}

	for i := 0; i < len(data); i += 3 {
		line1 := data[i]
		sum += part1(line1)
		fmt.Println(line1)

		line2 := data[i+1]
		sum += part1(line2)
		fmt.Println(line2)

		line3 := data[i+2]
		sum += part1(line3)
		fmt.Println(line3)

		sum2 += part2(line1, line2, line3)
	}

	fmt.Printf("p1. Sum = %d\n", sum)
	fmt.Printf("p2. Sum = %d\n", sum2)
}
