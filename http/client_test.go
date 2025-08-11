package http

import "testing"

func TestRequest(t *testing.T) {
	client, _ := New(
		&HttpClientOption{
			BaseUrl: "https://bing.com",
			Timeout: 60,
		},
	)
	resp, _ := client.Get("search?q=go+pkg")

	t.Logf("%s", string(resp.Body))

}

func TestPost(t *testing.T) {
	client, _ := New(
		&HttpClientOption{
			BaseUrl: "https://bing.com",
			Timeout: 60,
		},
	)
	resp, _ := client.Post("search", &HttpRequestStringBody{
		BodyString:        "q=go+pkg",
		ContentTypeString: Content_Type_Form,
	})

	t.Logf("%s", string(resp.Body))

}
