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

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filename, ok := map[string]string{
		"book-of-mormon":         "../books/book-of-mormon.json",
		"new-testament":          "../books/new-testament.json",
		"old-testament":          "../books/old-testament.json",
		"doctrine-and-covenants": "../books/doctrine-and-covenants.json",
		"pearl-of-great-price":   "../books/pearl-of-great-price.json",
	}[q.Get("book")]
	if !ok {
		http.Error(w, "invalid choice of book: "+q.Get("book"), http.StatusBadRequest)
		return
	}

	breakdown, ok := map[string][]reading.Breakdown{
		"chapter": []reading.Breakdown{reading.Book, reading.Chapter},
		"verse":   []reading.Breakdown{reading.Reference},
	}[q.Get("breakdown")]
	if !ok {
		http.Error(w, "invalid choice of breakdown: "+q.Get("breakdown"), http.StatusBadRequest)
		return
	}

	days, err := strconv.ParseInt(q.Get("days"), 10, 0)
	if err != nil {
		http.Error(w, "could not parse days: "+err.Error(), http.StatusBadRequest)
		return
	}

	out := json.NewEncoder(w)
	err = out.Encode(reading.Plan(
		filename,
		int(days),
		breakdown...,
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
