package githubclient

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	sdk "github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

func NewClient(getToken func() []byte) Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: string(getToken()),
	})
	tc := oauth2.NewClient(context.Background(), ts)

	return client{sdk.NewClient(tc)}
}

type client struct {
	c *sdk.Client
}

func (cl client) AddPRLabel(pr PRInfo, label string) error {
	_, _, err := cl.c.Issues.AddLabelsToIssue(
		context.Background(),
		pr.Org, pr.Repo, pr.Number, []string{label},
	)

	return err
}

func (cl client) RemovePRLabel(pr PRInfo, label string) error {
	r, err := cl.c.Issues.RemoveLabelForIssue(
		context.Background(),
		pr.Org, pr.Repo, pr.Number, label,
	)
	if err != nil && r != nil && r.StatusCode == 404 {
		return nil
	}

	return err
}

func (cl client) CreatePRComment(pr PRInfo, comment string) error {
	ic := sdk.IssueComment{
		Body: sdk.String(comment),
	}
	_, _, err := cl.c.Issues.CreateComment(
		context.Background(),
		pr.Org, pr.Repo, pr.Number, &ic,
	)

	return err
}

func (cl client) DeletePRComment(org, repo string, commentId int64) error {
	_, err := cl.c.Issues.DeleteComment(context.Background(), org, repo, commentId)

	return err
}

func (cl client) GetPRComments(pr PRInfo) ([]*sdk.IssueComment, error) {
	comments := []*sdk.IssueComment{}

	opt := &sdk.IssueListCommentsOptions{}
	opt.Page = 1

	for {
		v, resp, err := cl.c.Issues.ListComments(context.Background(), pr.Org, pr.Repo, pr.Number, opt)
		if err != nil {
			return comments, err
		}

		comments = append(comments, v...)

		link := parseLinks(resp.Header.Get("Link"))["next"]
		if link == "" {
			break
		}

		pagePath, err := url.Parse(link)
		if err != nil {
			break
		}

		p := pagePath.Query().Get("page")
		if p == "" {
			break
		}

		page, err := strconv.Atoi(p)
		if err != nil {
			break
		}
		opt.Page = page
	}

	return comments, nil
}

func (cl client) GetPRCommits(pr PRInfo) ([]*sdk.RepositoryCommit, error) {
	commits := []*sdk.RepositoryCommit{}

	f := func() error {
		opt := &sdk.ListOptions{}
		opt.Page = 1

		for {
			v, resp, err := cl.c.PullRequests.ListCommits(context.Background(), pr.Org, pr.Repo, pr.Number, nil)
			if err != nil {
				return err
			}

			commits = append(commits, v...)

			link := parseLinks(resp.Header.Get("Link"))["next"]
			if link == "" {
				break
			}

			pagePath, err := url.Parse(link)
			if err != nil {
				return fmt.Errorf("failed to parse 'next' link: %v", err)
			}

			p := pagePath.Query().Get("page")
			if p == "" {
				return fmt.Errorf("failed to get 'page' on link: %s", p)
			}

			page, err := strconv.Atoi(p)
			if err != nil {
				return err
			}

			opt.Page = page
		}

		return nil
	}

	err := f()

	return commits, err
}
