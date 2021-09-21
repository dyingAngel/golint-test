package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type question struct {
	question string
	answer   string
}

func loadQuestions() ([]question, error) {
	// Load csv and put in array
	csvFile, err := os.Open("problems.csv")
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	var questions []question
	for _, line := range csvLines {
		q := question{
			question: line[0],
			answer:   line[1],
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// Print Question and wait for answer
func getQuestionAnswer(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		// result = strings.Replace(result, "\n", "", -1)
		result = strings.ReplaceAll(result, "\n", "")
		if err != nil {
			fmt.Println(err)
		}

		input <- result
	}
}

func main() {
	// Time Limit
	timeLimit := 5
	timeEnd := time.Duration(timeLimit)
	counter := 0
	correct := 0

	questions, err := loadQuestions()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Sequence 1 : Welcome and Read Enter
	fmt.Println("------------------------------------------------")
	fmt.Println("Welcome to The Quiz")
	fmt.Println("You will have 30 second to answer the questions")
	fmt.Println("------------------------------------------------")
	fmt.Println("Please Press Enter To Start The Quiz")

	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err := reader.ReadRune()
		if char == 10 {
			break
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Start the quiz with 30s time limit.
	fmt.Println("---------  START !!! -----------")

	input := make(chan string)
	total := len(questions)

	go getQuestionAnswer(input)

	fmt.Println(questions[counter].question)
	for {
		select {
		// Input Received
		case i := <-input:
			// Validate Input
			if strings.Compare(questions[counter].answer, i) == 0 {
				correct++
			}

			counter++
			// Out of Questions
			if counter >= total {
				fmt.Println("-----------  ALL DONE -------------")
				fmt.Printf("Correct Answer : %v out of %v questions \n", correct, total)
				return
			}

			// If not end, print the next question.
			fmt.Println(questions[counter].question)
		// Timeout
		case <-time.After(timeEnd * time.Second):
			fmt.Println("-----------  TIME UP -------------")
			fmt.Printf("Correct Answer : %v out of %v questions \n", correct, total)
			return
		}
	}
}
