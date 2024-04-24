package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

var (
	methodGet  = "GET"
	methodPost = "POST"
	methodPut  = "PUT"
	methodDel  = "DELETE"

	//reqTimeOut = 30 * time.Second
	reqTimeOut = 15 * time.Second
)

func HttpGet(URL string) (*fasthttp.Response, error) {
	return HttpGetWithHeader(URL, nil)
}

func HttpGetWithHeader(URL string, header map[string]string) (*fasthttp.Response, error) {
	return httpReq(URL, methodGet, "", header, nil, reqTimeOut)
}

func HttpDel(URL string) (*fasthttp.Response, error) {
	return HttpDelWithHeader(URL, nil)
}

func HttpDelWithHeader(URL string, header map[string]string) (*fasthttp.Response, error) {
	return httpReq(URL, methodDel, "", header, nil, reqTimeOut)
}

func HttpPost(URL string, data interface{}) (*fasthttp.Response, error) {
	return HttpPostWithHeader(URL, data, nil)
}

func HttpPostWithHeader(URL string, data interface{}, header map[string]string) (*fasthttp.Response, error) {
	return httpReq(URL, methodPost, "", header, data, reqTimeOut)
}

func httpReq(url, method, contentType string, header map[string]string, data interface{}, timeout time.Duration) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		//fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetConnectionClose()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	if contentType != "" {
		req.Header.SetContentType(contentType)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	if data != nil {
		dataByte := make([]byte, 0)
		switch d := data.(type) {
		case map[string]string:
			args := &fasthttp.Args{}
			for k, v := range d {
				args.Add(k, v)
			}
			dataByte = []byte(args.String())
		case []byte:
			dataByte = d
		default:
			dataByte = []byte(data.(string))
		}
		req.SetBody(dataByte)
	}

	client := &fasthttp.Client{
		//TLSConfig:                     &tls.Config{InsecureSkipVerify: true},
		//NoDefaultUserAgentHeader:      true,
		MaxConnsPerHost:     10000,
		ReadTimeout:         timeout,
		WriteTimeout:        timeout,
		MaxIdleConnDuration: time.Minute,
		//DisableHeaderNamesNormalizing: true,
	}
	//if global.GVA_CONFIG.System.Env == "" || global.GVA_CONFIG.System.Env == "local" {
	//	client.Dial = FasthttpHTTPDialer("127.0.0.1:7890")
	//}

	if err := client.DoRedirects(req, resp, 5); err != nil {
		return nil, errors.New(fmt.Sprintf("HTTP %s fail url=%s Error=%s", method, url, err.Error()))
	}
	statusCode := resp.StatusCode()
	//if statusCode != 200 && statusCode != 302 && statusCode != 304 && statusCode != 206 {
	if statusCode != 200 && statusCode != 302 {
		return resp, errors.New(fmt.Sprintf("%d", statusCode))
	}
	return resp, nil
}

func FasthttpHTTPDialer(proxyAddr string) fasthttp.DialFunc {
	return func(addr string) (net.Conn, error) {
		conn, err := fasthttp.Dial(proxyAddr)
		if err != nil {
			return nil, err
		}

		req := "CONNECT " + addr + " HTTP/1.1\r\n"
		// req += "Proxy-Authorization: xxx\r\n"
		req += "\r\n"

		if _, err := conn.Write([]byte(req)); err != nil {
			return nil, err
		}

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		res.SkipBody = true

		if err := res.Read(bufio.NewReader(conn)); err != nil {
			conn.Close()
			return nil, err
		}
		if res.Header.StatusCode() != 200 {
			conn.Close()
			return nil, fmt.Errorf("could not connect to proxy")
		}
		return conn, nil
	}
}
