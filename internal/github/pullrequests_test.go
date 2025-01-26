package github_test

import (
	"log"
	"testing"

	"github.com/lakrizz/prollama/internal/github"
)

func TestGetPullRequestsByRepositoryName(t *testing.T) {
	pr, err := github.GetPullRequestsByRepositoryName("lakrizz/hooks")
	if err != nil {
		log.Panic("error", err)
		return
	}
	log.Println(pr)
}
