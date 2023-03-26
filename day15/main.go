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
var day = 15
var path = os.Getenv("ADVENT_HOME")

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Tuple struct {
	x, y int
}

func (t Tuple) String() string {
	return fmt.Sprintf("(%d,%d)", t.x, t.y)
}

type Sensor struct {
	source Tuple
	beacon Tuple
}

func TaxiCabDist(a Tuple, b Tuple) int {
	return AbsInt(a.x-b.x) + AbsInt(a.y-b.y)
}

func (s Sensor) String() string {
	return fmt.Sprintf(" ")
}

func (s Sensor) TaxiCabDist() int {
	return TaxiCabDist(s.source, s.beacon)
}

func (s Sensor) LeftBound() int {
	return s.source.x - s.TaxiCabDist()
}

func (s Sensor) RightBound() int {
	return s.source.x + s.TaxiCabDist()
}

func StringToTuple(str string) Tuple {
	tokens := strings.Split(str, ",")
	x := strings.TrimSpace(tokens[0])
	y := strings.TrimSpace(tokens[1])

	x = strings.TrimPrefix(x, "x=")
	y = strings.TrimPrefix(y, "y=")

	xVal, _ := strconv.Atoi(x)
	yVal, _ := strconv.Atoi(y)

	return Tuple{xVal, yVal}
}

func StringToSensor(dataLine string) Sensor {
	tokens := strings.Split(dataLine, ":")
	left := tokens[0]
	right := tokens[1]

	left = strings.TrimPrefix(left, "Sensor at ")
	right = strings.TrimPrefix(right, " closest beacon is at ")
	l := StringToTuple(left)
	r := StringToTuple(right)

	return Sensor{l, r}
}

func GetScannedSpots(sensors []Sensor, minX, maxX, atRow int) int {
	count := 0

	occupied := make(map[Tuple]bool)
	for _, sensor := range sensors {
		occupied[sensor.source] = true
		occupied[sensor.beacon] = true
	}

	for i := minX; i <= maxX; i++ {
		inRange := false
		spot := Tuple{i, atRow}
		if _, ok := occupied[spot]; ok {
			continue
		}
		for _, sensor := range sensors {
			if TaxiCabDist(spot, sensor.source) <= sensor.TaxiCabDist() {
				inRange = true
			}
			if inRange {
				break
			}
		}
		if inRange {
			count += 1
		}
	}

	return count
}

func NotInRange(sensors []Sensor, at Tuple) bool {
	for _, sensor := range sensors {
		if TaxiCabDist(at, sensor.source) <= sensor.TaxiCabDist() {
			return false
		}
	}
	return true
}

func CalculateTuningFrequency(t Tuple) int {
	return (t.x * 4000000) + t.y
}

func InBoundsChecker(maxX int, maxY int) func(Tuple) bool {
	return func(t Tuple) bool {
		return (0 <= t.x && t.x <= maxX) && (0 <= t.y && t.y <= maxY)
	}
}

func FindBeaconSpot(sensors []Sensor, maxX, maxY, atRow int) int {
	// find range with biggest area
	// navigate circumference of area to find spot
	for _, sensor := range sensors {
		sensorMax := sensor

		// fmt.Println(sensorMax.source)

		check := InBoundsChecker(maxX, maxY)

		dist := sensorMax.TaxiCabDist() + 1
		x, y := sensorMax.source.x+dist, sensorMax.source.y
		// up left

		for i := 0; i <= dist; i++ {
			at := Tuple{x - i, y + i}
			// fmt.Println(at)
			if check(at) && NotInRange(sensors, at) {
				return CalculateTuningFrequency(at)
			}
		}

		// down left
		x, y = sensorMax.source.x, sensorMax.source.y+dist

		for i := 0; i <= dist; i++ {
			at := Tuple{x - i, y - i}
			// fmt.Println(at)

			if check(at) && NotInRange(sensors, at) {
				return CalculateTuningFrequency(at)
			}
		}

		x, y = sensorMax.source.x-dist, sensorMax.source.y
		// down right
		for i := 0; i <= dist; i++ {
			at := Tuple{x + i, y - i}
			// fmt.Println(at)

			if check(at) && NotInRange(sensors, at) {
				return CalculateTuningFrequency(at)
			}
		}

		x, y = sensorMax.source.x, sensorMax.source.y-dist
		// up right
		for i := 0; i <= dist; i++ {
			at := Tuple{x + i, y + i}
			// fmt.Println(at)

			if check(at) && NotInRange(sensors, at) {
				return CalculateTuningFrequency(at)
			}
		}
	}

	return 0
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

	sensors := []Sensor{}
	for scanner.Scan() {
		line := scanner.Text()
		sensor := StringToSensor(line)
		sensors = append(sensors, sensor)
	}

	maxX := 0

	for _, sensor := range sensors {
		rightBound := sensor.RightBound()
		if maxX < rightBound {
			maxX = rightBound
		}
	}

	minX := maxX
	for _, sensor := range sensors {
		leftBound := sensor.LeftBound()
		if minX > leftBound {
			minX = leftBound
		}
	}

	fmt.Printf("l=%d, r=%d\n", minX, maxX)
	atRow := 2000000
	// atRow := 10

	count := 0
	count = GetScannedSpots(sensors, minX, maxX, atRow)
	fmt.Printf("count = %d\n", count)

	sizeX := 4000000
	sizeY := 4000000

	// sizeX := 20
	// sizeY := 20

	count = FindBeaconSpot(sensors, sizeX, sizeY, atRow)
	fmt.Printf("count = %d\n", count)
}
