package reading

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type verse struct {
	Book         string
	Chapter      json.Number
	FullTitle    string `json:"full_title"`
	LDSSlug      string `json:"lds_slug"`
	Reference    string
	Text         string
	Verse        json.Number
	LastModified string `json:"last_modified"`
	Section      string
	Subsubtitle  string
	Subtitle     string
	Version      string
}

// Breakdown will group the reading by a category.
type Breakdown struct {
	Name   string
	Select func(v verse) string
}

var Book = Breakdown{Name: "Book", Select: func(v verse) string { return v.Book }}
var Chapter = Breakdown{Name: "Chapter", Select: func(v verse) string { return v.Chapter.String() }}
var FullTitle = Breakdown{Name: "FullTitle", Select: func(v verse) string { return v.FullTitle }}
var LDSSlug = Breakdown{Name: "LDSSlug", Select: func(v verse) string { return v.LDSSlug }}
var Reference = Breakdown{Name: "Reference", Select: func(v verse) string { return v.Reference }}
var Text = Breakdown{Name: "Text", Select: func(v verse) string { return v.Text }}
var Verse = Breakdown{Name: "Verse", Select: func(v verse) string { return v.Verse.String() }}
var LastModified = Breakdown{Name: "LastModified", Select: func(v verse) string { return v.LastModified }}
var Section = Breakdown{Name: "Section", Select: func(v verse) string { return v.Section }}
var Subsubtitle = Breakdown{Name: "Subsubtitle", Select: func(v verse) string { return v.Subsubtitle }}
var Subtitle = Breakdown{Name: "Subtitle", Select: func(v verse) string { return v.Subtitle }}
var Version = Breakdown{Name: "Version", Select: func(v verse) string { return v.Version }}

// A Unit is the lowest unit of reading.
type Unit struct {
	Title  string `json:"title"`
	Weight float64
	Prev   int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func sessionTitle(a, b Unit) string {
	if a.Title == b.Title {
		return "Read " + a.Title
	}
	return "Read " + a.Title + " - " + b.Title
}

// Plan will create a reading plan.
func Plan(filename string, days int, breakdowns ...Breakdown) []Unit {
	u := load(filename, breakdowns...)
	var parts []string
	for _, b := range breakdowns {
		parts = append(parts, b.Name)
	}
	cachefile := "../cache/" + strings.Join(parts, "_") + "-" + path.Base(filename)
	var cache [][]Unit
	contents, err := ioutil.ReadFile(cachefile)
	if err != nil {
		/* do nothing */
	} else {
		_ = json.Unmarshal(contents, &cache)
	}
	curr := initialize(u, days)
	for i := 1; i < min(len(cache), len(curr)); i++ {
		curr[i] = cache[i]
	}

	result, curr := create(
		u,
		curr,
		maxint(len(cache), 1),
		days,
	)
	if len(curr) > len(cache) {
		contents, err := json.Marshal(curr)
		if err != nil {
			/* do nothing */
		} else {
			_ = ioutil.WriteFile(cachefile, contents, 0664) // permission
		}
	}
	return result
}

func initialize(u []Unit, n int) [][]Unit {
	h, w := len(u), n
	d := make([][]Unit, w)

	for i := 0; i < n; i++ {
		d[i] = make([]Unit, h)
	}

	for j := 0; j < h; j++ {
		d[0][j] = Unit{
			Title:  sessionTitle(u[0], u[j]),
			Weight: u[j].Weight * u[j].Weight,
			Prev:   -1,
		}
	}
	return d
}

func create(u []Unit, d [][]Unit, idx, n int) ([]Unit, [][]Unit) {
	h := len(u)
	for i := idx; i < n; i++ {
		for j := h - 1; j >= 0; j-- {
			readnothing := d[i-1][j].Weight
			d[i][j].Weight = readnothing
			d[i][j].Prev = j

			for k := j - 1; k >= 0; k-- {
				last := d[i-1][k].Weight
				words := u[j].Weight - u[k].Weight
				weight := ((words * words) + (float64(i) * last)) / (float64(i) + 1) // sum(words ^ 2) / n
				if weight > readnothing {
					break
				}
				if weight <= d[i][j].Weight {
					d[i][j] = Unit{
						Title:  sessionTitle(u[k+1], u[j]),
						Weight: weight,
						Prev:   k,
					}
				}
			}
		}
	}

	result := make([]Unit, n)
	j := h - 1
	for i := n - 1; i >= 0 && j >= 0; i-- {
		result[i] = d[i][j]
		j = result[i].Prev
	}
	return result, d
}

// loads a book to read into units.
func load(filename string, breakdowns ...Breakdown) []Unit {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	var verses []verse
	err = json.Unmarshal(f, &verses)
	if err != nil {
		log.Println(err)
		return nil
	}
	var result []Unit
	curr := Unit{}
	var running float64

	for _, v := range verses {
		var titles []string
		for _, b := range breakdowns {
			titles = append(titles, b.Select(v))
		}
		title := strings.Join(titles, " ")
		if curr.Title != title && curr.Weight != 0 {
			result = append(result, curr)
			curr = Unit{}
		}
		curr.Title = title
		running += float64(len(strings.Fields(v.Text)))
		curr.Weight = running
	}
	if curr.Weight != 0 {
		result = append(result, curr)
	}
	return result
}
