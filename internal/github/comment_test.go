package github_test

import (
	"log"
	"testing"

	"github.com/lakrizz/prollama/internal/github"
	"github.com/lakrizz/prollama/pkg/models"
)

func TestCommenting(t *testing.T) {
	cmts := []*models.Comment{
		&models.Comment{
			FileName: "backend/cmd/hooksim/main.go",
			Position:     45,
			Comment:  "Typo in error message: 'failde' should be 'failed'.",
		},
		&models.Comment{
			FileName: "backend/cmd/hooksim/main.go",
			Position:     60,
			Comment:  "Variable 'foo' is declared but not used. Consider removing it or using it.",
		},
		&models.Comment{
			FileName: "backend/cmd/hooksim/main.go",
			Position:     61,
			Comment:  "Variable 'bar' is declared but not used. Consider removing it or using it.",
		},
		&models.Comment{
			FileName: "backend/cmd/hooksim/main.go",
			Position:     70,
			Comment:  "Typo in hostname: 'lcoalhost' should be 'localhost'.",
		},
	}

	err := github.AddCommentsToPR("lakrizz/hooks", 76, cmts)
	if err != nil {
		log.Println("ERROR", err)
		panic(err)
	}
}
