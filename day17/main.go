package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var year = 2022
var day = 17
var path = os.Getenv("ADVENT_HOME")

type Pair struct {
	x, y int
}

type Rock struct {
	formation int
	offset    Pair
}

type RockType struct {
	formation             []Pair
	boundingBoxBottomLeft Pair
	boundingBoxTopRight   Pair
}

const (
	_ = iota
	LEFT
	DOWN
	RIGHT
)

const (
	_ = iota
	HORZ_BAR
	PLUS
	REVERSE_L
	VERT_BAR
	SQUARE
)

var rockA = []Pair{Pair{0, 0}, Pair{1, 0}, Pair{2, 0}, Pair{3, 0}}
var rockB = []Pair{Pair{1, 0}, Pair{0, 1}, Pair{1, 1}, Pair{2, 1}, Pair{1, 2}}
var rockC = []Pair{Pair{0, 0}, Pair{1, 0}, Pair{2, 0}, Pair{2, 1}, Pair{2, 2}}
var rockD = []Pair{Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}}
var rockE = []Pair{Pair{0, 0}, Pair{1, 0}, Pair{0, 1}, Pair{1, 1}}

func newPairOrigin() Pair {
	return Pair{0, 0}
}

func (r Rock) String() string {
	return fmt.Sprintf("[Type=%d, offset=%s]", r.formation, r.offset)
}

func newRockType(formation []Pair) RockType {
	bbbl, bbtr := getBoundingBox(formation, newPairOrigin())
	return RockType{formation: formation, boundingBoxBottomLeft: bbbl, boundingBoxTopRight: bbtr}
}

var rockTypeA = newRockType(rockA)
var rockTypeB = newRockType(rockB)
var rockTypeC = newRockType(rockC)
var rockTypeD = newRockType(rockD)
var rockTypeE = newRockType(rockE)

var ROCK_MAP = map[int]RockType{
	HORZ_BAR:  rockTypeA,
	PLUS:      rockTypeB,
	REVERSE_L: rockTypeC,
	VERT_BAR:  rockTypeD,
	SQUARE:    rockTypeE,
}

var ROCK_SET = []int{HORZ_BAR, PLUS, REVERSE_L, VERT_BAR, SQUARE}

func main() {

	dataPath := fmt.Sprintf("%s/%d/data/day%d/dataSample.txt", path, year, day)
	fmt.Println(dataPath)

	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	jetPattern := ""
	for scanner.Scan() {
		jetPattern = scanner.Text()
	}

	numberOfRocks := 2
	simulateRocks(numberOfRocks, jetPattern)

}

var leftBound = 0
var rightBound = 8

type State struct {
	rockStack []Rock
}

func (p Pair) add(p2 Pair) Pair {
	return Pair{p.x + p2.x, p.y + p2.y}
}

func (p Pair) String() string {
	return fmt.Sprintf("Pair::{%d, %d}", p.x, p.y)
}

func getHighestPoint(formation []Pair, offset Pair) Pair {
	hp := Pair{0, 0}
	for _, point := range formation {
		op := point.add(offset)
		if op.y > hp.y {
			hp = op
		}
	}
	return hp
}

func getLowestPoint(formation []Pair, offset Pair) Pair {
	hp := Pair{0, 10}.add(offset)
	for _, point := range formation {
		op := point.add(offset)
		if op.y < hp.y {
			hp = op
		}
	}
	return hp
}

func getLeftMostPoint(formation []Pair, offset Pair) Pair {
	hp := Pair{8, 0}.add(offset)
	for _, point := range formation {
		op := point.add(offset)
		if op.x < hp.x {
			hp = op
		}
	}
	return hp
}

func getRightMostPoint(formation []Pair, offset Pair) Pair {
	hp := Pair{0, 0}.add(offset)
	for _, point := range formation {
		op := point.add(offset)
		if op.x > hp.x {
			hp = op
		}
	}
	return hp
}

func getBoundingBox(formation []Pair, offset Pair) (Pair, Pair) {
	l := getLeftMostPoint(formation, offset).x
	rr := getRightMostPoint(formation, offset).x
	t := getHighestPoint(formation, offset).y
	b := getLowestPoint(formation, offset).y

	return Pair{l, b}, Pair{rr, t}
}

//-----------------------------

func (r *RockType) getBoundingBox(offset Pair) (Pair, Pair) {

	bl := r.boundingBoxBottomLeft
	tr := r.boundingBoxTopRight

	return bl.add(offset), tr.add(offset)
}

// -------------------

func simulateRocks(numOfRocks int, jetPattern string) {
	rockCount := 0
	tick := 0
	jetCount := 0
	state := State{rockStack: []Rock{}}
	jetPatternArray := []rune(jetPattern)

	for rockCount < numOfRocks {
		// spawn rock
		newRock := state.spawnRock(rockCount)
		for state.canRockFall(newRock) {
			fmt.Print("in loop")
			fmt.Println(newRock)
			if tick%2 == 0 {
				fmt.Print(" updating jets ")

				newRock = updateJetHittingRock(newRock, jetPatternArray, jetCount, state)
				jetCount++
			} else {
				fmt.Print(" updating falling ")

				newRock = updateFallingRock(newRock, state)
			}
			tick++
			fmt.Println()
		}

		// settle rock once it cannot fall anymore
		state.rockStack = append(state.rockStack, newRock)
		rockCount++
	}

	state.PrintPicture()
	state.PrintState()

}

func (s *State) spawnRock(time int) Rock {
	rockChoice := ROCK_SET[time%len(ROCK_SET)]

	lowYBound := s.getMaxY()
	offset := Pair{3, lowYBound + 4}
	newRock := Rock{formation: rockChoice, offset: offset}

	return newRock
}

