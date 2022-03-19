package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Prompt   string
	Solution string
}

var correctCounter, totalCounter int

func main() {
	var problemsFile = flag.String("f", "problems.csv", "CSV file with problems")
	var shuffle = flag.Bool("r", false, "Randomly shuffle the problems")
	var secondsToSolve = flag.Int("t", 30, "Time the player has to solve the problems")
	flag.Parse()

	problems := readProblems(*problemsFile)
	correctCounter, totalCounter = 0, 0

	go func() {
		time.Sleep(time.Duration(*secondsToSolve) * time.Second)
		fmt.Printf("\nTime's up! You correctly solved %v out of %v problems! Well done!\n",
			correctCounter, totalCounter)
		os.Exit(0)
	}()

	fmt.Println("----------- Math Quiz ------------")
	fmt.Printf("---- You have %v seconds. Go! ----\n", *secondsToSolve)
	for {
		if *shuffle {
			problems = shuffleProblems(problems)
		}
		runProblems(problems)
	}
}

func readProblems(filename string) []Problem {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	problemsCSV, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	problems := make([]Problem, 0, len(problemsCSV))
	for _, p := range problemsCSV {
		problems = append(problems, Problem{p[0], p[1]})
	}

	return problems
}

func shuffleProblems(problems []Problem) []Problem {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems),
		func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	return problems
}

func runProblems(problems []Problem) {
	reader := bufio.NewReader(os.Stdin)
	for _, p := range problems {
		fmt.Printf("Question %v: %v = ", totalCounter+1, p.Prompt)
		solution := readFromInput(reader)

		if solution == p.Solution {
			fmt.Println("Correct!")
			correctCounter++
		} else {
			fmt.Println("Wrong!")
		}
		totalCounter++
	}
}

func readFromInput(reader *bufio.Reader) string {
	solution, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(solution, "\n", "", -1)
}
