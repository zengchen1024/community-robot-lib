package options

import (
	"flag"
	"time"
)

type PluginOptions struct {
	Port         int
	GracePeriod  time.Duration
	PluginConfig string
}

func (o *PluginOptions) Validate() error {
	return nil
}

func (o *PluginOptions) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.Port, "port", 8888, "Port to listen on.")
	fs.StringVar(&o.PluginConfig, "plugin-config", "/etc/plugins/plugins.yaml", "Path to plugin config file.")
	fs.DurationVar(&o.GracePeriod, "grace-period", 180*time.Second, "On shutdown, try to handle remaining events for the specified duration.")
}
