package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
)

var year = 2022
var day = 12
var path = os.Getenv("ADVENT_HOME")

func findInGrid(grid [][]rune, find rune) (int, int) {
	i, j := 0, 0

	for r, row := range grid {
		for c, col := range row {
			if col == find {
				i, j = r, c
				break
			}
		}
	}

	return i, j
}

func contains(s []Pair, e Pair) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func inGrid(grid [][]rune, at Pair) bool {
	x, y := at.x, at.y
	return 0 <= x && x < len(grid) && 0 <= y && y < len(grid[x])
}

func canClimb(grid [][]rune, at Pair, to Pair) bool {
	if inGrid(grid, at) && inGrid(grid, to) {
		atRune := grid[at.x][at.y]
		toRune := grid[to.x][to.y]

		cond := toRune-atRune <= 1
		// fmt.Printf("atRune=%s toRune=%s cond=%t\n", string(atRune), string(toRune), cond)

		return cond
	}
	return false
}

func dfs(grid [][]rune, at Pair, end Pair, trail []Pair, count *int) {
	if inGrid(grid, at) && !contains(trail, at) {

		if at == end {
			fmt.Printf("found path %d count=%d\n", len(trail), *count)

			if len(trail) < *count {
				fmt.Printf("found shorter path %d\n", len(trail))
				*count = len(trail)
			}
			return
		} else {
			trail := append(trail, at)
			next := Pair{at.x + 1, at.y}
			if inGrid(grid, next) && canClimb(grid, at, next) {
				dfs(grid, next, end, trail, count)
			}

			next = Pair{at.x, at.y + 1}
			if inGrid(grid, next) && canClimb(grid, at, next) {
				dfs(grid, next, end, trail, count)
			}

			next = Pair{at.x - 1, at.y}
			if inGrid(grid, next) && canClimb(grid, at, next) {
				dfs(grid, next, end, trail, count)
			}

			next = Pair{at.x, at.y - 1}
			if inGrid(grid, next) && canClimb(grid, at, next) {
				dfs(grid, next, end, trail, count)
			}
		}
	}
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func h(grid [][]rune, at Pair, end Pair) int {
	return AbsInt(at.x-end.x) + AbsInt(at.y-end.y)
}

func getLowestFScore(openSet map[Pair]bool, fScore map[Pair]int) Pair {
	lowestKey := Pair{}
	for k, _ := range openSet {
		lowestKey = k
		break
	}

	for k, _ := range openSet {
		if getScore(fScore, k) < getScore(fScore, lowestKey) {
			lowestKey = k
		}
	}
	return lowestKey
}

func getNeighbors(grid [][]rune, at Pair) []Pair {
	neighs := []Pair{}
	next := Pair{at.x + 1, at.y}
	if inGrid(grid, next) && canClimb(grid, at, next) {
		neighs = append(neighs, next)
	}

	next = Pair{at.x, at.y + 1}
	if inGrid(grid, next) && canClimb(grid, at, next) {
		neighs = append(neighs, next)
	}

	next = Pair{at.x - 1, at.y}
	if inGrid(grid, next) && canClimb(grid, at, next) {
		neighs = append(neighs, next)
	}

	next = Pair{at.x, at.y - 1}
	if inGrid(grid, next) && canClimb(grid, at, next) {
		neighs = append(neighs, next)
	}
	return neighs
}

func containsKey(cameFrom map[Pair]Pair, at Pair) bool {
	_, ok := cameFrom[at]
	return ok
}

func reconstructPath(cameFrom map[Pair]Pair, at Pair) []Pair {
	path := []Pair{}
	path = append(path, at)
	for containsKey(cameFrom, at) {
		at = cameFrom[at]
		path = append(path, at)
	}
	return path
}

func getScore(score map[Pair]int, at Pair) int {
	if val, ok := score[at]; ok {
		return val
	} else {
		return math.MaxInt
	}
}

func astar(grid [][]rune, start Pair, end Pair, trail []Pair) ([]Pair, error) {
	openSet := make(map[Pair]bool)

	cameFrom := make(map[Pair]Pair)

	gScore := make(map[Pair]int)
	gScore[start] = 0

	fScore := make(map[Pair]int)
	fScore[start] = h(grid, start, end)

	openSet[start] = true

	for len(openSet) > 0 {
		current := getLowestFScore(openSet, fScore)
		// fmt.Println(current)
		if current == end {
			return reconstructPath(cameFrom, current), nil
		}

		delete(openSet, current)
		neighbors := getNeighbors(grid, current)
		// fmt.Println(neighbors)

		for _, neighbor := range neighbors {
			// fmt.Println(neighbor, gScore[neighbor])

			tentativeGScore := gScore[current] + 1
			if tentativeGScore < getScore(gScore, neighbor) {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + h(grid, neighbor, end)
				if _, ok := openSet[neighbor]; !ok {
					openSet[neighbor] = true
				}
			}

		}
	}

	return []Pair{}, errors.New("no path")

}

type Pair struct {
	x, y int
}

func (p Pair) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
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

	grid := [][]rune{}

	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		grid = append(grid, row)
	}

	sx, sy := findInGrid(grid, 'S')
	ex, ey := findInGrid(grid, 'E')

	start := Pair{sx, sy}
	end := Pair{ex, ey}
	grid[start.x][start.y] = 'a'
	grid[end.x][end.y] = 'z'
	// dfs(grid, start, end, []Pair{}, &count)

	// fmt.Printf("count is %d\n", count)

	path, err := astar(grid, start, end, []Pair{})
	if err == nil {
		fmt.Printf("count is %d\n", len(path)-1)
	} else {
		log.Fatal(err)
	}

	aStartPaths := []Pair{}

	for r, row := range grid {
		for c, col := range row {
			if col == 'a' {
				aStartPaths = append(aStartPaths, Pair{r, c})
			}
		}
	}

	// NOTE: this (and therefore also the previous one given we can use one computation) can be done more efficiently with dykstra
	minAPath := math.MaxInt
	for _, pair := range aStartPaths {
		path, err := astar(grid, pair, end, []Pair{})
		if err == nil {
			if (len(path) - 1) < minAPath {
				minAPath = (len(path) - 1)
			}
		}
	}

	fmt.Printf("min path length is %d\n", minAPath)

}
