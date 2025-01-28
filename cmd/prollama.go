package cmd

import (
	"context"
	"log"

	"github.com/lakrizz/prollama/internal/github"
	"github.com/lakrizz/prollama/pkg/hunks"
	"github.com/lakrizz/prollama/pkg/ollama"
)

func Prollama(ctx context.Context) error {
	o, err := ollama.New(ollama.WithContext(ctx))
	if err != nil {
		return err
	}

	repo := "lakrizz/hooks"
	pullRequests, err := github.GetPullRequestsByRepositoryName(repo)
	if err != nil {
		panic(err)
	}

	for _, pr := range pullRequests {
		log.Println("now reviewing pr", pr.Title)
		diffs, err := hunks.Parse(pr.Changes)
		if err != nil {
			panic(err)
		}

		comments, err := o.GenerateCommentsForPatch(diffs)
		if err != nil {
			panic(err)
		}

		err = github.AddCommentsToPR(repo, pr.Number, comments)
		if err != nil {
			return err
		}
	}
	return nil
}
