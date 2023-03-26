package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var year = 2022
var day = 8
var path = os.Getenv("ADVENT_HOME")

func checkTreeVisibleFromOutside(grid [][]int, r int, c int) bool {
	rLen := len(grid)
	cLen := 0
	if rLen > 0 {
		cLen = len(grid[0])
	}

	canBeSeen := true
	valueAtRC := grid[r][c]

	// check up
	for rr := r - 1; rr >= 0; rr-- {
		valueAt := grid[rr][c]
		if valueAt >= valueAtRC {
			canBeSeen = false
			break
		}
	}

	if canBeSeen {
		return canBeSeen
	}
	canBeSeen = true
	for rr := r + 1; rr < rLen; rr++ {
		valueAt := grid[rr][c]
		if valueAt >= valueAtRC {
			canBeSeen = false
			break
		}
	}

	if canBeSeen {
		return canBeSeen
	}
	canBeSeen = true

	for cc := c - 1; cc >= 0; cc-- {
		valueAt := grid[r][cc]
		if valueAt >= valueAtRC {
			canBeSeen = false
			break
		}
	}

	if canBeSeen {
		return canBeSeen
	}
	canBeSeen = true

	for cc := c + 1; cc < cLen; cc++ {
		valueAt := grid[r][cc]
		if valueAt >= valueAtRC {
			canBeSeen = false
			break
		}
	}

	return canBeSeen
}

func checkTreeScenicView(grid [][]int, r int, c int) int {

	scenicViewScore := 1
	rLen := len(grid)
	cLen := 0
	if rLen > 0 {
		cLen = len(grid[0])
	}

	valueAtRC := grid[r][c]

	// check up
	view := 0
	for rr := r - 1; rr >= 0; rr-- {
		valueAt := grid[rr][c]
		view += 1
		if valueAt >= valueAtRC {
			break
		}
	}
	scenicViewScore *= view

	view = 0
	for rr := r + 1; rr < rLen; rr++ {
		view += 1
		valueAt := grid[rr][c]
		if valueAt >= valueAtRC {
			break
		}
	}
	scenicViewScore *= view

	view = 0
	for cc := c - 1; cc >= 0; cc-- {
		valueAt := grid[r][cc]
		view += 1
		if valueAt >= valueAtRC {
			break
		}
	}
	scenicViewScore *= view

	view = 0
	for cc := c + 1; cc < cLen; cc++ {
		view += 1
		valueAt := grid[r][cc]
		if valueAt >= valueAtRC {
			break
		}
	}

	scenicViewScore *= view
	return scenicViewScore
}

func partOne(grid [][]int) int {
	rLen := len(grid)
	cLen := 0
	if rLen > 0 {
		cLen = len(grid[0])
	}

	outsidePerimeter := (2 * (rLen - 2)) + (2 * cLen)

	insideCount := 0
	for i := 1; i < rLen-1; i++ {
		for j := 1; j < cLen-1; j++ {
			if checkTreeVisibleFromOutside(grid, i, j) {
				insideCount += 1
			}
		}
	}

	return outsidePerimeter + insideCount
}

func partTwo(grid [][]int) int {
	rLen := len(grid)
	cLen := 0
	if rLen > 0 {
		cLen = len(grid[0])
	}

	max := 0
	for i := 1; i < rLen-1; i++ {
		for j := 1; j < cLen-1; j++ {
			sv := checkTreeScenicView(grid, i, j)
			if max < sv {
				max = sv
			}
		}
	}

	return max
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

	data := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		row := []int{}
		for _, v := range line {
			i, _ := strconv.Atoi(string(v))
			row = append(row, i)
		}
		data = append(data, row)
	}
	count := partOne(data)
	fmt.Printf("count = %d\n", count)

	max := partTwo(data)
	fmt.Printf("max = %d\n", max)

}
