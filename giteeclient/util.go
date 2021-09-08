package giteeclient

import (
	"strings"

	sdk "gitee.com/openeuler/go-gitee/gitee"
	"k8s.io/apimachinery/pkg/util/sets"
)

type PRInfo struct {
	Org     string
	Repo    string
	BaseRef string
	HeadSHA string
	Author  string
	Number  int
	Labels  sets.String
}

func (pr PRInfo) HasLabel(l string) bool {
	return pr.Labels.Has(l)
}

func GetPRInfoByPREvent(pre *sdk.PullRequestEvent) PRInfo {
	pr := pre.PullRequest

	return PRInfo{
		Org:     pre.Repository.Namespace,
		Repo:    pre.Repository.Path,
		BaseRef: pr.Base.Ref,
		HeadSHA: pr.Head.Sha,
		Author:  pr.User.Login,
		Number:  int(pr.Number),
		Labels:  getLabelFromEvent(pr.Labels),
	}
}

func getLabelFromEvent(labels []sdk.LabelHook) sets.String {
	m := sets.NewString()
	for i := range labels {
		m.Insert(labels[i].Name)
	}
	return m
}

//GetOwnerAndRepoByPREvent obtain the owner and repository name from the pullrequest's event
func GetOwnerAndRepoByPREvent(pre *sdk.PullRequestEvent) (string, string) {
	return pre.Repository.Namespace, pre.Repository.Path
}

//GetOwnerAndRepoByIssueEvent obtain the owner and repository name from the issue's event
func GetOwnerAndRepoByIssueEvent(issue *sdk.IssueEvent) (string, string) {
	return issue.Repository.Namespace, issue.Repository.Path
}

const (
	PRActionOpened              = "opened"
	PRActionClosed              = "closed"
	PRActionUpdatedLabel        = "update_label"
	PRActionChangedTargetBranch = "target_branch_changed"
	PRActionChangedSourceBranch = "source_branch_changed"
)

func GetPullRequestAction(e *sdk.PullRequestEvent) string {
	switch strings.ToLower(*(e.Action)) {
	case "open":
		return PRActionOpened

	case "update":
		switch strings.ToLower(*(e.ActionDesc)) {
		case "source_branch_changed": // change the pr's commits
			return PRActionChangedSourceBranch

		case "target_branch_changed": // change the branch to which this pr will be merged
			return PRActionChangedTargetBranch

		case "update_label":
			return PRActionUpdatedLabel
		}

	case "close":
		return PRActionClosed
	}

	return ""
}
