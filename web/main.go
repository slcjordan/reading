package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	fmt.Printf("%v %v %v %v = %v", filename(r), days(r), breakdowns(r), algorithm(r),
		reading.Plan(
			filename(r),
			days(r),
			algorithm(r),
			breakdowns(r)...,
		))
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
	var addr string
	var static string
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "the address to serve from")
	flag.StringVar(&static, "static", "dist", "the static directory")
	flag.Parse()
	defer os.RemoveAll(reading.CacheDirectory)

	r := mux.NewRouter()
	r.HandleFunc("/plan", handler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(static))))

	info.Println("serving at " + addr)
	http.ListenAndServe(addr,
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
		)(
			handlers.LoggingHandler(os.Stdout,
				handlers.CompressHandler(r),
			)))
}
