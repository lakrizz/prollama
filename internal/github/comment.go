package github

import (
	"github.com/lakrizz/prollama/pkg/gh"
	"github.com/lakrizz/prollama/pkg/models"
)

func AddCommentsToPR(repo string, prNumber int, comments []*models.Comment) error {
	err := gh.AddReview(repo, prNumber, comments)
	if err != nil {
		return err
	}

	return nil
}
