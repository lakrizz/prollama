package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/jonathanhecl/gollama"
	"github.com/k0kubun/pp"
	"github.com/sourcegraph/go-diff/diff"

	"github.com/lakrizz/prollama/pkg/models"
)

var Model = "qwen2.5-coder:14b"

const (
	gptPrompt = `
	Please review the following unified Diff file. Identify issues such as syntax errors, logical flaws, unhandled errors, code smells, and testing gaps. Suggest improvements, refactoring opportunities, or missing tests where applicable. Inform about State of the Art implementations, Best Practices and Design Patterns. Combine multiple findings for the same line into a single comment. 

Return all feedback as an array of JSON objects, where each object contains the fields:  
- 'line': The line number in the changed file, as indicated by the patch metadata.  
- 'body': A detailed explanation of the issue and actionable suggestions for improvement.  
- 'affected_line': a copy of the line this comment belongs to. Include all (this also applies to repeated instances) control characters and the leading '+' or '-'

If no issues are found, return an empty array ('[]').  

Ignore metadata lines that indicate information for the Diffpatch (e.g., lines that contain four numbers).
Assume all brackets, quotes, parantheses are closed at some point, so do not mark a missing closing or an unclosed pair as an error. 

This is the patch: %v`

	systemPrompt = "You are a principal software engineer with extensive experience in performing detailed code reviews, adhering to best practices, and ensuring code quality in production systems. Your role is to critically analyze Git patch files, identifying issues such as syntax errors, logical flaws, unhandled errors, code smells, and testing gaps. Provide actionable, well-reasoned feedback in JSON format with references to file names and line numbers, scoped strictly to the context of the patch. Focus on clarity, correctness, and optimization, ensuring that your feedback is concise, relevant, and insightful. Do not comment on correct or acceptable code."
)

func GetCommentsForPatch(diffs []*diff.FileDiff) ([]*models.Comment, error) {
	logFolder := filepath.Join("logs", fmt.Sprintf("%v", time.Now().Unix()))
	ctx := context.Background()

	res := make([]*models.Comment, 0)
	g := gollama.New(Model)
	g.ContextLength = 32768
	g.SystemPrompt = systemPrompt

	for _, unidiff := range diffs {
		if !checkFilenameAllowed(unidiff.NewName) {
			continue
		}

		log.Println("checking new diff", unidiff.NewName)

		// now check all the hunks
		for i, hunk := range unidiff.Hunks {
			log.Println("checking hunk", i+1, "of", len(unidiff.Hunks))
			output, err := g.Chat(ctx, fmt.Sprintf(gptPrompt, string(hunk.Body)))
			if err != nil {
				return nil, err
			}
			// convert to our struct
			cmt := []*models.Comment{}
			err = json.Unmarshal(removeJSONTag(output.Content), &cmt)
			if err != nil {
				log.Println(fmt.Errorf("error unmarshalling ollama response: %w", err))
				pp.Println(removeJSONTag(output.Content))
				err = os.WriteFile(filepath.Join(logFolder, fmt.Sprintf("F%vH%v", filepath.Base(unidiff.NewName), i)), []byte(output.Content), 0777)
				if err != nil {
					log.Println(err)
				}
				continue
			}

			// now fix all the lines etc
			for _, c := range cmt {
				affectedLineNumber := findLineNumberInHunk(c.AffectedLine, string(hunk.Body))
				if affectedLineNumber == -1 {
					continue
				}

				c.FileName = unidiff.NewName[2:]
				c.EndLine = int(hunk.NewStartLine) + affectedLineNumber
			}

			log.Println("adding", len(cmt), "new comments")
			res = append(res, cmt...)
		}
	}

	log.Println("found", len(res), "comments for this patch")
	return res, nil
}

// checkFilenameAllowed checks if the given filename is allowed based on specific rules.
//
// The function returns false if the filename has a blocked extension or contains any
// of the blocked partial strings. Otherwise, it returns true.
func checkFilenameAllowed(filename string) bool {
	blockedExtensions := []string{".yaml", ".toml", ".xml", ".json", ""}
	blockedPartials := []string{".min.", ".gen.", ".d."}

	ext := filepath.Ext(filename)
	if slices.Contains(blockedExtensions, ext) {
		return false
	}

	if slices.ContainsFunc(blockedPartials, func(s string) bool {
		return strings.Contains(filename, s)
	}) {
		return false
	}

	return true
}
func removeJSONTag(input string) []byte {
	j := strings.TrimPrefix(input, "```json")
	return []byte(strings.TrimSuffix(j, "```"))
}

// findLineNumberInHunk searches for the line number of a specific line (needle) within a block of text (haystack).
// It splits the haystack into lines and compares each line with the needle, ignoring control characters.
// If a match is found, it returns the line number (0-based index). If no match is found, it returns -1.
func findLineNumberInHunk(needle string, haystack string) int {
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
