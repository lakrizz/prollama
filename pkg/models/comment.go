package models

type Comment struct {
	FileName     string `json:"path,omitempty"`
	StartLine    int    `json:"startLine,omitempty"`
	EndLine      int    `json:"line,omitempty"`
	AffectedLine string `json:"affected_line,omitempty"`
	Comment      string `json:"body,omitempty"`
}

type GithubComment struct {
	FileName  string `json:"path,omitempty"`
	StartLine int    `json:"startLine,omitempty"` // this will be taken from the hunk
	EndLine   int    `json:"line,omitempty"`      // this will be taken from the hunk
	Comment   string `json:"body,omitempty"`
	Side      string `json:"side,omitempty"`
}

func (c *Comment) GetGitHubCompatible() *GithubComment {
	return &GithubComment{
		FileName:  c.FileName,
		StartLine: c.StartLine,
		EndLine:   c.EndLine,
		Comment:   c.Comment,
		Side:      "RIGHT",
	}
}
