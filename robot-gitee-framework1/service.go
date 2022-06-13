package framework

import (
	"net/http"
	"strconv"
	"time"

	"github.com/opensourceways/community-robot-lib/interrupts"
)

type Service interface {
	RegisterIssueHandler(IssueHandler)
	RegisterPullRequestHandler(PullRequestHandler)
	RegisterPushEventHandler(PushEventHandler)
	RegisterNoteEventHandler(NoteEventHandler)

	Run(int, time.Duration)
}

func NewService() Service {
	return &service{}
}

type service struct {
	handlers
}

func (s *service) Run(port int, timeout time.Duration) {
	d := dispatcher{h: s.handlers}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		d.wait()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/gitee-hook", &d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(port)}

	interrupts.ListenAndServe(httpServer, timeout)
}
