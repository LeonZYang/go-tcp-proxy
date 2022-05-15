package proxy

import (
	"fmt"
)

var conf = new(Config)

type Config struct {
	Debug bool                   `toml:"debug"`
	Proxy map[string]ProxyConfig `toml:"proxy"`
}

type ProxyConfig struct {
	Enabled bool         `toml:"enabled"`
	Listen  ListenConfig `toml:"listen"`
	Remote  ListenConfig `toml:"remote"`
}

type ListenConfig struct {
	Addr     string `toml:"addr"`
	TLS      bool   `toml:"tls"`
	Ca       string `toml:"ca"`
	PrivFile string `toml:"privFile"`
	PubFile  string `toml:"pubFile"`
}

func (conf *Config) SetDefault() {
	*conf = Config{}
}

func (conf *Config) Print() {
}

func (conf *Config) GetRPCAddr() string {
	return ""
}

func (conf Config) IsDebug() bool {
	return conf.Debug
}

func (conf *Config) Init() {
}

func (conf Config) Version() {
	fmt.Println("0.0.1")
}

func GetConfig() *Config {
	return conf
}
