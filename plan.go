package reading

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	Section      json.Number
	Subsubtitle  string
	Subtitle     string
	Version      json.Number
}

// Breakdown will group the reading by a category.
type Breakdown func(v verse) string

func Book(v verse) string         { return v.Book }
func Chapter(v verse) string      { return v.Chapter.String() }
func FullTitle(v verse) string    { return v.FullTitle }
func LDSSlug(v verse) string      { return v.LDSSlug }
func Reference(v verse) string    { return v.Reference }
func Text(v verse) string         { return v.Text }
func Verse(v verse) string        { return v.Verse.String() }
func LastModified(v verse) string { return v.LastModified }
func Section(v verse) string      { return v.Section.String() }
func Subsubtitle(v verse) string  { return v.Subsubtitle }
func Subtitle(v verse) string     { return v.Subtitle }
func Version(v verse) string      { return v.Version.String() }

// A Unit is the lowest unit of reading.
type Unit struct {
	Title  string `json:"title"`
	Weight float64
	Prev   int
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func toByteArray(f float64) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(b, binary.LittleEndian, f)
	return buf.Bytes()
}

func hashToString(units []Unit) string {
	h := md5.New()
	for _, u := range units {
		h.Write([]byte(u.Title))
		h.Write(toByteArray(u.Weight))
		h.Write([]byte{u.Prev})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func maybeLoadCache(filename string) [][]Unit {
	var cache [][]Unit
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(err, "while reading", filename)
		}
		return nil
	}
	err = json.Unmarshal(contents, &cache)
	if err != nil {
		log.Println(err, "while unmarshalling", filename)
		return nil
	}
	return cache
}

func maybeSaveCache(filename string, cache [][]Unit) {
	if len(cache) > 200 {
		cache = cache[:200]
	}
	contents, err := json.Marshal(cache)
	if err != nil {
		log.Println(err, "while marshalling cache")
		return
	}
	err = ioutil.WriteFile(filename, contents, 0664) // permission
	if err != nil {
		log.Println(err, "while saving", filename)
		return
	}
}

func sessionTitle(a, b Unit) string {
	if a.Title == b.Title {
		return "Read " + a.Title
	}
	return "Read " + a.Title + " - " + b.Title
}

func buildDynamicPlan(u []Unit, cache [][]Unit, days int) [][]Unit {
	plan := make([][]Unit, days+1)

	for i := 0; i < minInt(len(cache), days+1); i++ {
		plan[i] = cache[i]
	}
	for i := len(cache); i < days+1; i++ {
		plan[i] = make([]Unit, len(u))
	}
	for i := maxInt(len(cache), 1); i < days+1; i++ {
		for j := len(u) - 1; j >= 0; j-- {
			plan[i][j].Weight = u[j].Weight * u[j].Weight // initialize to reading everything.
			plan[i][j].Title = sessionTitle(u[0], u[j])

			for k := j - 1; k >= 0; k-- {
				words := u[j].Weight - u[k].Weight
				weight := ((words * words) + (float64(i-1) * plan[i-1][k].Weight)) / float64(i) // sum(words ^ 2) / n
				if weight > plan[i-1][j].Weight {
					break
				}
				if weight <= plan[i][j].Weight {
					plan[i][j] = Unit{
						Title:  sessionTitle(u[k+1], u[j]),
						Weight: weight,
						Prev:   k,
					}
				}
			}
		}
	}
	return plan
}

func addTitles(u []Unit) []Unit {
	for i := range u {
		u[i].Title = sessionTitle(u[i], u[i])
	}
	return u
}

// The Algorithm used to minimize the cost.
type Algorithm func([]Unit, int) []Unit

// Dynamic will try to minimize the maximum session word count.
func Dynamic(u []Unit, days int) []Unit {
	if days >= len(u) {
		return addTitles(u)
	}
	filename := "../cache/" + hashToString(u) + ".json"
	cache := maybeLoadCache(filename)
	w := days + 1
	newInfo := w > len(cache)
	plan := buildDynamicPlan(u, cache, days)
	if newInfo {
		maybeSaveCache(filename, plan)
	}

	result := make([]Unit, days)
	j := len(u) - 1
	for i := days; i >= 1; i-- {
		result[i-1] = plan[i][j]
		j = result[i-1].Prev
	}
	return result

}

// Greedy will read until the day's portion has been fulfilled.
func Greedy(u []Unit, days int) []Unit {
	if days >= len(u) {
		return addTitles(u)
	}

	total := u[len(u)-1].Weight
	result := make([]Unit, days)

	var curr int
	var prev int
	for i := 0; i < days; i++ {
		for ; u[curr].Weight < total*(float64(i+1)/float64(days)); curr++ {
		}
		result[i] = Unit{
			Title: sessionTitle(u[prev], u[curr]),
		}
		prev = curr
	}
	return result
}

func breakdown(verses []verse, breakdowns ...Breakdown) []Unit {
	var result []Unit
	curr := Unit{}
	var running float64

	for _, v := range verses {
		var titles []string
		for _, b := range breakdowns {
			titles = append(titles, b(v))
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

// loads a book to read into units.
func load(filename string) []verse {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err, "while reading", filename)
		return nil
	}
	var verses []verse
	err = json.Unmarshal(f, &verses)
	if err != nil {
		log.Println(err, "while unmarshalling", filename)
		return nil
	}
	return verses
}

// Plan will create a reading plan.
func Plan(filename string, days int, a Algorithm, breakdowns ...Breakdown) []Unit {
	if filename == "" || a == nil || days < 0 {
		return nil
	}
	v := load(filename)
	u := breakdown(v, breakdowns...)
	return a(u, days)
}
