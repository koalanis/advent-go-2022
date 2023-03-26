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
var day = 14
var path = os.Getenv("ADVENT_HOME")

type Tuple struct {
	x, y int
}

func (t Tuple) String() string {
	return fmt.Sprintf("(%d,%d)", t.x, t.y)
}

func TupleFromStr(str string) Tuple {
	tokens := strings.Split(str, ",")
	valX, _ := strconv.Atoi(tokens[0])
	valY, _ := strconv.Atoi(tokens[1])
	return Tuple{valX, valY}
}

func getLineBetweenPaths(t1 Tuple, t2 Tuple) []Tuple {
	path := []Tuple{}

	if t1.x == t2.x {
		sy := t1.y
		ey := t2.y
		if sy < ey {
			for i := sy; i <= ey; i += 1 {
				path = append(path, Tuple{t1.x, i})
			}

		} else {
			for i := sy; i >= ey; i -= 1 {
				path = append(path, Tuple{t1.x, i})
			}
		}

	} else {
		sx := t1.x
		ex := t2.x
		if sx < ex {
			for i := sx; i <= ex; i += 1 {
				path = append(path, Tuple{i, t1.y})
			}

		} else {
			for i := sx; i >= ex; i -= 1 {
				path = append(path, Tuple{i, t1.y})
			}
		}
	}

	return path
}

func getPath(dataLine string) []Tuple {
	path := []Tuple{}
	DELIMETER := " -> "
	tokens := strings.Split(dataLine, DELIMETER)
	segments := []Tuple{}
	for _, token := range tokens {
		tuple := TupleFromStr(token)
		segments = append(segments, tuple)
	}

	for i := 0; i < len(segments)-1; i++ {
		left, right := segments[i], segments[i+1]
		straight := getLineBetweenPaths(left, right)
		for _, t := range straight {
			path = append(path, t)
		}
	}

	return path
}

func constructGridMap(data []string) map[Tuple]string {
	grid := make(map[Tuple]string)

	for _, dataLine := range data {
		for _, tuple := range getPath(dataLine) {
			grid[tuple] = "#"
		}
	}

	return grid
}

func getDepthBounds(m map[Tuple]string) int {
	maxY := 0
	for key, _ := range m {

		if maxY < key.y {
			maxY = key.y
		}
	}
	return maxY
}

func FallDown(sand Tuple, grid map[Tuple]string) Tuple {
	belowL, belowM, belowR := Tuple{sand.x - 1, sand.y + 1}, Tuple{sand.x, sand.y + 1}, Tuple{sand.x + 1, sand.y + 1}
	_, okL := grid[belowL]
	_, okM := grid[belowM]

	if !okM {
		// fmt.Printf("current = %s, next = %s\n", sand, belowM)
		return belowM
	}

	if !okL {
		// fmt.Printf("current = %s, next = %s\n", sand, belowL)
		return belowL
	}
	// fmt.Printf("current = %s, next = %s\n", sand, belowR)
	return belowR
}

func CanFallDown(sand Tuple, grid map[Tuple]string, maxY int) (bool, bool) {

	nextY := sand.y + 1
	if nextY > maxY {
		return true, true
	}
	belowL, belowM, belowR := Tuple{sand.x - 1, nextY}, Tuple{sand.x, nextY}, Tuple{sand.x + 1, nextY}
	_, okL := grid[belowL]
	_, okM := grid[belowM]
	_, okR := grid[belowR]

	return !okL || !okM || !okR, false
}

func PrintGrid(grid map[Tuple]string) {
	maxX, maxY := 0, 0

	for key, _ := range grid {
		if maxX < key.x {
			maxX = key.x
		}
		if maxY < key.y {
			maxY = key.y
		}
	}

	minX, minY := maxX, maxY

	for key, _ := range grid {
		if minX > key.x {
			minX = key.x
		}
		if minY > key.y {
			minY = key.y
		}
	}
	out := ""
	for j := minY; j <= maxY; j++ {
		line := ""
		for i := minX; i <= maxX; i++ {
			tuple := Tuple{i, j}
			val, exists := grid[tuple]

			if !exists {
				line += "."
			} else {
				line += val
			}
		}
		out += line + "\n"
	}

	fmt.Print(out)

}

func simulateP1(grid map[Tuple]string, source Tuple, maxDepth int) int {
	end := false
	for !end {

		falling := true
		sand := Tuple{source.x, source.y}

		for falling {
			if canFall, outOfBounds := CanFallDown(sand, grid, maxDepth); !outOfBounds {
				if canFall {
					sand = FallDown(sand, grid)
				} else {
					grid[sand] = "o"
					falling = false
				}
			} else {
				falling = false
				end = true
			}
		}
	}

	rockCount := 0
	for _, val := range grid {
		if val == "o" {
			rockCount += 1
		}
	}

	// fmt.Println()
	// PrintGrid(grid)
	// fmt.Println()

	return rockCount
}

func CanFallDownWithFloor(sand Tuple, grid map[Tuple]string, floor int) bool {
	nextY := sand.y + 1

	if nextY == floor {
		return false
	}

	belowL, belowM, belowR := Tuple{sand.x - 1, nextY}, Tuple{sand.x, nextY}, Tuple{sand.x + 1, nextY}
	_, okL := grid[belowL]
	_, okM := grid[belowM]
	_, okR := grid[belowR]

	return !okL || !okM || !okR
}

func simulateP2(grid map[Tuple]string, source Tuple, maxDepth int) int {
	end := false

	floorY := maxDepth + 2

	for !end {

		falling := true
		sand := Tuple{source.x, source.y}

		_, sourceBlocked := grid[sand]
		if sourceBlocked {
			end = true
			break
		}

		for falling {
			if canFall := CanFallDownWithFloor(sand, grid, floorY); canFall {
				sand = FallDown(sand, grid)
			} else {
				grid[sand] = "o"
				falling = false
			}
		}

	}

	rockCount := 0
	for _, val := range grid {
		if val == "o" {
			rockCount += 1
		}
	}

	// fmt.Println()
	// PrintGrid(grid)
	// fmt.Println()

	return rockCount
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

	for scanner.Scan() {
		line := scanner.Text()
		data = append(data, line)
	}

	gridMap := constructGridMap(data)

	maxY := getDepthBounds(gridMap)
	fmt.Printf("maxY = %d\n", maxY)

	sandCount := simulateP1(gridMap, Tuple{500, 0}, maxY)
	fmt.Printf("sandCount = %d\n", sandCount)

	sandCount = simulateP2(gridMap, Tuple{500, 0}, maxY)
	fmt.Printf("sandCount = %d\n", sandCount)
}
