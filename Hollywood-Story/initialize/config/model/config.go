package model

type Config struct {
	App   *App   `yaml:"app"`
	Proxy *Proxy `yaml:"proxy"`
	AI    *AI    `yaml:"ai"`
	Mongo *Mongo `yaml:"mongo"`
	Zap   *Zap   `yaml:"zap"`
}
