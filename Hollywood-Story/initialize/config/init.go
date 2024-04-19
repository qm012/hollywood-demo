package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strings"
	"vland.live/app/global"
	"vland.live/app/initialize/config/model"
)

func readFile(filePath string) (config *model.Config, err error) {
	if filePath == "" {
		return nil, errors.New("filePath is null")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("init config file err: %v", err.Error())
	}
	all, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read config file data err: %v", err.Error())
	}
	config = new(model.Config)
	err = yaml.Unmarshal(all, config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file data err: %v", err.Error())
	}

	return
}

func profileActive(configActive string) (active string) {
	// 读取系统变量
	// Env > config
	log.Println("从系统环境变量中读取激活环境")
	active = os.Getenv("VERSE_ACTIVE")
	if len(active) == 0 {
		log.Println("未从系统环境变量中读取到激活环境，继续从配置文件中读取。")
		// 从配置文件读取环境
		active = configActive
		if len(active) == 0 {
			log.Println("未从配置文件中读取到激活环境")
			panic("env [ACTIVE] cant be empty")
		}
	}
	log.Printf("已加载到激活环境【%v】\n", active)
	return
}

func Init() {
	// 读取应用主体配置文件
	log.Printf("正在加载应用的配置文件【%v】\n", "./app.yaml")
	config, err := readFile("./app.yaml")
	if err != nil {
		log.Fatalln(`readFile("./app.yaml") failed`, err)
	}
	// 读取对应环境的
	active := profileActive(config.App.Active)

	filePath := fmt.Sprintf("./app-%v.yaml", active)
	log.Printf("正在加载环境的配置文件【%v】\n", filePath)
	envServer, err := readFile(filePath)
	if err != nil {
		log.Fatalln("readFile filePath failed", filePath, err)
	}

	config.Zap.Level = envServer.Zap.Level

	// 代理
	{
		if envServer.Proxy != nil {
			config.Proxy = envServer.Proxy
		}
	}

	// mongo
	{
		if envServer.Mongo.Path != "" {
			config.Mongo.Path = envServer.Mongo.Path
		}
		if envServer.Mongo.MaxPoolSize != 0 {
			config.Mongo.MaxPoolSize = envServer.Mongo.MaxPoolSize
		}
		config.Mongo.Username = envServer.Mongo.Username
		config.Mongo.Password = envServer.Mongo.Password
	}

	// 重新赋值获取到的active
	config.App.Active = active
	global.Config = config

	//if config.App.Release() || config.App.Us() {
	//gin.SetMode(gin.ReleaseMode)
	//}
	// 打印信息
	printAppInfo(config.App)
}

func printAppInfo(app *model.App) {
	str := `

		应用基本信息
	名称:  %s
	端口:  %d
	环境:  %s
	版本:  %s
	作者:  %#v
	仓库地址:  %s

`
	log.Printf(str, app.Name, app.Port, app.Active,
		app.Version, strings.Join(app.Authors, ","), app.Repository,
	)
}
