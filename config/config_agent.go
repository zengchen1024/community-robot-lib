package config

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

// The returned value of NewPluginConfig must be a pointer.
// Otherwise, it should deep copy the config when reading it.
type NewPluginConfig func() PluginConfig

type ConfigAgent struct {
	mut    sync.RWMutex
	c      PluginConfig
	b      NewPluginConfig
	stop   chan bool
	md5Sum string
}

func NewConfigAgent(b NewPluginConfig) ConfigAgent {
	return ConfigAgent{b: b}
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

	ca.mut.Lock()
	ca.md5Sum = v
	ca.c = c
	ca.mut.Unlock()

	return nil
}

func (ca *ConfigAgent) GetConfig() (string, PluginConfig) {
	ca.mut.RLock()
	c := ca.c // copy the pointer
	v := ca.md5Sum
	ca.mut.RUnlock()

	return v, c
}

// Start starts polling path for plugin config. If the first attempt fails,
// then start returns the error. Future errors will halt updates but not stop.
func (ca *ConfigAgent) Start(path string) error {
	if err := ca.load(path); err != nil {
		return err
	}

	l := logrus.WithField("path", path)

	ticker := time.Tick(1 * time.Minute)
	go func() {
		for {
			//TODO is it right
			select {
			case <-ticker:
				if err := ca.load(path); err != nil {
					l.WithError(err).Error("Error loading plugin config.")
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
