package diff

import (
	"regexp"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/sourcegraph/go-diff/diff"
)

func Parse(input string) ([]*diff.FileDiff, error) {
	// first we need to split the big diff file into respective diffsets
	// each diffset is prefixed by two lines: one that strarts with 'diff --git' and one that starts with 'index'
	diffs := []*diff.FileDiff{}
	for _, unidiff := range splitUnidiff(input) {
		if strings.HasPrefix(unidiff, "Binary files") {
			continue
		}

		parsed, err := diff.ParseFileDiff([]byte(unidiff))
		if err != nil {
			pp.Println(unidiff)
			return nil, err
		}
		diffs = append(diffs, parsed)
	}
	return diffs, nil
}

// thanks chatgpt :D
func splitUnidiff(input string) []string {
	re := regexp.MustCompile(`(?m)^diff --git[\s\S]*?^index[^\n]*\n`)
	matches := re.FindAllStringIndex(input, -1)

	if matches == nil {
		return nil
	}

	unidiffs := []string{}

	for i, match := range matches {
		start := match[1]
		var end int
		if i+1 < len(matches) {
			end = matches[i+1][0]
		} else {
			end = len(input)
		}
		unidiffs = append(unidiffs, input[start:end])
	}

	return unidiffs
}
