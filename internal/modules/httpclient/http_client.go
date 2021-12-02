package httpclient

// http-client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

type MarkDownModel struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

type DingMsg struct {
	MsgType  string   `json:"msgtype,omitempty"`
	Markdown MarkDownModel `json:"markdown,omitempty"`
}

func Get(url string, timeout int) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}

	return request(req, timeout)
}

func PostParams(url string, params string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return request(req, timeout)
}

func PostJson(url string, body interface{}, timeout int) ResponseWrapper {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return request(req, timeout)
}

func request(req *http.Request, timeout int) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		logger.Errorf("http response: %+v", wrapper)
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header
	logger.Infof("http response: %+v", wrapper)
	return wrapper
}

func setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "golang/gocron")
}

func createRequestError(err error) ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return ResponseWrapper{0, errorMessage, make(http.Header)}
}
