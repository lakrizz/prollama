package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jonathanhecl/gollama"
	"github.com/k0kubun/pp"
	"github.com/sourcegraph/go-diff/diff"

	"github.com/lakrizz/prollama/pkg/hunks"
	"github.com/lakrizz/prollama/pkg/models"
)

func (o *Ollama) GenerateCommentsForPatch(diffs []*diff.FileDiff) ([]*models.Comment, error) {
	logFolder := filepath.Join("logs", fmt.Sprintf("%v", time.Now().Unix()))
	ctx := context.Background()

	res := make([]*models.Comment, 0)
	g := gollama.New(o.cfg.Model)
	g.ContextLength = 32768
	g.SystemPrompt = o.systemPrompt

	for _, unidiff := range diffs {
		if !checkFilenameAllowed(unidiff.NewName) {
			continue
		}

		log.Println("checking new diff", unidiff.NewName)

		// now check all the hunks
		for i, hunk := range unidiff.Hunks {
			log.Println("checking hunk", i+1, "of", len(unidiff.Hunks))
			output, err := g.Chat(ctx, fmt.Sprintf(o.userPrompt, string(hunk.Body)))
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
				affectedLineNumber := hunks.FindLineNumberInHunk(c.AffectedLine, string(hunk.Body))
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
