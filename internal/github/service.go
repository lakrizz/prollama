package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lakrizz/prollama/config"
	"github.com/lakrizz/prollama/pkg/gh"
)

type Service struct {
	Config *config.Config
}

type GitHubServiceOption func(*Service) error

func New(opts ...GitHubServiceOption) (*Service, error) {
	s := &Service{}

	for _, o := range opts {
		err := o(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Config.Repo == "" {
		if s.Config.Debug {
			slog.Debug("no repo given, trying to figure it out")
		}

		err := s.AutoSelectRepo()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func WithContext(ctx context.Context) GitHubServiceOption {
	return func(s *Service) error {
		if cfg, ok := config.FromContext(ctx); ok {
			s.Config = cfg
			return nil
		}

		return fmt.Errorf("context is missing config")
	}
}

func (s *Service) AutoSelectRepo() error {
	repo, err := gh.GetRepositoryName()
	if err != nil {
		return err
	}

	s.Config.Repo = repo
	if s.Config.Debug {
		slog.Debug("found repository", "repository_name", repo)
	}
	return nil
}
