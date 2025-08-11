package download

import "g3u8/http"

var client *http.HttpClient

func GetClient() *http.HttpClient {
	if client == nil {
		_client, err := http.New(&http.HttpClientOption{})
		if err != nil {
			panic(err)
		}
		client = _client
	}
	return client
}
