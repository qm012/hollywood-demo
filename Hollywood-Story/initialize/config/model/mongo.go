package model

type Mongo struct {
	Path        string `yaml:"path"`
	Database    string `yaml:"database"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	MaxPoolSize uint64 `yaml:"max-pool-size"`
}
