package main

import (
	"errors"

	"github.com/opensourceways/community-robot-lib/config"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

// TODO: set botName
const botName = ""

type iClient interface {
}

func newRobot(cli iClient, gc func() *configuration) *robot {
	return &robot{cli: cli, gc: gc}
}

type robot struct {
	agent *config.ConfigAgent
	gc    func() *configuration
	cli   iClient
}

func (bot *robot) getConfig() (*configuration, error) {
	if c := bot.gc(); c != nil {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand PR event, delete this function.
	return nil
}

func (bot *robot) handleIssueEvent(e *sdk.IssueEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Issue event, delete this function.
	return nil
}

func (bot *robot) handlePushEvent(e *sdk.PushEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Push event, delete this function.
	return nil
}

func (bot *robot) handleNoteEvent(e *sdk.NoteEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Note event, delete this function.
	return nil
}
