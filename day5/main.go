package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var year = 2022
var day = 5
var path = os.Getenv("ADVENT_HOME")

func parseDataCargo(dataCargo []string) [][]string {
	stacks := [][]string{}
	numOfStacks := 9
	for i := 0; i < numOfStacks; i++ {
		stacks = append(stacks, []string{})
	}

	firstRow := dataCargo[0]
	fmt.Printf("%d\n", len(firstRow))
	for _, str := range dataCargo {
		stackCount := 0
		for i := 1; i < len(str); i += 4 {
			value := str[i]
			if !unicode.IsSpace(rune(value)) {
				stacks[stackCount] = append(stacks[stackCount], string(value))
			}
			stackCount += 1
		}
	}

	// for i := 0; i < len(stacks); i++ {
	// 	fmt.Println(stacks[i])
	// }

	// fmt.Printf("%d\n\n", (len(firstRow)-1)/3)

	// for i := 0; i < len(firstRow)
	return stacks
}

type instruction struct {
	amount, from, to int
}

func (i instruction) String() string {
	return fmt.Sprintf("(move:%d, from:%d, to:%d)", i.amount, i.from, i.to)
}

func parseDataInstructions(data []string) []instruction {
	output := []instruction{}

	for _, str := range data {
		tokens := strings.Split(str, " ")
		// 1,3,5
		a, _ := strconv.Atoi(tokens[1])
		b, _ := strconv.Atoi(tokens[3])
		c, _ := strconv.Atoi(tokens[5])

		inst := instruction{a, b - 1, c - 1}
		// fmt.Println(inst)

		output = append(output, inst)
	}
	return output
}

func handleInstruction(cargo [][]string, inst instruction) [][]string {
	amount := inst.amount
	from := inst.from
	to := inst.to
	// fmt.Println(inst)

	for i := 0; i < amount; i++ {
		if len(cargo[from]) > 0 {
			// grab box and remove it from current stack
			box := cargo[from][0]
			cargo[from] = cargo[from][1:]
			cargo[to] = append([]string{box}, cargo[to]...)
		}
	}

	return cargo

}

func printCargo(cargo [][]string) {
	for i := 0; i < len(cargo); i++ {
		fmt.Printf("%d ", i)
		fmt.Println(cargo[i])
	}

}

func handleInstructionV2(cargo [][]string, inst instruction) [][]string {

	printCargo(cargo)
	amount := inst.amount
	from := inst.from
	to := inst.to
	fmt.Println(inst)
	if amount > 0 {
		boxStack := cargo[from][:amount]
		fmt.Printf("BoxStack = ")
		fmt.Println(boxStack)
		if len(boxStack) > 0 {
			rest := cargo[from][amount:]
			fmt.Printf("new cargo[from] = ")
			fmt.Println(rest)
			cargo[from] = append([]string{}, rest...)
			cargo[to] = append(boxStack, cargo[to]...)

			fmt.Printf("cargo[from] = ")
			fmt.Println(cargo[from])
			fmt.Printf("cargo[to] = ")
			fmt.Println(cargo[to])
		}
	}

	printCargo(cargo)

	return cargo

}

func processPartOne(cargo [][]string, instructions []instruction) [][]string {
	for _, inst := range instructions {
		cargo = handleInstruction(cargo, inst)
	}

	printCargo(cargo)

	return cargo

}

func processPartTwo(cargo [][]string, instructions []instruction) [][]string {
	for _, inst := range instructions {
		fmt.Println("----------------------")
		cargo = handleInstructionV2(cargo, inst)
		fmt.Println("----------------------")
	}

	printCargo(cargo)

	return cargo

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

	dataCargo := []string{}
	dataInstructions := []string{}

	const (
		_ = iota
		CARGO_MODE
		INSTRUCTION_MODE
	)

	scanMode := CARGO_MODE

	for scanner.Scan() {
		originalTextData := scanner.Text()
		textData := strings.TrimSpace(originalTextData)
		if len(textData) == 0 {
			scanMode = INSTRUCTION_MODE
			dataCargo = dataCargo[:len(dataCargo)-1]
		} else {
			if scanMode == CARGO_MODE {
				dataCargo = append(dataCargo, originalTextData)
			} else {
				dataInstructions = append(dataInstructions, originalTextData)
			}
		}
	}

	fmt.Println("Data Cargo")
	// fmt.Println("----------------------")
	// fmt.Println("Data Instructions")
	// fmt.Println(strings.Join(dataInstructions, "\n"))

	cargo := parseDataCargo(dataCargo)
	printCargo(cargo)
	instructions := parseDataInstructions(dataInstructions)
	finalCargo := processPartOne(cargo, instructions)
	top := []string{}
	for _, l := range finalCargo {
		top = append(top, l[0])
	}
	fmt.Printf("p1. == %s\n", strings.Join(top, ""))

	cargo = parseDataCargo(dataCargo)
	finalCargo = processPartTwo(cargo, instructions)
	top = []string{}
	for _, l := range finalCargo {
		if len(l) > 0 {
			top = append(top, l[0])

		} else {
			top = append(top, " ")
		}
	}

	fmt.Printf("p2. == %s\n", strings.Join(top, ""))
}
