package model

type App struct {
	Active             string   `yaml:"active"`
	Name               string   `yaml:"name"`
	Port               int      `yaml:"port"`
	Version            string   `yaml:"version"`
	Authors            []string `yaml:"authors"`
	Repository         string   `yaml:"repository"`
	AdminAuthVerifyUrl string   `yaml:"admin-auth-verify-url"`
}

func (a *App) Release() bool {
	return a.Active == "release"
}
func (a *App) Us() bool {
	return a.Active == "us"
}

func (a *App) Develop() bool {
	return a.Active == "develop"
}

func (a *App) Test() bool {
	return a.Active == "test"
}
