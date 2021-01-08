package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func quiz(records [][]string, timer *time.Timer) (int, int) {
	var hits, responses int = 0, 0

	for _, record := range records {
		question := record[0]
		answer := record[1]

		fmt.Printf("%s = ", question)

		answerCh := make(chan string)
		go func() {
			var resp string
			fmt.Scanf("%s\n", &resp)
			answerCh <- resp
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime is up.")
			return hits, responses
		case resp := <-answerCh:
			if resp == answer {
				hits++
			}
			responses++
		}
	}
	return hits, responses
}

func main() {
	f, err := os.Open("problems.csv")
	check(err)
	r := io.Reader(f)
	reader := csv.NewReader(r)

	records, err := reader.ReadAll()
	check(err)

	total := len(records)

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	hits, responses := quiz(records, timer)

	fmt.Printf("%d answered, %d correctly of %d total\n", responses, hits, total)
}
