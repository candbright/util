package xssh

type Config struct {
	isLocal  bool
	host     string
	user     string
	password string
	port     uint16
}

func NewConfig(isLocal bool, host string, user string, password string, port uint16) Config {
	return Config{isLocal: isLocal, host: host, user: user, password: password, port: port}
}

func SimpleConfig(host string) Config {
	return Config{host: host}
}

func (c *Config) IsLocal() bool {
	if c.host == "" {
		c.host = "127.0.0.1"
	}
	if c.host == "127.0.0.1" {
		c.isLocal = true
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
