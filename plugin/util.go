package plugin

import (
	"fmt"

	sdk "gitee.com/openeuler/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-plugin/giteeclient"
)

func checkNoteEvent(e *sdk.NoteEvent) error {
	eventType := "note event"
	ne := giteeclient.NewNoteEventWrapper(e)
	if ne.Comment == nil {
		return fmtCheckError(eventType, "comment")
	}
	if ne.IsPullRequest() {
		if err := checkPullRequestHook(ne.PullRequest, eventType); err != nil {
			return err
		}
	}
	if ne.IsIssue() && ne.Issue == nil {
		return fmtCheckError(eventType, "issue")
	}
	return checkRepository(e.Repository, eventType)
}

func checkIssueEvent(e *sdk.IssueEvent) error {
	eventType := "issue event"
	if e.Issue == nil {
		return fmtCheckError(eventType, "issue")
	}
	return checkRepository(e.Repository, eventType)
}

func checkPullRequestEvent(e *sdk.PullRequestEvent) error {
	eventType := "pull request event"
	if err := checkPullRequestHook(e.PullRequest, eventType); err != nil {
		return err
	}
	return checkRepository(e.Repository, eventType)
}

func checkPullRequestHook(pr *sdk.PullRequestHook, eventType string) error {
	if pr == nil {
		return fmtCheckError(eventType, "pull_request")
	}
	if pr.Head == nil || pr.Base == nil {
		return fmtCheckError(eventType, "pull_request.head or pull_request.base")
	}
	return nil
}

func checkRepository(rep *sdk.ProjectHook, eventType string) error {
	if rep == nil {
		return fmtCheckError(eventType, "pull_request")
	}
	if rep.Namespace == "" || rep.Path == "" {
		return fmtCheckError(eventType, "pull_request.namespace or pull_request.path")
	}
	return nil
}

func fmtCheckError(eventType, field string) error {
	return fmt.Errorf("%s is illegal: the %s field is empty", eventType, field)
}
