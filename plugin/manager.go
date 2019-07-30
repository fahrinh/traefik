package plugin

import (
	"fmt"
	"github.com/containous/traefik/log"
	"github.com/containous/traefik/safe"
	goplugin "plugin"
)



type Manager struct {
	// 1) matchers
	matchers map[string]*IMatcher

	// 2) middlewares
}

// NewManager builds a new manager
func NewManager() *Manager {
	return &Manager{
		matchers: map[string]*IMatcher{},
	}
}

// Load loads a plugin
func (m *Manager) Load(pluginId string, pluginPath string) error {
	errChan := make(chan error)
	defer close(errChan)

	safe.GoWithRecover(func() {
		logSuffix := fmt.Sprintf("[%s][%s]", pluginId, pluginPath)

		p, err := goplugin.Open(pluginPath)
		if err != nil {
			errChan <- fmt.Errorf("%s error opening plugin: %s", logSuffix, err)
			return
		}

		loader, err := p.Lookup("Load")
		if err != nil {
			errChan <- fmt.Errorf("%s error finding Load() interface{} function: %s", logSuffix, err)
			return
		}
		load, ok := loader.(func() interface{})
		if !ok {
			errChan <- fmt.Errorf("%s plugin does not implement Load() interface{} function", logSuffix)
			return
		}
		instance := load()
		if instance == nil {
			errChan <- fmt.Errorf("%s plugin does Load() nil instance", logSuffix)
			return
		}

		matcherInstance, ok := instance.(IMatcher)
		if !ok {
			errChan <- fmt.Errorf("%s plugin does not implement any plugin interface", logSuffix)
			return
		}

		m.matchers[pluginId] = &matcherInstance

		errChan <- nil
		return
	}, func(err interface{}) {
		log.Errorf("Error in plugin Go routine: %s", err)
	})

	if err, ok := <-errChan; ok {
		return err
	}
	return nil
}

// return a list of all IMatcher
func (m *Manager) GetMatchers() map[string]*IMatcher {
	return m.matchers
}
