package core

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InitHeaders(options Options) []map[string]string {
	var headers []map[string]string
	head := map[string]string{
		"UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	}
	headers = append(headers, head)
	headers = append(headers, map[string]string{"Accept": "*/*"})
	headers = append(headers, map[string]string{"AcceptLang": "en-US,en;q=0.8"})
	return headers
}

// SendGET just send GET request
func SendGET(url string, options Options) (Request, Response) {
	req := Request{
		Method:   "GET",
		URL:      url,
		Headers:  InitHeaders(options),
		Redirect: options.Redirect,
	}

	resp, _ := JustSend(options, req)
	return req, resp
}

// SendPOST just send POST request
func SendPOST(url string, options Options) (Request, Response) {
	req := Request{
		Method:   "POST",
		URL:      url,
		Headers:  InitHeaders(options),
		Redirect: options.Redirect,
	}

	resp, _ := JustSend(options, req)
	return req, resp
}

// JustSend just sending request
func JustSend(options Options, req Request) (res Response, err error) {
	if req.Method == "" {
		req.Method = "GET"
	}
	method := req.Method
	url := req.URL
	body := req.Body
	headers := GetHeaders(req)

	timeout := options.Timeout
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	// update it again
	var newHeader []map[string]string
	for k, v := range headers {
		element := make(map[string]string)
		element[k] = v
		newHeader = append(newHeader, element)
	}
	req.Headers = newHeader

	// disable log when retry
	logger := logrus.New()
	if !options.Debug {
		logger.Out = ioutil.Discard
	}

	client := resty.New()
	client.SetLogger(logger)
	client.SetTransport(&http.Transport{
		MaxIdleConns:          100,
		MaxConnsPerHost:       1000,
		IdleConnTimeout:       time.Duration(timeout) * time.Second,
		ExpectContinueTimeout: time.Duration(timeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
		TLSHandshakeTimeout:   time.Duration(timeout) * time.Second,
		DisableCompression:    true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	})

	client.SetHeaders(headers)
	client.SetCloseConnection(true)
	if options.Proxy != "" {
		client.SetProxy(options.Proxy)
	}
	// override proxy
	if req.Proxy != "" && req.Proxy != "blank" {
		client.SetProxy(req.Proxy)
	}
	if options.Retry > 0 {
		client.SetRetryCount(options.Retry)
	}
	client.SetTimeout(time.Duration(timeout) * time.Second)
	client.SetRetryWaitTime(time.Duration(timeout/2) * time.Second)
	client.SetRetryMaxWaitTime(time.Duration(timeout) * time.Second)
	timeStart := time.Now()
	// redirect policy
	if req.Redirect == false {
		client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			// keep the header the same
			// client.SetHeaders(headers)

			res.StatusCode = req.Response.StatusCode
			res.Status = req.Response.Status
			resp := req.Response
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ErrorF("%v", err)
			}
			bodyString := string(bodyBytes)
			resLength := len(bodyString)
			// format the headers
			var resHeaders []map[string]string
			for k, v := range resp.Header {
				element := make(map[string]string)
				//fmt.Printf("%v: %v\n", k, v)
				element[k] = strings.Join(v[:], "")
				resLength += len(fmt.Sprintf("%s: %s\n", k, strings.Join(v[:], "")))
				if k == "Location" {
					res.Location = strings.Join(v[:], "")
				}
				if k == "Set-Cookie" {
					res.Cookies = strings.Join(v[:], " ")
				}
				resHeaders = append(resHeaders, element)
			}

			// response time in second
			resTime := time.Since(timeStart).Seconds()
			resHeaders = append(resHeaders,
				map[string]string{"Total Length": strconv.Itoa(resLength)},
				map[string]string{"Response Time": fmt.Sprintf("%f", resTime)},
			)

			// set some variable
			res.Headers = resHeaders
			res.StatusCode = resp.StatusCode
			res.Status = fmt.Sprintf("%v %v", resp.Status, resp.Proto)
			res.Body = bodyString
			res.ResponseTime = resTime
			res.Length = resLength
			// beautify
			res.Beautify = BeautifyResponse(res)
			return errors.New("auto redirect is disabled")
		}))

		client.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return false
			},
		)
	} else {
		client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			// keep the header the same
			client.SetHeaders(headers)
			return nil
		}))
	}

	var resp *resty.Response
	// really sending things here
	method = strings.ToLower(strings.TrimSpace(method))
	switch method {
	case "get":
		resp, err = client.R().
			SetBody([]byte(body)).
			Get(url)
		break
	case "post":
		resp, err = client.R().
			SetBody([]byte(body)).
			Post(url)
		break
	case "head":
		resp, err = client.R().
			SetBody([]byte(body)).
			Head(url)
		break
	case "options":
		resp, err = client.R().
			SetBody([]byte(body)).
			Options(url)
		break
	case "patch":
		resp, err = client.R().
			SetBody([]byte(body)).
			Patch(url)
		break
	case "put":
		resp, err = client.R().
			SetBody([]byte(body)).
			Put(url)
		break
	case "delete":
		resp, err = client.R().
			SetBody([]byte(body)).
			Delete(url)
		break
	}

	// in case we want to get redirect stuff
	if res.StatusCode != 0 {
		return res, nil
	}

	if err != nil || resp == nil {
		ErrorF("%v %v", url, err)
		return Response{}, err
	}

	return ParseResponse(*resp), nil
}

