package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	cli := Cli{}
	cli.Run()
}

type Cli struct {}

type Problem struct {
	question string
	answer string
}

type Quiz struct {
	Problems []*Problem
}

func validateArgs() {
	if len(os.Args) < 2 {
		Usages()
		os.Exit(0)
	}
}

func Usages(){
	fmt.Println("Usage:")
	fmt.Println("-csv `file.csv` -limit `time`")
}

func (cli *Cli) Run() error {
	validateArgs()

	csvFile := flag.String("csv", "problems.csv", "file name")
	limit := flag.Int("timelimit", 30, "time limit")

	flag.Parse()
	if _, err := os.Stat(*csvFile); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(*csvFile)
	defer file.Close()

	if err != nil {
		log.Fatalf("Failed to open the CSV file: %s", *file)
	}
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("Failed to open csv file")
	}

	quiz := getProblems(records)
	fmt.Println("YOU HAVE 10 QUESTIONS\nANSWER AS FAST AS YOU CAN\nSTART...")
	c := start(quiz,limit)
	fmt.Printf("You scored %d correct out of %d\n", c,len(quiz.Problems))
	return nil
}

func getProblems(records [][]string) *Quiz {
	quiz := &Quiz{}

	for _, value := range records {
		q := value[0]
		a := value[1]
		p := &Problem{q, a }
		quiz.Problems = append(quiz.Problems, p)
	}
	return quiz
}

func start(quiz *Quiz,limit *int) int{
	var correct int
	timer := time.NewTimer(time.Duration(*limit) * time.Second)

	QuizLoop:
		for index, value := range quiz.Problems {
			ansCh := make(chan string)
			fmt.Printf("Problem #%d: %s = ", index+1, value.question)

			go func() {
				var ans string
				fmt.Scanf("%s\n", &ans)
				ansCh <- ans
			}()

			select {
			case <- timer.C:
				fmt.Println()
				break QuizLoop
			case ans := <- ansCh:
				if ans == value.answer {
					correct++
				}
			}
		}
	return correct
}




