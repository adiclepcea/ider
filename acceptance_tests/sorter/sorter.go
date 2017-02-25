package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type intArray []int64

func (s intArray) Len() int           { return len(s) }
func (s intArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s intArray) Less(i, j int) bool { return s[i] < s[j] }

func main() {
	ids := make([]int64, 50000000)
	var pos int64
	for i := 0; i < 50; i++ {
		fmt.Println(i + 1)
		file, err := os.Open(fmt.Sprintf("/data/%d.txt", i+1))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			id, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				panic(err)
			}
			ids[pos] = id
			pos++
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	fmt.Println("Sorting")

	sort.Sort(intArray(ids))

	fmt.Println("Writing result to file")

	fileOut, err := os.OpenFile("/data/out.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	defer fileOut.Close()

	for i := 0; i < 50000000; i++ {
		fileOut.WriteString(fmt.Sprintf("%d\n", ids[i]))
	}
}
