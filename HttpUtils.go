package utils

import (
	"net/url"
	"sync"
	"time"

	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var (
	once      sync.Once
	Client    *fasthttp.PipelineClient
	ClientMap sync.Map
)

const (
	TIMEOUTS                  = 3
	ClientMaxConns            = 500
	ClientMaxPendingRequests  = 1024
	ClientMaxIdleConnDuration = 10 * time.Second
)

func DoHttpTimeout(array map[string]interface{}, method string, params ...string) ([]byte, int, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	if len(params) == 0 {
		return response.Body(), response.StatusCode(), nil
	}
	request.Header.SetMethod(method)
	var arg = &fasthttp.Args{}
	urls := params[0]
	if len(array) > 0 {
		for k, v := range array {
			arg.Set(k, StringMust(v))
		}
	}
	switch method {
	case "GET":
		if arg.Len() > 0 {
			urls = urls + "?" + arg.String()
		}
	case "POST":
		arg.WriteTo(request.BodyWriter())
		request.Header.SetContentType("application/x-www-form-urlencoded")
		if len(params) == 2 {
			bodyJson, _ := jsoniter.Marshal(array)
			request.SetBodyString(string(bodyJson))
			request.Header.SetContentType("application/json;charset=utf-8")
		}
	}

	request.Header.Set("Connection", "keep-alive")
	request.SetRequestURI(urls)

	pipelineClient := getPipelineClient(urls)
	for {
		if err := pipelineClient.DoTimeout(request, response, TIMEOUTS*1000*time.Millisecond); err != nil {
			if err == fasthttp.ErrPipelineOverflow {
				continue
			}
		}
		break
	}
	return response.Body(), response.StatusCode(), nil
}

func newClientAndSetToMap(url string) *fasthttp.PipelineClient {
	if Client == nil {
		once.Do(func() {
			host, scheme := getHostFroURL(url)
			Client = &fasthttp.PipelineClient{
				Addr:                host,
				MaxConns:            ClientMaxConns,
				MaxPendingRequests:  ClientMaxPendingRequests,
				MaxIdleConnDuration: ClientMaxIdleConnDuration,
				IsTLS:               On(scheme == "http", false, true).(bool),
			}
			ClientMap.Store(url, Client)
			ClientMap.Store(host, Client)
		})
	}
	return Client
}

func getHostFroURL(urls string) (string, string) {
	u, _ := url.Parse(urls)
	return u.Host, u.Scheme
}
func getClientFromMap(key string) *fasthttp.PipelineClient {
	if value, ok := ClientMap.Load(key); ok {
		return value.(*fasthttp.PipelineClient)
	}
	return Client
}
func getPipelineClient(url string) *fasthttp.PipelineClient {
	host, _ := getHostFroURL(url)
	client := getClientFromMap(host)
	if client != nil {
		return client
	}
	client = newClientAndSetToMap(url)
	return client
}
