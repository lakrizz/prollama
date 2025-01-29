package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lakrizz/prollama/internal/github"
	"github.com/lakrizz/prollama/pkg/hunks"
	"github.com/lakrizz/prollama/pkg/ollama"
)

func Prollama(ctx context.Context) error {
	o, err := ollama.New(ollama.WithContext(ctx))
	if err != nil {
		return err
	}

	gh, err := github.New(github.WithContext(ctx))
	if err != nil {
		return err
	}

	pullRequests, err := gh.GetPullRequests()
	if err != nil {
		return err
	}

	if len(pullRequests) == 0 {
		slog.Info("No Pull Requests found")
		return nil
	}

	for _, pr := range pullRequests {
		slog.Info(fmt.Sprintf("Reviewing PR '%v'", pr.Title), "id", pr.Number)

		diffs, err := hunks.Parse(pr.Changes)
		if err != nil {
			panic(err)
		}

		comments, err := o.GenerateCommentsForPatch(diffs)
		if err != nil {
			panic(err)
		}

		err = gh.AddCommentsToPR(pr.Number, comments)
		if err != nil {
			return err
		}
	}
	return nil
}
