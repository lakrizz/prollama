package main

import (
	"log"

	"github.com/lakrizz/prollama/internal/github"
	"github.com/lakrizz/prollama/pkg/diff"
	"github.com/lakrizz/prollama/pkg/ollama"
)

func main() {
	repo := "lakrizz/hooks"
	pullRequests, err := github.GetPullRequestsByRepositoryName(repo)
	if err != nil {
		panic(err)
	}

	for _, pr := range pullRequests {
		log.Println("now reviewing pr", pr.Title)
		diffs, err := diff.Parse(pr.Changes)
		if err != nil {
			panic(err)
		}

		comments, err := ollama.GetCommentsForPatch(diffs)
		if err != nil {
			panic(err)
		}

		// pp.Println(comments)

		err = github.AddCommentsToPR(repo, pr.Number, comments)
		if err != nil {
			return
		}
	}
}
