package plugin

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

type PluginConfig interface {
	Validate() error
	SetDefault()
}

type NewPluginConfig func() PluginConfig

type ConfigAgent struct {
	mut    sync.RWMutex
	c      PluginConfig
	b      NewPluginConfig
	stop   chan bool
	md5Sum string
}

func newConfigAgent(b NewPluginConfig) *ConfigAgent {
	return &ConfigAgent{b: b}
}

func (ca *ConfigAgent) load(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	v := fmt.Sprintf("%x", md5.Sum(b))
	if ca.md5Sum == v {
		return nil
	}

	c := ca.b()
	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}

	c.SetDefault()

	if err := c.Validate(); err != nil {
		return err
	}

	ca.md5Sum = v

	ca.mut.Lock()
	ca.c = c
	ca.mut.Unlock()

	return nil
}

func (ca *ConfigAgent) GetConfig() PluginConfig {
	ca.mut.Lock()
	v := ca.c
	ca.mut.Unlock()

	return v
}

// Start starts polling path for plugin config. If the first attempt fails,
// then start returns the error. Future errors will halt updates but not stop.
func (ca *ConfigAgent) Start(path string) error {
	if err := ca.load(path); err != nil {
		return err
	}

	ticker := time.Tick(1 * time.Minute)
	go func() {
		for {
			//TODO is it right
			select {
			case <-ticker:
				if err := ca.load(path); err != nil {
					logrus.WithField("path", path).WithError(err).Error("Error loading plugin config.")
				}
			case <-ca.stop:
				break
			}
		}
	}()

	return nil
}

func (ca *ConfigAgent) Stop() {
	ca.stop <- true
}
