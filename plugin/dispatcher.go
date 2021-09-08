package plugin

import (
	"encoding/json"
	"sync"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

const (
	logFieldPR   = "pr"
	logFieldOrg  = "org"
	logFieldRepo = "repo"
)

type dispatcher struct {
	c *ConfigAgent
	h *handlers

	// Tracks running handlers for graceful shutdown
	wg sync.WaitGroup
}

func (d *dispatcher) Wait() {
	d.wg.Wait() // Handle remaining requests
}

func (d *dispatcher) Dispatch(eventType string, payload []byte, l *logrus.Entry) error {
	switch eventType {
	case "Note Hook":
		if d.h.noteEventHandler == nil {
			return nil
		}

		var e gitee.NoteEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		if err := checkNoteEvent(&e); err != nil {
			return err
		}
		d.wg.Add(1)
		go d.handleNoteEvent(&e, l)

	case "Issue Hook":
		if d.h.issueHandlers == nil {
			return nil
		}

		var ie gitee.IssueEvent
		if err := json.Unmarshal(payload, &ie); err != nil {
			return err
		}
		if err := checkIssueEvent(&ie); err != nil {
			return err
		}
		d.wg.Add(1)
		go d.handleIssueEvent(&ie, l)

	case "Merge Request Hook":
		if d.h.pullRequestHandler == nil {
			return nil
		}

		var pr gitee.PullRequestEvent
		if err := json.Unmarshal(payload, &pr); err != nil {
			return err
		}
		if err := checkPullRequestEvent(&pr); err != nil {
			return err
		}
		d.wg.Add(1)
		go d.handlePullRequestEvent(&pr, l)

	case "Push Hook":
		if d.h.pushEventHandler == nil {
			return nil
		}

		var pe gitee.PushEvent
		if err := json.Unmarshal(payload, &pe); err != nil {
			return err
		}
		if err := checkRepository(pe.Repository, "push event"); err != nil {
			return err
		}
		d.wg.Add(1)
		go d.handlePushEvent(&pe, l)

	default:
		l.Debug("Ignoring unhandled event type")
	}
	return nil
}

func (d *dispatcher) handlePullRequestEvent(pr *gitee.PullRequestEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  pr.Repository.Namespace,
		logFieldRepo: pr.Repository.Path,
		logFieldPR:   pr.PullRequest.Number,
		"author":     pr.PullRequest.User.Login,
		"url":        pr.PullRequest.HtmlUrl,
	})
	l.Infof("Pull request %s.", *pr.Action)

	if err := d.h.pullRequestHandler(pr, d.c.GetConfig(), l); err != nil {
		l.WithError(err).Error("Error handling PullRequestEvent.")
	}
}

func (d *dispatcher) handleIssueEvent(i *gitee.IssueEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  i.Repository.Namespace,
		logFieldRepo: i.Repository.Path,
		logFieldPR:   i.Issue.Number,
		"author":     i.Issue.User.Login,
		"url":        i.Issue.HtmlUrl,
	})
	l.Infof("Issue %s.", *i.Action)

	if err := d.h.issueHandlers(i, d.c.GetConfig(), l); err != nil {
		l.WithError(err).Error("Error handling IssueEvent.")
	}
}

func (d *dispatcher) handlePushEvent(pe *gitee.PushEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  pe.Repository.Namespace,
		logFieldRepo: pe.Repository.Path,
		"ref":        pe.Ref,
		"head":       pe.After,
	})
	l.Info("Push event.")

	if err := d.h.pushEventHandler(pe, d.c.GetConfig(), l); err != nil {
		l.WithError(err).Error("Error handling PushEvent.")
	}
}

func (d *dispatcher) handleNoteEvent(e *gitee.NoteEvent, l *logrus.Entry) {
	defer d.wg.Done()

	var n interface{}
	switch *(e.NoteableType) {
	case "PullRequest":
		n = e.PullRequest.Number
	case "Issue":
		n = e.Issue.Number
	}

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  e.Repository.Namespace,
		logFieldRepo: e.Repository.Path,
		logFieldPR:   n,
		"review":     e.Comment.Id,
		"commenter":  e.Comment.User.Login,
		"url":        e.Comment.HtmlUrl,
	})
	l.Infof("Note %s.", *e.Action)

	if err := d.h.noteEventHandler(e, d.c.GetConfig(), l); err != nil {
		l.WithError(err).Error("Error handling NoteEvent.")
	}
}
