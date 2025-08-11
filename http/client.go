package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClientCreateError struct {
	error
	message string
}

func (e *HttpClientCreateError) Error() string {
	return e.message
}

// httpclient
type HttpClient struct {
	client *http.Client
	option *HttpClientOption
}

type HttpClientOption struct {
	BaseUrl string
	Timeout int
	Headers map[string]string
	Cookies []*http.Cookie
	Proxy   string
}

func newClient(option *HttpClientOption) (*http.Client, error) {

	transport := http.DefaultTransport
	if option.Proxy != "" {
		_url, err := url.Parse(option.Proxy)
		if err != nil {
			return nil, &HttpClientCreateError{message: fmt.Sprintf("Invalid proxy url: %s,error: %s ", option.Proxy, err.Error())}
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(_url),
		}

	}

	return &http.Client{
		Timeout:   time.Duration(option.Timeout) * time.Second,
		Transport: transport,
	}, nil
}

func New(option *HttpClientOption) (*HttpClient, error) {
	client, err := newClient(option)
	if err != nil {
		return nil, err
	}
	return &HttpClient{
		client: client,
		option: option,
	}, nil
}

func (c *HttpClient) buildUrl(url string) string {
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		return url
	}
	url = strings.TrimPrefix(url, "/")
	return fmt.Sprintf("%s/%s", c.option.BaseUrl, url)
}

func (c *HttpClient) Do(request *HttpRequest) (*HttpResponse, error) {
	return c.DoWithContext(context.Background(), request)
}

func (c *HttpClient) DoWithContext(ctx context.Context, request *HttpRequest) (*HttpResponse, error) {
	var body io.Reader
	var contentType string
	if request.Body != nil {
		body = request.Body.Body()
		contentType = request.Body.ContentType()
	}
	req, err := http.NewRequestWithContext(ctx, request.Method, c.buildUrl(request.Url), body)
	if err != nil {
		return nil, err
	}
	if c.option.Headers != nil {
		for k, v := range c.option.Headers {
			if k != "" && v != "" {
				req.Header.Set(k, v)
			}
		}
	}

	if request.Header != nil {
		for k, v := range request.Header {
			if k != "" && v != "" {
				req.Header.Set(k, v)
			}
		}
	}

	if contentType != "" {
		req.Header.Set(Header_Content_Type, contentType)
	}
	if c.option.Cookies != nil {
		for _, cookie := range c.option.Cookies {
			req.AddCookie(cookie)
		}
	}

	response, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	httpResponse, err := wrapResponse(response)

	if err != nil {
		return nil, err
	}

	return httpResponse, nil
}

// Do Method

func (c *HttpClient) Get(url string) (*HttpResponse, error) {
	return c.Do(&HttpRequest{
		Method: Http_Get,
		Url:    c.buildUrl(url),
	})
}
func (c *HttpClient) Post(url string, body HttpRequestBody) (*HttpResponse, error) {
	return c.Do(&HttpRequest{
		Method: Http_Post,
		Url:    c.buildUrl(url),
		Body:   body,
	})
}
func (c *HttpClient) PostForm(url string, data url.Values) (*HttpResponse, error) {
	return c.Do(&HttpRequest{
		Method: Http_Post,
		Url:    c.buildUrl(url),
		Body: &HttpRequestStringBody{
			data.Encode(),
			Content_Type_Form,
		},
	})
}

func wrapResponse(response *http.Response) (*HttpResponse, error) {
	var headers = make(HttpHeader)

	for k, v := range response.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	defer response.Body.Close()

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, response.Body)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		Body:       buffer.Bytes(),
		Header:     headers,
		Proto:      response.Proto,
		StatusCode: response.StatusCode,
	}, nil
}
