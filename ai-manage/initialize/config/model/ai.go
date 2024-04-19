package model

type AI struct {
	Enable string  `yaml:"enable"`
	Openai *OpenAI `yaml:"openai"`
	Azure  *Azure  `yaml:"azure"`
}

type OpenAI struct {
	Key string `yaml:"key"`
}

type Azure struct {
	Key     string `yaml:"key"`
	BaseUrl string `yaml:"base-url"`
	OrgId   string `yaml:"org-id"`
}
