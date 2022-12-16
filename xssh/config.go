package xssh

type Config struct {
	isLocal    bool
	host       string
	user       string
	password   string
	port       uint16
	localHosts []string
}

func NewConfig(isLocal bool, host string, user string, password string, port uint16) Config {
	return Config{isLocal: isLocal, host: host, user: user, password: password, port: port}
}

func SimpleConfig(host string) Config {
	return Config{host: host}
}

func (c *Config) AddLocalHost(host string) {
	if c.localHosts == nil {
		c.localHosts = make([]string, 0)
		c.localHosts = append(c.localHosts, host)
		return
	}
	for _, localHost := range c.localHosts {
		if localHost == host {
			return
		}
	}
	c.localHosts = append(c.localHosts, host)
}

func (c *Config) IsLocal() bool {
	if c.host == "" {
		c.host = "127.0.0.1"
	}
	if c.host == "127.0.0.1" {
		c.isLocal = true
		return c.isLocal
	}
	if c.localHosts != nil {
		for _, localHost := range c.localHosts {
			if localHost == c.host {
				c.isLocal = true
				return c.isLocal
			}
		}
	}
	return c.isLocal
}

func (c *Config) Host() string {
	if c.host == "" {
		c.host = "127.0.0.1"
	}
	return c.host
}

func (c *Config) User() string {
	if c.user == "" {
		c.user = "root"
	}
	return c.user
}

func (c *Config) Port() uint16 {
	if c.port == 0 {
		c.port = 22
	}
	return c.port
}

func (c *Config) Password() string {
	return c.password
}
