package plugin

import (
	"flag"
	"time"
)

type Options struct {
	port         int
	gracePeriod  time.Duration
	pluginConfig string
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.port, "port", 8888, "Port to listen on.")
	fs.StringVar(&o.pluginConfig, "plugin-config", "/etc/plugins/plugins.yaml", "Path to plugin config file.")
	fs.DurationVar(&o.gracePeriod, "grace-period", 180*time.Second, "On shutdown, try to handle remaining events for the specified duration.")
}
