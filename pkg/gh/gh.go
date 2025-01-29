package gh

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/k0kubun/pp"

	"github.com/lakrizz/prollama/pkg/models"
)

type pullrequest struct {
	ID int `json:"id,omitempty"`
}

func execute(input []string) ([]byte, error) {
	cmd := exec.Command("gh", input...)
	return cmd.CombinedOutput()
}

func GetPullRequests(repo string) ([]byte, error) {
	args := []string{
		"pr",
		"list",
		"--json", "title,body,id,state,number",
		"--repo", repo,
	}

	res, err := execute(args)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetPullRequestDiff(repo string, number int) ([]byte, error) {
	args := []string{
		"pr",
		"diff",
		fmt.Sprintf("%v", number),
		"--repo", repo,
		"--color", "never",
	}

	res, err := execute(args)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type Review struct {
	Body     string                  `json:"body,omitempty"`
	Event    string                  `json:"event,omitempty"`
	Comments []*models.GithubComment `json:"threads,omitempty"`
}

func AddReview(repo string, prNumber int, comments []*models.Comment) error {
	args := []string{
		"api",
		fmt.Sprintf("/repos/%v/pulls/%v/reviews", repo, prNumber),
		// "--method", "POST",
		// "-H", "Accept: application/vnd.github+json",
		// "-H", "X-GitHub-Api-Version: 2022-11-28",
		"--input", createInputFile(comments),
	}

	res, err := execute(args)
	if err != nil {
		pp.Println(string(res))
		return err
	}

	pp.Println("response", string(res))
	prrrr := &pullrequest{}
	err = json.Unmarshal(res, &prrrr)
	if err != nil {
		return err
	}

	// now we want to submit this review
	_, err = execute([]string{
		"api",
		fmt.Sprintf("/repos/%v/pulls/%v/reviews/%v/events", repo, prNumber, prrrr.ID),
		"-f", "event=COMMENT",
	})
	if err != nil {
		return err
	}
	log.Println("pr submitted")

	return nil
}

func createInputFile(comments []*models.Comment) string {
	r := &Review{
		Body:     "This is an automated PR Review using prollama, feel free to apply these changes.",
		Event:    "COMMENT",
		Comments: ghComments(comments),
	}

	dat, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
		return ""
	}

	pp.Println("json input file for review creation content:", string(dat))

	f, err := os.CreateTemp("", "comments")
	if err != nil {
		return ""
	}

	_, err = f.Write(dat)
	if err != nil {
		return ""
	}

	return f.Name()
}

func ghComments(c []*models.Comment) []*models.GithubComment {
	res := make([]*models.GithubComment, len(c))
	for i, v := range c {
		res[i] = v.GetGitHubCompatible()
	}

	return res
}

func GetRepositoryName() (string, error) {
	resp, err := execute([]string{
		"repo",
		"view",
		"--json", "url",
		"-q", ".url",
	})

	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(string(resp), "https://github.com/") {
		return "", fmt.Errorf("not a github repository")
	}

	// cut the newline at the end
	r := string(resp)[:len(resp)-1]

	return strings.TrimPrefix(r, "https://github.com/"), nil
}
