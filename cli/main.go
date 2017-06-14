package main

import (
	"fmt"
	"os"

	"github.com/slcjordan/reading"
)

func main() {
	var days int
	fmt.Print("How many days to read d&c? ")
	_, err := fmt.Fscanf(os.Stdin, "%d", &days)
	if err != nil {
		fmt.Println(err)
		return
	}
	var idx int
	fmt.Print("How would you like that broken down [(1) Chapter (2) Verse]? ")
	_, err = fmt.Fscanf(os.Stdin, "%d", &idx)
	if err != nil {
		fmt.Println(err)
		return
	}
	if idx > 2 {
		idx = 2
	}
	if idx < 1 {
		idx = 1
	}
	fmt.Print("okay, let me plan that out for you....")
	b := [][]reading.Breakdown{
		[]reading.Breakdown{reading.Book, reading.Chapter},
		[]reading.Breakdown{reading.Reference},
	}[idx-1]
	p := reading.Plan(
		"../books/pearl-of-great-price.json",
		days,
		reading.Greedy,
		b...,
	)
	fmt.Printf("in order to read the d&c in %d days:\n", days)
	for i, session := range p {
		fmt.Printf("day %d: %s\n", i+1, session.Title)
	}
}
