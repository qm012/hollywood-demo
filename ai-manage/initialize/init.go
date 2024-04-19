package initialize

import (
	"context"
	"fmt"
	"log"
	"vland.live/app/global"
	"vland.live/app/initialize/config"
	"vland.live/app/initialize/router"
)

func Run(ctx context.Context) {
	{
		config.Init()   // 初始化配置信息
		initHttp(ctx)   // 初始化全局的http请求客户端
		initLogger(ctx) // 初始化日志客户端
		initOpenai()    // 初始化openai的客户端
		initMongo(ctx)  // 初始化mongo客户端
	}
	r := router.Init() // 初始化路由
	go func() {
		// 启动http服务
		var port = fmt.Sprintf(":%d", global.Config.App.Port)
		if err := r.Run(port); err != nil {
			log.Fatalln(err)
		}
	}()
}