func (s *State) getMaxY() int {
	hp := 0
	for _, rock := range s.rockStack {
		op := rock.getMaxY()
		if op > hp {
			hp = op
		}
	}
	return hp
}

func updateFallingRock(rock Rock, state State) Rock {
	shift := newPairOrigin()
	shift.y = -1

	newRock := rock
	newRock.offset = newRock.offset.add(shift)
	fmt.Print("updateFallingRock", rock, newRock)

	if newRock.getMinY() < 0 {
		return rock
	}

	if checkNewRockPlacement(newRock, state) {
		fmt.Print("new rock location is ok")

		return newRock
	}
	fmt.Print("new rock location is INVALID")

	return rock
}

func (s State) String() string {
	return fmt.Sprintf("%s", s.rockStack)
}

func (s State) PrintState() {
	fmt.Println(s)
}

func updateJetHittingRock(rock Rock, jetPattern []rune, jetTick int, state State) Rock {
	dir := jetPattern[jetTick%len(jetPattern)]
	dirEnum := LEFT
	if dir == '<' {
		dirEnum = LEFT
	} else if dir == '>' {
		dirEnum = RIGHT
	}

	shift := newPairOrigin()
	if dirEnum == LEFT {
		shift.x = -1
	} else if dirEnum == RIGHT {
		shift.x = 1
	}

	newRock := rock
	newRock.offset = newRock.offset.add(shift)

	if newRock.getMinX() <= 0 || newRock.getMaxX() >= 8 {
		return rock
	}

	if checkNewRockPlacement(newRock, state) {
		return newRock
	}

	return rock
}

func checkNewRockPlacement(rock Rock, state State) bool {
	for _, other := range state.rockStack {
		if boundingBoxesCollide(rock, other) {
			fmt.Print("COLLISION")
			if rockFormationIntersect(rock, other) {
				fmt.Print("INTERSECTION")
				return false
			}
		}
	}

	return true
}

func (r *Rock) getRockType() RockType {
	return ROCK_MAP[r.formation]
}

func (r *Rock) getBoundingBox() (Pair, Pair) {
	rockType := r.getRockType()
	return rockType.getBoundingBox(r.offset)
}

func (r *Rock) getMaxY() int {
	return ROCK_MAP[r.formation].boundingBoxTopRight.add(r.offset).y
}

func (r *Rock) getMinY() int {
	return ROCK_MAP[r.formation].boundingBoxBottomLeft.add(r.offset).y
}

func (r *Rock) getMaxX() int {
	return ROCK_MAP[r.formation].boundingBoxTopRight.add(r.offset).x
}

func (r *Rock) getMinX() int {
	return ROCK_MAP[r.formation].boundingBoxBottomLeft.add(r.offset).x
}

func (s *State) canRockMove(rock Rock, direction int) bool {

	return false
}

func (s State) canRockFall(rock Rock) bool {

	return rock.getMinY() > 0 && checkNewRockPlacement(updateFallingRock(rock, s), s)
}

func (r Rock) getPairsInSpace() []Pair {
	pairsInSpace := []Pair{}
	for _, pair := range ROCK_MAP[r.formation].formation {
		pairInSpace := pair.add(r.offset)
		pairsInSpace = append(pairsInSpace, pairInSpace)
	}
	return pairsInSpace
}

func (s State) PrintPicture() {
	picture := []string{}
	pairMap := make(map[int][]Pair)
	for _, rock := range s.rockStack {
		for _, pair := range rock.getPairsInSpace() {
			yLoc := pair.y
			if _, ok := pairMap[yLoc]; !ok {
				pairMap[yLoc] = []Pair{}
			}
			pairMap[yLoc] = append(pairMap[yLoc], pair)
		}
	}

	for i := s.getMaxY() + 4; i >= 0; i-- {
		row := pairMap[i]
		rowStr := createRow(row)
		picture = append(picture, rowStr)
	}

	acc := ""
	for i := s.getMaxY() + 4; i >= 0; i-- {
		acc += picture[i] + "\n"
	}

	fmt.Println(acc)
}

func createRow(pairs []Pair) string {
	xs := make(map[int]bool)
	for _, pair := range pairs {
		xs[pair.x] = true
	}

	acc := ""
	for i := 0; i <= 8; i++ {
		if i == 0 || i == 8 {
			acc += "|"
		} else if _, ok := xs[i]; ok {
			acc += "#"
		} else {
			acc += "."
		}
	}
	return acc
}

func inInterval(left, right, subject int) bool {
	return left <= subject && subject <= right
}

func boundingBoxesCollide(a, b Rock) bool {
	aMinBB, aMaxBB := a.getBoundingBox()
	aMinX, aMinY := aMinBB.x, aMinBB.y
	aMaxX, aMaxY := aMaxBB.x, aMaxBB.y

	bMinBB, bMaxBB := b.getBoundingBox()
	bMinX, bMinY := bMinBB.x, bMinBB.y
	bMaxX, bMaxY := bMaxBB.x, bMaxBB.y

	xIntersects := inInterval(aMinX, aMaxX, bMinX) || inInterval(bMinX, bMaxX, aMinX)
	yIntersects := inInterval(aMinY, aMaxY, bMinY) || inInterval(bMinY, bMaxY, aMinY)

	return xIntersects && yIntersects
}

func rockFormationIntersect(a, b Rock) bool {
	dict := map[string]bool{}
	for _, pair := range ROCK_MAP[a.formation].formation {
		pairInSpace := pair.add(a.offset)
		dict[pairInSpace.String()] = true
	}

	for _, pair := range ROCK_MAP[b.formation].formation {
		pairInSpace := pair.add(b.offset)
		if _, ok := dict[pairInSpace.String()]; ok {
			return true
		} else {
			dict[pairInSpace.String()] = true
		}
	}
	return false
}
