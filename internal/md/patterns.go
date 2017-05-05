package md

import (
	"encoding/json"
	"io/ioutil"
)

// Tag .
type Tag struct {
	Name       string `json:"name"`
	OpenWrap   string `json:"open_wrap"`
	CloseWrap  string `json:"close_wrap"`
	LinePrefix string `json:"line_prefix"`
}

// Patterns .
type Patterns struct {
	EmptyTags     []Tag `json:"empty_tags"`
	ContainerTags []Tag `json:"container_tags"`
}

func emptyPattern(tag string) string {
	return "<" + tag + "\\b([^>]*)\\/?>"
}

func containerPattern(tag string) string {
	return "<" + tag + "\\b([^>]*)>([\\s\\S]*?)<\\/" + tag + ">"
}

// LoadPatterns .
func LoadPatterns(filename string) *Patterns {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	pat := &Patterns{}
	json.Unmarshal(buf, pat)
	return pat
}
