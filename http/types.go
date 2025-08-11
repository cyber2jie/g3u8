package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
)

type HttpHeader map[string]string

func (h HttpHeader) Get(key string) string {
	header := h[key]
	return header
}
func (h HttpHeader) Set(key string, value string) {
	h[key] = value
}

type HttpRequestBody interface {
	Body() io.Reader
	ContentType() string
}

// http entity
type HttpRequestStringBody struct {
	BodyString        string
	ContentTypeString string
}

func (b *HttpRequestStringBody) Body() io.Reader {
	return strings.NewReader(b.BodyString)
}
func (b *HttpRequestStringBody) ContentType() string {
	return b.ContentTypeString
}

type HttpRequestMultipartBody struct {
	writer *multipart.Writer
	buffer *bytes.Buffer
}

func NewMultipartBody() *HttpRequestMultipartBody {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	return &HttpRequestMultipartBody{
		writer: writer,
		buffer: buffer,
	}
}

func (b *HttpRequestMultipartBody) Body() io.Reader {
	b.writer.Close()
	return b.buffer
}
func (b *HttpRequestMultipartBody) ContentType() string {
	return b.writer.FormDataContentType()
}

func (b *HttpRequestMultipartBody) AddFile(fieldName string, fileName string, file io.Reader) error {
	part, err := b.writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}
func (b *HttpRequestMultipartBody) AddField(fieldName string, fieldValue string) error {
	return b.writer.WriteField(fieldName, fieldValue)
}

// httprequest,httpresponse

type HttpRequest struct {
	Method string
	Url    string
	Header HttpHeader
	Body   HttpRequestBody
}

type HttpResponse struct {
	StatusCode int
	Proto      string
	Header     HttpHeader
	Body       []byte
}
