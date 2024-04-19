package model

type Zap struct {
	Level     string `yaml:"level"`
	Filename  string `yaml:"filename"`
	MaxSize   int    `yaml:"max-size"`
	MaxBackup int    `yaml:"max-backup"`
	MaxAge    int    `yaml:"max-age"`
}
