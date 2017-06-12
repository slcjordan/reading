package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/slcjordan/reading"
)

var info = log.New(os.Stdout, "", log.LstdFlags)

func filename(r *http.Request) string {
	return map[string]string{
		"book-of-mormon":         "../books/book-of-mormon.json",
		"new-testament":          "../books/new-testament.json",
		"old-testament":          "../books/old-testament.json",
		"doctrine-and-covenants": "../books/doctrine-and-covenants.json",
		"pearl-of-great-price":   "../books/pearl-of-great-price.json",
	}[r.URL.Query().Get("book")]
}

func breakdowns(r *http.Request) []reading.Breakdown {
	return map[string][]reading.Breakdown{
		"chapter": []reading.Breakdown{reading.Book, reading.Chapter},
		"verse":   []reading.Breakdown{reading.Reference},
	}[r.URL.Query().Get("breakdown")]
}

func days(r *http.Request) int {
	days, err := strconv.ParseInt(r.URL.Query().Get("days"), 10, 0)
	if err != nil {
		return 0
	}
	return int(days)
}

func algorithm(r *http.Request) reading.Algorithm {
	return map[string]reading.Algorithm{
		"chapter": reading.Dynamic,
		"verse":   reading.Greedy,
	}[r.URL.Query().Get("breakdown")]
}

func handler(w http.ResponseWriter, r *http.Request) {
	out := json.NewEncoder(w)
	err := out.Encode(
		reading.Plan(
			filename(r),
			days(r),
			algorithm(r),
			breakdowns(r)...,
		))

	if err != nil {
		info.Println(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("site"))
	http.HandleFunc("/plan", handler)
	http.Handle("/", fs)
	addr := ":8080"
	info.Println("serving at " + addr)
	http.ListenAndServe(addr, nil)
}
