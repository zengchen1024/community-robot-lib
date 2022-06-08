package githubclient

import (
	"fmt"

	sdk "github.com/google/go-github/v36/github"
)

type PRInfo struct {
	Org    string
	Repo   string
	Number int
}

func (p PRInfo) String() string {
	return fmt.Sprintf("%s/%s:%d", p.Org, p.Repo, p.Number)
}

// Client interface for Github API
type Client interface {
	AddPRLabel(pr PRInfo, label string) error
	RemovePRLabel(pr PRInfo, label string) error
	CreatePRComment(pr PRInfo, comment string) error
	DeletePRComment(org, repo string, ID int64) error
	GetPRCommits(pr PRInfo) ([]*sdk.RepositoryCommit, error)
	GetPRComments(pr PRInfo) ([]*sdk.IssueComment, error)
}
