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
var day = 13
var path = os.Getenv("ADVENT_HOME")

const (
	_ = iota
	RIGHT_ORDER
	WRONG_ORDER
	CONTINUE
)

func printEnum(i int) string {
	if i == RIGHT_ORDER {
		return "RIGHT_ORDER"
	}

	if i == WRONG_ORDER {
		return "WRONG_ORDER"
	}

	return "CONTINUE"
}

func splitList(str string) []string {
	count := 0
	commas := []int{}

	for i, r := range str {
		if string(r) == "[" {
			count += 1
		}
		if string(r) == "]" {
			count -= 1
		}
		if string(r) == "," {
			if count == 0 {
				commas = append(commas, i)
			}
		}
	}

	output := []string{}
	i := 0

	for _, indexOfComma := range commas {
		sub := str[i:indexOfComma]
		output = append(output, sub)
		i = indexOfComma + 1
	}
	output = append(output, str[i:])

	return output
}

func getValues(str string) []string {
	val := strings.TrimPrefix(str, "[")
	val = strings.TrimSuffix(val, "]")

	vals := splitList(val)
	out := []string{}
	for _, v := range vals {
		if len(strings.TrimSpace(v)) > 0 {
			out = append(out, v)
		}
	}
	return out
}

func isList(str string) bool {
	return len(str) >= 2 && strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")
}

func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func wrapAsList(str string) string {
	return "[" + str + "]"
}

func compareValues(left string, right string) int {
	if isNumber(left) && isNumber(right) {
		lVal, _ := strconv.Atoi(left)
		rVal, _ := strconv.Atoi(right)
		if lVal < rVal {
			return RIGHT_ORDER
		} else if lVal > rVal {
			return WRONG_ORDER
		} else {
			return CONTINUE
		}
	} else if isList(left) && isList(right) {
		return compareList(left, right)
	} else {
		nLeft, nRight := left, right
		if isNumber(nLeft) {
			nLeft = wrapAsList(nLeft)
		}
		if isNumber(nRight) {
			nRight = wrapAsList(nRight)
		}
		return compareList(nLeft, nRight)
	}
}

func compareList(left string, right string) int {
	leftVals := getValues(left)
	rightVals := getValues(right)

	l := len(leftVals)
	r := len(rightVals)
	// fmt.Printf("compareList %s ||| %s\n", leftVals, rightVals)

	i := 0

	if l == 0 && r > 0 {
		return RIGHT_ORDER
	}

	if l > 0 && r == 0 {
		return WRONG_ORDER
	}

	for i < l && i < r {
		val := compareValues(leftVals[i], rightVals[i])
		if val == RIGHT_ORDER || val == WRONG_ORDER {
			return val
		}
		i++
	}

	if i == l && i == r {
		return CONTINUE
	}

	if i == l && i < r {
		return RIGHT_ORDER
	}

	return WRONG_ORDER
}

func isRightOrder(left string, right string) (bool, int) {
	output := compareList(left, right)
	return output == RIGHT_ORDER, output
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
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			data = append(data, scanner.Text())
		}
	}

	acc := 0
	count := 1
	for i := 0; i < len(data); i += 2 {
		msg := data[i]
		test := data[i+1]
		ok, _ := isRightOrder(msg, test)
		if ok {
			acc += count
		}
		count += 1
	}

	fmt.Printf("sum = %d\n", acc)

	// part 2
	// add divider tokens
	LEFT_DIV := "[[2]]"
	RIGHT_DIV := "[[6]]"
	data = append(data, LEFT_DIV)
	data = append(data, RIGHT_DIV)

	sort.Slice(data, func(i, j int) bool {
		seqI := data[i]
		seqJ := data[j]
		ok, _ := isRightOrder(seqI, seqJ)
		return ok
	})

	ld := 0
	rd := 0

	for i := 0; i < len(data); i++ {
		if data[i] == LEFT_DIV {
			ld = i
		}
		if data[i] == RIGHT_DIV {
			rd = i
		}
	}

	// shift for 1st based index
	rd++
	ld++
	fmt.Printf("decoderKey (ld,rd)=(%d,%d) = %d\n", ld, rd, ld*rd)

}
