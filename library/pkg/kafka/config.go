package kafka

import "github.com/BurntSushi/toml"

type Config struct {
	Addr []string `toml:"addr"`
}

func (c *Config) Set(s string) (err error) {
	return toml.Unmarshal([]byte(s), c)
}
