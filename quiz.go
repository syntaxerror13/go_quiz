package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "flag, if present, the quiz will be shuffled before run")
	flag.Parse()

	f, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Println("Cannot open CSV file")
		os.Exit(1)
	}
	r := io.Reader(f)
	reader := csv.NewReader(r)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Cannot read CSV file")
		os.Exit(1)
	}

	total := len(records)

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(total, func(i, j int) { records[i], records[j] = records[j], records[i] })
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	defer timer.Stop()

	hits, responses := quiz(records, timer)

	fmt.Printf("%d answered, %d correctly of %d total\n", responses, hits, total)
}
