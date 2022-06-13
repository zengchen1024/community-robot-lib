package framework

import (
	"net/http"
	"strconv"

	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/opensourceways/community-robot-lib/options"
)

type Service interface {
	RegisterIssueHandler(IssueHandler)
	RegisterPullRequestHandler(PullRequestHandler)
	RegisterPushEventHandler(PushEventHandler)
	RegisterNoteEventHandler(NoteEventHandler)

	Run(options.ServiceOptions)
}

func NewService() Service {
	return &service{}
}

type service struct {
	handlers
}

func (s *service) Run(o options.ServiceOptions) {
	d := dispatcher{h: s.handlers}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		d.wait()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/gitee-hook", &d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(o.Port)}

	interrupts.ListenAndServe(httpServer, o.GracePeriod)
}
