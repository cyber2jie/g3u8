package download

import (
	"context"
	"g3u8/config"
	"g3u8/http"
	"log"
)

var defaultM3u8DownloadContext *M3u8DownloadContext

func InitContext(ctx context.Context, g3u8Config *config.G3u8Config, options M3u8DownloadOptions, outDir string) {

	client, err := getHttpClientFromConfig(g3u8Config, options)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defaultM3u8DownloadContext = &M3u8DownloadContext{
		Config:  g3u8Config,
		Client:  client,
		Options: options,
		Context: ctx,
		OutDir:  outDir,
	}
}

func getHttpClientFromConfig(g3u8Config *config.G3u8Config, options M3u8DownloadOptions) (*http.HttpClient, error) {

	headers := make(map[string]string)
	var baseUrl, proxy string
	var timeout int

	if g3u8Config.Http.Timeout > 0 {
		timeout = g3u8Config.Http.Timeout
	}
	if g3u8Config.Http.Headers != nil && len(g3u8Config.Http.Headers) > 0 {
		for _, header := range g3u8Config.Http.Headers {
			if header.Name == "" || header.Value == "" {
				continue
			}
			headers[header.Name] = header.Value
		}
	}

	if g3u8Config.Proxy.Enable {
		if g3u8Config.Proxy.Proxy != nil && *g3u8Config.Proxy.Proxy != "" {
			proxy = *g3u8Config.Proxy.Proxy
			log.Printf("Using proxy: %s", proxy)
		}
	}

	var clientOption = &http.HttpClientOption{
		BaseUrl: baseUrl,
		Timeout: timeout,
		Headers: headers,
		Proxy:   proxy,
	}

	client, err := http.New(clientOption)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetDefaultM3u8DownloadContext() *M3u8DownloadContext {
	return defaultM3u8DownloadContext
}

type M3u8DownloadContext struct {
	Config  *config.G3u8Config
	Client  *http.HttpClient
	Options M3u8DownloadOptions
	Context context.Context
	OutDir  string
}
