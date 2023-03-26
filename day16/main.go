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
var day = 16
var path = os.Getenv("ADVENT_HOME")

type Valve struct {
	name        string
	flowRate    int
	connections []string
}

type Path struct {
	flow  int
	nodes []string
}

// Floyd-Warshall algo
// produces shortest path between all nodes
func findShortestPathPairs(valveMap map[string]Valve) map[string]map[string]int {
	inf := len(valveMap)*(len(valveMap)-1)/2 + 1

	grid := map[string]map[string]int{}

	// create grid with default value of inf (largest possible value)
	for key := range valveMap {
		grid[key] = make(map[string]int)
		for key2 := range valveMap {
			grid[key][key2] = inf
		}
	}

	// set weights of graph (aka weight of 1 for adjacent nodes)
	for key := range valveMap {
		for _, connected := range valveMap[key].connections {
			grid[key][connected] = 1
		}
	}

	// set reflexive edges to 0
	for key := range valveMap {
		grid[key][key] = 0
	}

	// find shorttest path pairs
	for k := range valveMap {
		for i := range valveMap {
			for j := range valveMap {
				if grid[i][j] > (grid[i][k] + grid[k][j]) {
					grid[i][j] = grid[i][k] + grid[k][j]
				}
			}
		}
	}

	return grid
}

var shortestPathPairs map[string]map[string]int
var valveMap map[string]Valve

func main() {
	fmt.Println(path)
	dataPath := fmt.Sprintf("%s/%d/data/day%d/data.txt", path, year, day)
	fmt.Println(dataPath)

	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	valveMap = make(map[string]Valve)

	for scanner.Scan() {
		output := scanner.Text()
		tokens := strings.Split(output, " ")
		name := tokens[1]
		rate := tokens[4]
		rate = strings.TrimPrefix(rate, "rate=")
		rate = strings.TrimSuffix(rate, ";")
		rateVal, _ := strconv.Atoi(rate)
		connections := []string{}
		for i := 9; i < len(tokens); i++ {
			val := tokens[i]
			val = strings.TrimSuffix(val, ",")
			connections = append(connections, val)
		}
		valve := Valve{name, rateVal, connections}
		valveMap[valve.name] = valve
	}

	// fmt.Println(valveMap)

	shortestPathPairs = findShortestPathPairs(valveMap)

	// from valveMap, remove nodes that have zero flow
	// we dont need them anymore, because the shortestPathPairs map will gives us the distance from any two nodes
	for key := range valveMap {
		if valveMap[key].flowRate == 0 {
			delete(valveMap, key)
		}
	}
	// fmt.Println(valveMap)

	// solve p1
	maxRate := solveP1()
	fmt.Printf("max Rate : %d\n", maxRate)

	// solve p2
	maxRate = solveP2()
	fmt.Printf("max Rate : %d\n", maxRate)
}

// part 1
// You need to find the traversal path through the graph that produces the highest released pressure
func solveP1() int {
	// find the list of paths that can be starting with the node "AA", with 30 seconds
	paths := findMaxPath("AA", 30, Path{0, make([]string, 0)}, make(map[string]bool))
	// find the path with the largest accumulated flow
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].flow > paths[j].flow
	})
	return paths[0].flow
}

// this method builds off part 1. now that two traversals will be done at the same time, you compute the list
// of possible traversal paths and compare them to see the set of two paths that could be done with no overlap.
// comparing the set of traversal pairs that could be done with no overlap, you find the pair which produces the
// highest combined released pressure. that is the answer to Part 2
func solveP2() int {
	paths := findMaxPath("AA", 26, newPath(), make(map[string]bool))
	maxRate := 0
	for _, path := range paths {
		if len(path.nodes) > 0 {
			for _, other := range paths {
				possibleRate := path.flow + other.flow
				if maxRate < possibleRate && len(other.nodes) > 0 && differentPaths(path, other) {
					maxRate = possibleRate
				}
			}
		}
	}

	return maxRate
}

// with recursive backtracking, this function computes the potential paths of valve activation
// and the totalFlowRate during such traversal.
func findMaxPath(at string, time int, path Path, seenMap map[string]bool) []Path {
	paths := []Path{path}

	// iterate over all nodes that increase pressure on valve release
	for nextNode := range valveMap {
		// compute the time it will take to get from at node to nextNode
		nextTime := time - shortestPathPairs[at][nextNode] - 1
		// if nextNode has been traversed (and therefore opened) or if there is no time to get to the next node,
		// skip and try another potential path
		if seenMap[nextNode] || nextTime <= 0 {
			continue
		}

		// make a copy of the current path and increment based on increase flow rate
		nextMap := copyMap(seenMap)
		nextMap[nextNode] = true
		newPath := path.copy()
		newPath.addToPath(nextTime*valveMap[nextNode].flowRate, nextNode)

		// add newly created path along with all paths that can be created _off_ of it to list
		paths = append(paths, findMaxPath(nextNode, nextTime, newPath, nextMap)...)

	}

	return paths
}

// given two paths, return true if their traversal path is exclusive
func differentPaths(a Path, b Path) bool {
	m := make(map[string]bool)
	for _, node := range a.nodes {
		m[node] = true
	}

	for _, node := range b.nodes {
		if _, ok := m[node]; ok {
			return false
		}
	}
	return true
}

func (p *Path) addToPath(flow int, valve string) {
	p.flow += flow
	p.nodes = append(p.nodes, valve)
}

func (p Path) copy() Path {
	nodes := make([]string, len(p.nodes))
	copy(nodes, p.nodes)
	return Path{p.flow, nodes}
}

func copyMap(mm map[string]bool) map[string]bool {
	mmm := map[string]bool{}
	for key := range mm {
		mmm[key] = mm[key]
	}
	return mmm
}

func newPath() Path {
	return Path{0, make([]string, 0)}
}
