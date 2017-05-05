package md

import (
	"fmt"
	"regexp"
)

func replace(re *regexp.Regexp, content string, t Tag) string {
	content = re.ReplaceAllStringFunc(content, func(html string) string {
		matches := re.FindStringSubmatch(html)
		if len(matches) > 2 {
			inner := matches[2]
			if t.LinePrefix != "" {
				inner = regexp.MustCompile(`(?m)^`).ReplaceAllString(inner, t.LinePrefix)
				fmt.Println(inner)
			}
			return t.OpenWrap + inner + t.CloseWrap
		}
		return t.OpenWrap + t.CloseWrap
	})
	return content
}

// Replace .
func Replace(content string, pat *Patterns) string {
	var re *regexp.Regexp
	for _, t := range pat.EmptyTags {
		re = regexp.MustCompile(emptyPattern(t.Name))
		content = replace(re, content, t)
	}
	for _, t := range pat.ContainerTags {
		re = regexp.MustCompile(containerPattern(t.Name))
		content = replace(re, content, t)
	}
	return content
}
