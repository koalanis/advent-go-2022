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
var day = 11
var path = os.Getenv("ADVENT_HOME")

type MonkeyOperation func(int64) int64

type MonkeyTest func(int64) bool

func identity() MonkeyOperation {
	fmt.Println("Using Identity, check code")
	return func(i int64) int64 {
		return i
	}
}

func addOperation(addBy int64) MonkeyOperation {
	return func(i int64) int64 {
		return i + addBy
	}
}

func multiplyOperation(multBy int64) MonkeyOperation {
	return func(i int64) int64 {
		return i * multBy
	}
}

func squareOperation() MonkeyOperation {
	return func(i int64) int64 {
		return i * i
	}
}

func divisibleByTest(divisibleBy int64) MonkeyTest {
	return func(i int64) bool {
		return i%divisibleBy == 0
	}
}

type Monkey struct {
	itemQueue []int64
	operation MonkeyOperation
	test      MonkeyTest
	tossMap   map[bool]int

	number         int
	itemsInspected int
}

func intSliceToString(nums []int64) string {
	valuesText := []string{}

	return fmt.Sprintf("[%s]", strings.Join(valuesText, ","))
}

func boolMapToString(m map[bool]int) string {
	valuesText := []string{}

	for k := range m {
		valuesText = append(valuesText, fmt.Sprintf("%t=%d", k, m[k]))
	}

	return fmt.Sprintf("[%s]", strings.Join(valuesText, ","))
}

func (m Monkey) String() string {
	return fmt.Sprintf("Monkey{num:%d\nitems=%s\nitemsInspected=%d\ntossMap=%s\n}", m.number, intSliceToString(m.itemQueue), m.itemsInspected, boolMapToString(m.tossMap))
}

func getStartingItems(line string) []int64 {
	data := strings.TrimPrefix(line, "Starting items:")
	data = strings.TrimSpace(data)
	items := []int64{}
	tokens := strings.Split(data, ",")
	for _, str := range tokens {
		s := strings.TrimSpace(str)
		val, _ := strconv.Atoi(s)
		items = append(items, int64(val))
	}
	return items
}

func getOperation(line string) MonkeyOperation {
	data := strings.TrimPrefix(line, "Operation: new =")
	data = strings.TrimSpace(data)

	if data == "old * old" {
		return squareOperation()
	} else if strings.HasPrefix(data, "old *") {
		tokens := strings.Split(data, "*")
		operand := strings.TrimSpace(tokens[1])
		val, _ := strconv.Atoi(operand)
		return multiplyOperation(int64(val))
	} else if strings.HasPrefix(data, "old +") {
		tokens := strings.Split(data, "+")
		operand := strings.TrimSpace(tokens[1])
		val, _ := strconv.Atoi(operand)
		return addOperation(int64(val))
	}
	return identity()
}

func getTest(testLine string, trueCase string, falseCase string) (int, MonkeyTest, map[bool]int) {
	data := strings.TrimSpace(strings.TrimPrefix(testLine, "Test: divisible by"))
	divisor, _ := strconv.Atoi(data)
	testOp := divisibleByTest(int64(divisor))
	testMap := make(map[bool]int)
	data = strings.TrimSpace(strings.TrimPrefix(trueCase, "If true: throw to monkey"))
	val, _ := strconv.Atoi(data)
	testMap[true] = val

	data = strings.TrimSpace(strings.TrimPrefix(falseCase, "If false: throw to monkey"))
	val, _ = strconv.Atoi(data)
	testMap[false] = val

	return divisor, testOp, testMap
}

func performOperation(data []string, rounds int, worry bool) {
	monkeys := []Monkey{}

	monkeyNum := 0
	divisors := []int{}
	for i := 0; i < len(data); i += 7 {
		monkey := Monkey{}
		startingItems := getStartingItems(data[i+1])
		operation := getOperation(data[i+2])
		divisor, test, tossMap := getTest(data[i+3], data[i+4], data[i+5])
		divisors = append(divisors, divisor)
		monkey.number = monkeyNum
		monkey.itemQueue = startingItems
		monkey.operation = operation
		monkey.test = test
		monkey.tossMap = tossMap
		monkey.itemsInspected = 0

		monkeys = append(monkeys, monkey)
		monkeyNum += 1
	}

	lcd := int64(1)
	for i := range divisors {
		lcd *= int64(divisors[i])
	}
	fmt.Printf("lcd = %d", lcd)
	round := 0
	totalRounds := rounds

	for round < totalRounds {
		// fmt.Printf("Round %d\n", round)
		for i := range monkeys {
			// fmt.Printf("Monkey%d\n", i)
			currentMonk := monkeys[i]
			for len(currentMonk.itemQueue) > 0 {
				item, queue := currentMonk.itemQueue[0], currentMonk.itemQueue[1:]
				currentMonk.itemQueue = queue
				nv := currentMonk.operation(item)
				if worry {
					nv = nv / 3
				} else {
					nv = nv % lcd
				}
				testVal := currentMonk.test(nv)
				throwTo := currentMonk.tossMap[testVal]
				monkeys[throwTo].itemQueue = append(monkeys[throwTo].itemQueue, nv)
				currentMonk.itemsInspected += 1
			}
			monkeys[i] = currentMonk
		}

		round++
	}

	inspections := []int{}
	for i := range monkeys {
		monk := monkeys[i]
		inspections = append(inspections, monk.itemsInspected)
	}

	sort.Slice(inspections, func(i, j int) bool {
		return inspections[j] < inspections[i]
	})

	fmt.Printf("monkeyBusiness = %d\n", inspections[0]*inspections[1])
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
		data = append(data, strings.TrimSpace(scanner.Text()))
	}

	performOperation(data, 20, true)
	performOperation(data, 10000, false)

}
