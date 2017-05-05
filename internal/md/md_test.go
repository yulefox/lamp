package md

import (
	"io/ioutil"
	"testing"
)

func load(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return buf
}

func TestPattern(t *testing.T) {
	html := string(load("test.html"))
	pat := LoadPatterns("patterns.json")
	Replace(html, pat)
}
