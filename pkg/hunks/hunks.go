package hunks

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

// findLineNumberInHunk searches for the line number of a specific line (needle) within a block of text (haystack).
// It splits the haystack into lines and compares each line with the needle, ignoring control characters.
// If a match is found, it returns the line number (0-based index). If no match is found, it returns -1.
func FindLineNumberInHunk(needle string, haystack string) int {
	hs := strings.Split(haystack, "\n")
	for i, line := range hs {
		if compareIgnoringControlChars(line, needle) {
			return i - 1
		}
	}
	return -1
}

// compareIgnoringControlChars compares two strings for equality while ignoring ASCII control characters and the DEL character.
// It uses a regular expression to remove these non-printable characters from both input strings before performing the comparison.
func compareIgnoringControlChars(first, second string) bool {
	// This regular expression matches any byte with a value from 0 to 31 (inclusive) or 127.
	// These bytes are part of the ASCII control characters and the DEL character, which are not printable characters.
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	normalizedStr1 := re.ReplaceAllString(first, "")
	normalizedStr2 := re.ReplaceAllString(second, "")
	return normalizedStr1 == normalizedStr2
}
