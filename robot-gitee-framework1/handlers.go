package framework

import (
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

// IssueHandler defines the function contract for a gitee.IssueEvent handler.
type IssueHandler func(e *sdk.IssueEvent, log *logrus.Entry) error

// PullRequestHandler defines the function contract for a sdk.PullRequestEvent handler.
type PullRequestHandler func(e *sdk.PullRequestEvent, log *logrus.Entry) error

// PushEventHandler defines the function contract for a sdk.PushEvent handler.
type PushEventHandler func(e *sdk.PushEvent, log *logrus.Entry) error

// NoteEventHandler defines the function contract for a sdk.NoteEvent handler.
type NoteEventHandler func(e *sdk.NoteEvent, log *logrus.Entry) error

type handlers struct {
	pullRequestHandler PullRequestHandler
	pushEventHandler   PushEventHandler
	noteEventHandler   NoteEventHandler
	issueHandler       IssueHandler
}

// RegisterIssueHandler registers a plugin's gitee.IssueEvent handler.
func (h *handlers) RegisterIssueHandler(fn IssueHandler) {
	h.issueHandler = fn
}

// RegisterPullRequestHandler registers a plugin's gitee.PullRequestEvent handler.
func (h *handlers) RegisterPullRequestHandler(fn PullRequestHandler) {
	h.pullRequestHandler = fn
}

// RegisterPushEventHandler registers a plugin's gitee.PushEvent handler.
func (h *handlers) RegisterPushEventHandler(fn PushEventHandler) {
	h.pushEventHandler = fn
}

// RegisterNoteEventHandler registers a plugin's gitee.NoteEvent handler.
func (h *handlers) RegisterNoteEventHandler(fn NoteEventHandler) {
	h.noteEventHandler = fn
}

func (h *handlers) getHandler() (r map[string]func(payload []byte, l *logrus.Entry)) {
	r = make(map[string]func(payload []byte, l *logrus.Entry))

	if h.noteEventHandler != nil {
		r[sdk.EventTypeNote] = h.handleNoteEvent
	}

	if h.issueHandler != nil {
		r[sdk.EventTypeIssue] = h.handleIssueEvent
	}

	if h.pullRequestHandler != nil {
		r[sdk.EventTypePR] = h.handlePullRequestEvent
	}

	if h.pushEventHandler == nil {
		r[sdk.EventTypePush] = h.handlePushEvent
	}

	return
}

func (h *handlers) handlePullRequestEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToPREvent(payload)
	if err != nil {
		l.Errorf("convert to PREvent, err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.PullRequest.HtmlUrl,
		logFieldAction: e.GetActionDesc(),
	})

	if err := h.pullRequestHandler(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}

func (h *handlers) handleIssueEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToIssueEvent(payload)
	if err != nil {
		l.Errorf("convert to IssueEvent, err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.Issue.HtmlUrl,
		logFieldAction: *e.Action,
	})

	if err := h.issueHandler(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}

func (h *handlers) handlePushEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToPushEvent(payload)
	if err != nil {
		l.Errorf("convert to PushEvent, err: ", err.Error())

		return
	}

	org, repo := e.GetOrgRepo()

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"ref":        e.Ref,
		"head":       e.After,
	})

	if err := h.pushEventHandler(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}

func (h *handlers) handleNoteEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToNoteEvent(payload)
	if err != nil {
		l.Errorf("convert to NoteEvent, err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		"commenter":    e.Comment.User.Login,
		logFieldURL:    e.Comment.HtmlUrl,
		logFieldAction: *e.Action,
	})

	if err := h.noteEventHandler(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}
