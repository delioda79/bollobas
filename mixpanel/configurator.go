package mixpanel

type Configurator struct {
	s map[string]interface{}
}

func (c *Configurator) Configure(s map[string]interface{}) {
	c.s = s
}

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
