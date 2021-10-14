package quizgame

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func QuizGame() {
	fileName := flag.String("file", "problems.csv", "quiz Q&A csv file")
	timeLimit := flag.Int("limit", 30, "quiz time limit (s)")
	flag.Parse()

	f := "quizgame/" + *fileName

	file, err := os.Open(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	problems := parseRecords(records)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	for i, p := range problems {
		fmt.Printf("Question #%d: %s = ", i+1, p.q)
		respCh := make(chan string)
		go func() {
			var resp string
			fmt.Scanf("%s", &resp)
			respCh <- resp
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nRan out of time\n")
			return
		case resp := <-respCh:
			if resp == p.a {
				correct++
			}
		}
	}
	fmt.Printf("Score: %d / %d\n", correct, len(problems))
}

func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))
	for i, record := range records {
		problems[i] = problem{
			q: record[0],
			a: record[1],
		}
	}
	return problems
}
