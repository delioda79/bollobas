package mixpanel

import "sync"

//Configurator is a settings configuratin manager for mixpanel related handlers
type Configurator struct {
	s  map[string]interface{}
	mx *sync.Mutex
}

//Configure updates the configuration
func (c *Configurator) Configure(s map[string]interface{}) {
	c.mx.Lock()
	c.s = s
	c.mx.Unlock()
}

//Check returns a boolean stating if the message can be processed or not, based on the settings and the message itself
func (c *Configurator) Check(msg []byte) bool {
	c.mx.Lock()
	mp, ok := c.s["mixpanel_passenger_enabled"]
	c.mx.Unlock()
	if !ok {
		return false
	}

	r, ok := mp.(bool)
	if !ok {
		return false
	}

	return r
}

// NewConfigurator returns a new configurator
func NewConfigurator() *Configurator {
	mx := &sync.Mutex{}
	return &Configurator{mx: mx, s: map[string]interface{}{}}
}
