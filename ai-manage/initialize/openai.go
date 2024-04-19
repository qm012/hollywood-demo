package initialize

import (
	"github.com/sashabaranov/go-openai"
	"log"
	"vland.live/app/global"
)

func initOpenai() {

	// 微软的openai客户端初始化
	{
		azureClientConfig := openai.DefaultAzureConfig(global.Config.AI.Azure.Key, global.Config.AI.Azure.BaseUrl)
		azureClientConfig.AzureModelMapperFunc = func(mode string) string {
			azureModelMapping := map[string]string{
				openai.GPT3Dot5Turbo: "xxx", // model name
				//openai.AdaEmbeddingV2.String(): "xxx",// model name
				openai.AdaEmbeddingV2.String(): "xxx", // model name
				openai.GPT3Dot5Turbo1106:       "xxx", // model name
			}
			return azureModelMapping[mode]
		}

		azureClientConfig.HTTPClient = global.HttpClient
		global.AIAzureClient = openai.NewClientWithConfig(azureClientConfig)
		log.Println("AIAzureClient init success")
	}

	// openai 的客户端初始化
	{
		openaiClientConfig := openai.DefaultConfig(global.Config.AI.Openai.Key)
		openaiClientConfig.HTTPClient = global.HttpClient
		global.AIOpenaiClient = openai.NewClientWithConfig(openaiClientConfig)
		log.Println("AIOpenaiClient init success")
	}

	// 处理通用的AI客户端（基于配置文件启动哪一个）
	switch global.Config.AI.Enable {
	case "openai":
		global.AIClient = global.AIOpenaiClient
	case "azure":
		global.AIClient = global.AIAzureClient
	}

	var AIClientText = "AIClient init success（base: %s）\n"
	if global.AIClient == nil {
		AIClientText = "AIClient init failed（base: %s）\n"
	}
	log.Printf(AIClientText, global.Config.AI.Enable)
}
