package giteeclient

import (
	"encoding/json"
	"fmt"

	sdk "gitee.com/openeuler/go-gitee/gitee"
)

func ConvertToNoteEvent(payload []byte) (e sdk.NoteEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}

	err = checkNoteEvent(&e)
	return
}

func checkNoteEvent(e *sdk.NoteEvent) error {
	eventType := EventTypeNote

	ne := NewNoteEventWrapper(e)

	if ne.Comment == nil {
		return fmtCheckError(eventType, "Comment")
	}

	if ne.IsPullRequest() {
		err := checkPullRequestHook(ne.PullRequest, eventType, "PullRequest")
		if err != nil {
			return err
		}
	}

	if ne.IsIssue() && ne.Issue == nil {
		return fmtCheckError(eventType, "Issue")
	}

	return checkRepository(e.Repository, eventType)
}

func ConvertToIssueEvent(payload []byte) (e sdk.IssueEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}

	err = checkIssueEvent(&e)
	return
}

func checkIssueEvent(e *sdk.IssueEvent) error {
	eventType := EventTypeIssue

	if e.Issue == nil {
		return fmtCheckError(eventType, "Issue")
	}

	return checkRepository(e.Repository, eventType)
}

func ConvertToPREvent(payload []byte) (e sdk.PullRequestEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}

	err = checkPullRequestEvent(&e)
	return
}

func checkPullRequestEvent(e *sdk.PullRequestEvent) error {
	eventType := EventTypePR

	if err := checkPullRequestHook(e.PullRequest, eventType, "PullRequest"); err != nil {
		return err
	}

	return checkRepository(e.Repository, eventType)
}

func checkPullRequestHook(pr *sdk.PullRequestHook, eventType, field string) error {
	if pr == nil {
		return fmtCheckError(eventType, field)
	}

	if pr.Head == nil || pr.Base == nil {
		return fmtCheckError(eventType, field+".Head or "+field+".Base")
	}
	return nil
}

func checkRepository(rep *sdk.ProjectHook, eventType string) error {
	field := "Repository"

	if rep == nil {
		return fmtCheckError(eventType, field)
	}

	org, repo := getOrgRepo(rep)
	if org == "" || repo == "" {
		return fmtCheckError(eventType, field+" .org or .repo")
	}
	return nil
}

func ConvertToPushEvent(payload []byte) (e sdk.PushEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}

	err = checkRepository(e.Repository, EventTypePush)
	return
}

func fmtCheckError(eventType, field string) error {
	return fmt.Errorf("%s is illegal: the field of '%s' is empty", eventType, field)
}
