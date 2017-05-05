package md

import (
	"fmt"
	"testing"
)

func TestLoadPatterns(t *testing.T) {
	pat := LoadPatterns("patterns.json")
	fmt.Printf("%+v\n", pat)
}
