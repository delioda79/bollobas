package mixpanel

//Configurator is a settings configuratin manager for mixpanel related handlers
type Configurator struct {
	s map[string]interface{}
}

//Configure updates the configuration
func (c *Configurator) Configure(s map[string]interface{}) {
	c.s = s
}

//Check returns a boolean stating if the message can be processed or not, based on the settings and the message itself
func (c *Configurator) Check(msg []byte) bool {
	mp, ok := c.s["mixpanel_passenger_enabled"]
	if !ok {
		return false
	}

	r, ok := mp.(bool)
	if !ok {
		return false
	}

	return r
}