// ParseResponse field to Response
func ParseResponse(resp resty.Response) (res Response) {
	// var res libs.Response
	resLength := len(string(resp.Body()))
	// format the headers
	var resHeaders []map[string]string
	for k, v := range resp.RawResponse.Header {
		element := make(map[string]string)
		element[k] = strings.Join(v[:], "")
		resLength += len(fmt.Sprintf("%s: %s\n", k, strings.Join(v[:], "")))
		if k == "Location" {
			res.Location = strings.Join(v[:], "")
		}
		if k == "Set-Cookie" {
			res.Cookies = strings.Join(v[:], " ")
		}
		resHeaders = append(resHeaders, element)
	}
	// response time in second
	resTime := float64(resp.Time()) / float64(time.Second)
	resHeaders = append(resHeaders,
		map[string]string{"Total Length": strconv.Itoa(resLength)},
		map[string]string{"Response Time": fmt.Sprintf("%f", resTime)},
	)

	// set some variable
	res.Headers = resHeaders
	res.StatusCode = resp.StatusCode()
	res.Status = fmt.Sprintf("%v %v", resp.Status(), resp.RawResponse.Proto)
	res.Body = string(resp.Body())
	res.ResponseTime = resTime
	res.Length = resLength
	// beautify
	res.Beautify = BeautifyResponse(res)
	res.BeautifyHeader = BeautifyHeaders(res)
	return res
}

// GetHeaders generate headers if not provide
func GetHeaders(req Request) map[string]string {
	// random user agent
	UserAgens := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3941.0 Safari/537.36",
		"Mozilla/5.0 (X11; U; Windows NT 6; en-US) AppleWebKit/534.12 (KHTML, like Gecko) Chrome/9.0.587.0 Safari/534.12",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
	}

	headers := make(map[string]string)
	if len(req.Headers) == 0 {
		rand.Seed(time.Now().Unix())
		headers["User-Agent"] = UserAgens[rand.Intn(len(UserAgens))]
		return headers
	}

	for _, header := range req.Headers {
		for key, value := range header {
			headers[key] = value
		}
	}

	rand.Seed(time.Now().Unix())
	// append user agent in case you didn't set user-agent
	if headers["User-Agent"] == "" {
		rand.Seed(time.Now().Unix())
		headers["User-Agent"] = UserAgens[rand.Intn(len(UserAgens))]
	}
	return headers
}

// BeautifyHeaders beautify response headers
func BeautifyHeaders(res Response) string {
	beautifyHeader := fmt.Sprintf("%v \n", res.Status)
	for _, header := range res.Headers {
		for key, value := range header {
			beautifyHeader += fmt.Sprintf("%v: %v\n", key, value)
		}
	}
	return beautifyHeader
}

// BeautifyResponse beautify response
func BeautifyResponse(res Response) string {
	var beautifyRes string
	beautifyRes += fmt.Sprintf("%v \n", res.Status)

	for _, header := range res.Headers {
		for key, value := range header {
			beautifyRes += fmt.Sprintf("%v: %v\n", key, value)
		}
	}

	beautifyRes += fmt.Sprintf("\n%v\n", res.Body)
	return beautifyRes
}
