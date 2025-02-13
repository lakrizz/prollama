package github

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lakrizz/prollama/pkg/gh"
)

// this package is used to retrieve all open pull requests for either a given repo
// or the base repository in `pwd`

type PullRequest struct {
	Repo string `json:"-"`

	Title   string `json:"title,omitempty"`
	Body    string `json:"body,omitempty"`
	State   string `json:"state,omitempty"`
	ID      string `json:"id,omitempty"`
	Number  int    `json:"number,omitempty"`
	Changes string `json:"-"`
}

func (s *Service) GetPullRequests() ([]*PullRequest, error) {
	dat, err := gh.GetPullRequests(s.Config.Repo)
	if err != nil {
		return nil, fmt.Errorf("error getting pull requests: %w", err)
	}

	pr := []*PullRequest{}
	err = json.Unmarshal(dat, &pr)
	if err != nil {
		return nil, fmt.Errorf("GetPullRequestsByRepositoryName: error unmarshalling response: %w", err)
	}

	for _, req := range pr {
		req.Repo = s.Config.Repo
		err := req.getDiffs()
		if err != nil {
			if s.Config.Debug {
				slog.Debug("error getting diff, not aborting!", "error", err.Error())
			}
			continue
		}
	}

	if s.Config.Debug {
		slog.Debug("pull requests", "amount", len(pr))
	}

	return pr, nil
}

func (pr *PullRequest) getDiffs() error {
	dat, err := gh.GetPullRequestDiff(pr.Repo, pr.Number)
	if err != nil {
		return err
	}

	pr.Changes = string(dat)
	return nil
}
