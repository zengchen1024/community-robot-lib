package plugin

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-plugin/plugin/interrupts"
)

type HandlerRegitster interface {
	RegisterIssueHandler(IssueHandler)
	RegisterPullRequestHandler(PullRequestHandler)
	RegisterPushEventHandler(PushEventHandler)
	RegisterNoteEventHandler(NoteEventHandler)
}

type Plugin interface {
	NewPluginConfig() PluginConfig
	RegisterEventHandler(HandlerRegitster)
	Exit()
}

func Run(p Plugin, o Options) {
	agent := newConfigAgent(p.NewPluginConfig)
	if err := agent.Start(o.pluginConfig); err != nil {
		return
	}

	h := handlers{}
	p.RegisterEventHandler(&h)

	d := &dispatcher{c: agent, h: &h}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		d.Wait()
		p.Exit()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/gitee-hook", d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(o.port)}

	interrupts.ListenAndServe(httpServer, o.gracePeriod)
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventType, eventGUID, payload := parseRequest(w, r)

	l := logrus.WithFields(
		logrus.Fields{
			"event-type": eventType,
			"event_id":   eventGUID,
		},
	)

	if err := d.Dispatch(eventType, payload, l); err != nil {
		l.WithError(err).Error("Error parsing event.")
	}
}

func parseRequest(w http.ResponseWriter, r *http.Request) (eventType string, uuid string, payload []byte) {
	defer r.Body.Close()

	resp := func(code int, msg string) {
		http.Error(w, msg, code)
	}

	if eventType = r.Header.Get("X-Gitee-Event"); eventType == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-Gitee-Event Header")
		return
	}

	if uuid = r.Header.Get("X-Gitee-Timestamp"); uuid == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-Gitee-Timestamp Header")
		return
	}

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp(http.StatusInternalServerError, "500 Internal Server Error: Failed to read request body")
		return "", "", nil
	}
	payload = v

	return
}
