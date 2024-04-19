package initialize

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"
	"vland.live/app/global"
)

func initHttp(ctx context.Context) {

	var httpProxyAddress string
	if global.Config.Proxy != nil {
		httpProxyAddress = global.Config.Proxy.Http
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	if httpProxyAddress != "" {
		proxyURL, _ := url.Parse(httpProxyAddress)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}
	global.HttpClient = httpClient
	log.Println("HttpClient init success")
}
