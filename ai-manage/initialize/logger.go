package initialize

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
	"vland.live/app/global"
)

func initLogger(ctx context.Context) {
	logger := global.Config.Zap
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logger.Filename,
		MaxSize:    logger.MaxSize,
		MaxAge:     logger.MaxAge,
		MaxBackups: logger.MaxBackup,
	})
	l := new(zapcore.Level)
	if err := l.UnmarshalText([]byte(logger.Level)); err != nil {
		panic(fmt.Sprintf("init logger error: %v", err.Error()))
	}
	core := zapcore.NewTee(
		zapcore.NewCore(getEncoder(), writeSyncer, l),
		zapcore.NewCore(getEncoder(), os.Stdout, l))
	global.Logger = zap.New(core, zap.AddCaller())
	log.Println("Zap logger Client init success")
	zap.ReplaceGlobals(global.Logger)

	go func() {
		<-ctx.Done()
		log.Println("global.Config.Zap: 接收到了停机请求")
		time.Sleep(time.Millisecond * 10)
		log.Println("global.Config.Zap: close")
	}()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
