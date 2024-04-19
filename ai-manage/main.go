package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
	"vland.live/app/global"
	"vland.live/app/initialize"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.TODO())
	initialize.Run(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("接收到中断信号，开始优雅停机...\n\t\t    注意：正在关闭相关资源，请勿强制退出...")
	cancelFn() // 发出停止信号

	i := 3
	for i > 0 {
		log.Printf("剩余等待时间：%d秒...\n", i)
		i--
		time.Sleep(time.Second)
	}
	log.Printf("优雅停机已完成。感谢使用%s的服务！\n", global.Config.App.Name)
}
