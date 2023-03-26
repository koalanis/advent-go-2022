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
var day = 7
var path = os.Getenv("ADVENT_HOME")

type File struct {
	name string
	size int
}

type Dir struct {
	name  string
	files map[string]*File
	dirs  map[string]*Dir
	prev  *Dir
	size  int
}

type FS struct {
	root *Dir
}

type Command struct {
	text   string
	output []string
}

func (d Dir) String() string {
	format := "[ name=%s size=%d files=%d dirs=%s]"
	return fmt.Sprintf(format, d.name, d.size, len(d.files), fmt.Sprintln(d.dirs))
}

func (f File) String() string {
	return fmt.Sprintf("<%s:%d>\n", f.name, f.size)
}

func (c Command) String() string {
	return fmt.Sprintf("text=%s output=[%s]\n", c.text, strings.Join(c.output, ","))
}

func handleCommand(fs *FS, currentDirectory *Dir, command Command) *Dir {
	if strings.HasPrefix(command.text, "cd") {
		tokens := strings.Split(command.text, " ")
		dir := tokens[1]
		if dir == "/" {
			currentDirectory = (*fs).root
		} else if dir == ".." {
			if currentDirectory.prev != nil {
				currentDirectory = currentDirectory.prev
			}
		} else {
			// create directory
			nextDir := &Dir{dir, make(map[string]*File), make(map[string]*Dir), currentDirectory, 0}
			currentDirectory.dirs[dir] = nextDir
			currentDirectory = nextDir
		}
	} else if strings.HasPrefix(command.text, "ls") {
		for _, str := range command.output {
			if !strings.HasPrefix(str, "dir") {
				tokens := strings.Split(str, " ")
				size, _ := strconv.Atoi(tokens[0])
				filename := tokens[1]
				file := &File{filename, size}
				currentDirectory.files[filename] = file
			}
		}
	}
	return currentDirectory
}

func computeSize(dir *Dir, sizeList *[]int) int {
	if dir == nil {
		return 0
	}

	fileSize := 0
	for _, file := range (*dir).files {
		if file != nil {
			fileSize += (*file).size
		}
	}

	dirSize := 0
	for _, d := range (*dir).dirs {
		if d != nil {
			dirSize += computeSize(d, sizeList)
		}
	}
	size := fileSize + dirSize
	(*dir).size = size
	*sizeList = append(*sizeList, size)
	return size
}

func main() {

	dataPath := fmt.Sprintf("%s/%d/data/day%d/data.txt", path, year, day)

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

	cmds := []*Command{}

	var currCommand *Command = nil
	for _, str := range data {
		if strings.HasPrefix(str, "$") {
			currCommand = &Command{}
			currCommand.text = strings.TrimSpace(str[1:])
			currCommand.output = make([]string, 0)
			cmds = append(cmds, currCommand)
		} else {
			currCommand.output = append(currCommand.output, str)
		}
	}

	// create root folder
	fs := &FS{}
	fs.root = &Dir{"/", make(map[string]*File), make(map[string]*Dir), nil, 0}

	var currentDirectory *Dir

	for _, cmd := range cmds {
		currentDirectory = handleCommand(fs, currentDirectory, *cmd)
	}

	sizeList := &[]int{}
	computeSize(fs.root, sizeList)

	cap := 100000
	sum := 0
	maxSize := 0
	for _, i := range *sizeList {

		// find max
		if maxSize < i {
			maxSize = i
		}

		if i <= cap {
			sum += i
		}
	}
	fmt.Printf("totalSize = %d, maxSize = %d\n", sum, maxSize)

	//  part 2
	totalDiskSpace := 70000000
	spaceNeededForUpdate := 30000000
	unusedSpace := totalDiskSpace - maxSize
	spaceStillNeeded := spaceNeededForUpdate - unusedSpace

	minSizeCanDelete := maxSize
	for _, i := range *sizeList {
		if i >= spaceStillNeeded {
			if i <= minSizeCanDelete {
				minSizeCanDelete = i
			}
		}
	}

	fmt.Printf("minSizeCanDelete = %d\n", minSizeCanDelete)

}
