package global

import (
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"vland.live/app/initialize/config/model"
)

var (
	AIAzureClient  *openai.Client // openai的客户端
	AIOpenaiClient *openai.Client // 微软的openai客户端
	AIClient       *openai.Client // 通用的AI客户端
	Config         *model.Config  // 环境配置变量
	HttpClient     *http.Client   // 通用http请求客户端
	Mongo          *mongo.Client  // mongo存储
	Logger         *zap.Logger    // 日志
)
