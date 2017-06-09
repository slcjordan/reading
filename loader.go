package reading

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/slcjordan/reading/plan"
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

type Breakdown func(verse) string

func Book(v verse) string         { return v.Book }
func Chapter(v verse) string      { return v.Chapter.String() }
func FullTitle(v verse) string    { return v.FullTitle }
func LDSSlug(v verse) string      { return v.LDSSlug }
func Reference(v verse) string    { return v.Reference }
func Text(v verse) string         { return v.Text }
func Verse(v verse) string        { return v.Verse.String() }
func LastModified(v verse) string { return v.LastModified }
func Section(v verse) string      { return v.Section }
func Subsubtitle(v verse) string  { return v.Subsubtitle }
func Subtitle(v verse) string     { return v.Subtitle }
func Version(v verse) string      { return v.Version }

func Load(filename string, breakdowns ...Breakdown) []plan.Unit {
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
	var result []plan.Unit
	curr := plan.Unit{}

	for _, v := range verses {
		var names []string
		for _, b := range breakdowns {
			names = append(names, b(v))
		}
		name := strings.Join(names, " ")
		if curr.Name != name && curr.Weight != 0 {
			result = append(result, curr)
			curr = plan.Unit{}
		}
		curr.Name = name
		curr.Weight += int64(len(strings.Fields(v.Text)))
	}
	if curr.Weight != 0 {
		result = append(result, curr)
	}
	return result
}
