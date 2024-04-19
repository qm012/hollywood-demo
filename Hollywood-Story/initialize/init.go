package initialize

import (
	"context"
	"fmt"
	"log"
	"vland.live/app/global"
	"vland.live/app/initialize/config"
	"vland.live/app/initialize/router"
	"vland.live/app/internal/constant"
)

func Run(ctx context.Context) {
	{
		config.Init()               // 初始化配置信息
		initLogger(ctx)             // 初始化日志客户端
		initOpenai()                // 初始化openai的客户端
		initMongo(ctx)              // 初始化mongo客户端
		constant.InitCaptureEvent() // 初始化事件
	}

	r := router.Init() // 初始化路由
	// 启动http服务
	var port = fmt.Sprintf(":%d", global.Config.App.Port)
	if err := r.Run(port); err != nil {
		log.Fatalln(err)
	}
}
