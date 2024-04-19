package initialize

import (
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"net/url"
	"time"
	"vland.live/app/global"
)

func initOpenai() {
	var httpProxyAddress string
	if global.Config.Proxy != nil {
		httpProxyAddress = global.Config.Proxy.Http
	}
	// 微软的openai客户端初始化
	{
		azureClientConfig := openai.DefaultAzureConfig(global.Config.AI.Azure.Key, global.Config.AI.Azure.BaseUrl)
		azureClientConfig.AzureModelMapperFunc = func(mode string) string {
			azureModelMapping := map[string]string{
				openai.GPT3Dot5Turbo:           "xxx", // model name
				openai.AdaEmbeddingV2.String(): "xxx", // model name
			}
			return azureModelMapping[mode]
		}

		// 设置代理
		if httpProxyAddress != "" {
			proxyURL, _ := url.Parse(httpProxyAddress)
			azureClientConfig.HTTPClient = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
			}
		} else {
			azureClientConfig.HTTPClient = &http.Client{
				Timeout: 60 * time.Second,
			}
		}
		global.AIAzureClient = openai.NewClientWithConfig(azureClientConfig)
		log.Println("AIAzureClient init success")
	}

	// openai 的客户端初始化
	{
		openaiClientConfig := openai.DefaultConfig(global.Config.AI.Openai.Key)
		if httpProxyAddress != "" {
			proxy, _ := url.Parse(httpProxyAddress)
			httpClient := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxy),
				},
			}
			openaiClientConfig.HTTPClient = httpClient
		} else {
			openaiClientConfig.HTTPClient = &http.Client{Timeout: 60 * time.Second}
			global.Logger.Info("openaiClientConfig.HTTPClient超时时间已设置")
		}
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
