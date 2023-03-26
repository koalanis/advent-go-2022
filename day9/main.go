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
var day = 9
var path = os.Getenv("ADVENT_HOME")

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Vec struct {
	x, y int
}

func (v Vec) String() string {
	return fmt.Sprintf("(%d,%d)", v.x, v.y)
}

type Segment struct {
	head, tail Vec
	tailTrail  map[Vec]bool
}

func (r Segment) String() string {
	return fmt.Sprintf("%s-->%s", r.tail, r.head)
}

func GridLength(head Vec, tail Vec) int {
	length := 0
	if head.x == tail.x {
		length = AbsInt(head.y - tail.y)
	} else if head.y == tail.y {
		length = AbsInt(head.x - tail.x)
	} else {
		x1, y1 := head.x, head.y
		x2, y2 := tail.x, tail.y
		dx, dy := x1-x2, y1-y2

		ix := 0
		iy := 0
		if dx != 0 {
			ix = dx / AbsInt(dx)
		}
		if dy != 0 {
			iy = dy / AbsInt(dy)
		}

		dist := 0
		for dx != 0 && dy != 0 {
			dx -= ix
			dy -= iy
			dist += 1
		}

		for dx != 0 {
			dx -= ix
			dist += 1
		}

		for dy != 0 {
			dy -= iy
			dist += 1
		}

		length = dist
	}
	return length
}

func DeltaVector(head Vec, tail Vec) Vec {
	x1, y1 := head.x, head.y
	x2, y2 := tail.x, tail.y
	dx, dy := x1-x2, y1-y2
	ix := 0
	iy := 0
	if dx != 0 {
		ix = dx / AbsInt(dx)
	}
	if dy != 0 {
		iy = dy / AbsInt(dy)
	}
	return Vec{ix, iy}
}

func handleRopeMove(rope *Rope, dir string, steps int) {
	if steps > 0 {
		(*rope).HandleRopeMoves(dir, steps)
	}
}

type Rope struct {
	length    int
	segments  []Vec
	tailTrail map[Vec]bool
}

func (rope *Rope) GetTail() Vec {
	return rope.segments[len(rope.segments)-1]
}

func NewRope(length int) Rope {
	rope := Rope{}
	rope.length = length
	rope.segments = make([]Vec, length)
	rope.tailTrail = make(map[Vec]bool)
	for i := 0; i < length; i++ {
		rope.segments[i] = Vec{0, 0}
	}

	rope.tailTrail[rope.GetTail()] = true
	return rope
}

func NewRopeAt(tail Vec, head Vec) Rope {
	rope := NewRope(2)
	rope.segments[0] = head
	rope.segments[1] = tail
	return rope
}

func (r *Rope) HandleRopeMove(move string) {
	if move == "U" {
		r.segments[0].y += 1
	} else if move == "L" {
		r.segments[0].x -= 1
	} else if move == "R" {
		r.segments[0].x += 1
	} else if move == "D" {
		r.segments[0].y -= 1
	}
	if r.RopeTailShouldUpdate() {
		r.HandleRopeUpdate()
	}
}

func (r *Rope) HandleRopeUpdate() {
	if r.RopeTailShouldUpdate() {
		t := 1

		for t < r.length {
			for GridLength(r.segments[t-1], r.segments[t]) > 1 {
				prev, at := r.segments[t-1], r.segments[t]
				delta := DeltaVector(prev, at)
				updated := Vec{at.x + delta.x, at.y + delta.y}
				r.segments[t] = updated
			}
			t += 1
		}
		r.tailTrail[r.GetTail()] = true
	}
}

func (r *Rope) HandleRopeMoves(move string, steps int) {
	if steps > 0 {
		for i := 0; i < steps; i++ {
			r.HandleRopeMove(move)
		}
	}
}

func (r Rope) RopeTailShouldUpdate() bool {
	head, next := r.segments[0], r.segments[1]
	return GridLength(head, next) > 1
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

	r1 := NewRope(2)
	rope1 := &r1

	r2 := NewRope(10)
	rope2 := &r2

	for scanner.Scan() {
		text := scanner.Text()
		tokens := strings.Split(text, " ")
		move := strings.TrimSpace(tokens[0])
		steps, err := strconv.Atoi(tokens[1])
		if err == nil {
			handleRopeMove(rope1, move, steps)
			handleRopeMove(rope2, move, steps)
		}
	}

	fmt.Printf("count = %d \n", len((*rope1).tailTrail))
	fmt.Printf("count = %d \n", len((*rope2).tailTrail))

}
