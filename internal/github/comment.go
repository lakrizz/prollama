package github

import (
	"github.com/lakrizz/prollama/pkg/gh"
	"github.com/lakrizz/prollama/pkg/models"
)

func (s *Service) AddCommentsToPR(prNumber int, comments []*models.Comment) error {
	err := gh.AddReview(s.Config.Repo, prNumber, comments)
	if err != nil {
		return err
	}

	return nil
}
