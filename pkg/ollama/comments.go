package ollama

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jonathanhecl/gollama"
	"github.com/k0kubun/pp"
	"github.com/sourcegraph/go-diff/diff"

	"github.com/lakrizz/prollama/pkg/models"
)

func (o *Ollama) GenerateCommentsForPatch(diffs []*diff.FileDiff) ([]*models.Comment, error) {
	ctx := context.Background()

	res := make([]*models.Comment, 0)
	g := gollama.New(o.cfg.Model)
	g.ContextLength = 32768
	g.SystemPrompt = o.systemPrompt

	for _, unidiff := range diffs {
		if !checkFilenameAllowed(unidiff.NewName) {
			continue
		}

		// now check all the hunks
		for _, hunk := range unidiff.Hunks {
			output, err := g.Chat(ctx, fmt.Sprintf(o.userPrompt, string(hunk.Body)))
			if err != nil {
				return nil, err
			}

			// convert to our struct
			cmt := &models.Comment{
				FileName:  unidiff.NewName[2:],
				StartLine: int(hunk.NewStartLine),
				Comment:   output.Content,
			}

			pp.Println(cmt)
			res = append(res, cmt)
		}

	}

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
