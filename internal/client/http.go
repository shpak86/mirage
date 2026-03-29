package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/enetx/g"
	"github.com/enetx/http"
	"github.com/enetx/surf"
)

type Request struct {
	Fingerprint string
	Method      string
	Url         string
	Headers     map[string][]string
	Cookies     []string
	Body        []byte
}

type Response struct {
	RawRequest  *http.Request
	RawResponse *http.Response
}

type HttpClient struct {
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) Do(request Request) (response *Response, err error) {
	builder := surf.NewClient().Builder()
	browser, Os, _ := strings.Cut(request.Fingerprint, "-")
	if browser == "" {
		browser = "chrome"
	}
	if Os == "" {
		Os = "windows"
	}
	var imporsonate *surf.Impersonate
	// OS
	switch Os {
	case "android":
		imporsonate = builder.Impersonate().Android()
	case "windows":
		imporsonate = builder.Impersonate().Windows()
	case "linux":
		imporsonate = builder.Impersonate().Linux()
	case "macos":
		imporsonate = builder.Impersonate().MacOS()
	default:
		return nil, fmt.Errorf("unknown os %s", Os)
	}

	// Browser
	switch browser {
	case "chrome":
		imporsonate.Chrome()
		builder.JA().Chrome()
	case "chrome58":
		imporsonate.Chrome()
		builder.JA().Chrome58()
	case "chrome62":
		imporsonate.Chrome()
		builder.JA().Chrome62()
	case "chrome70":
		imporsonate.Chrome()
		builder.JA().Chrome70()
	case "chrome72":
		imporsonate.Chrome()
		builder.JA().Chrome72()
	case "chrome83":
		imporsonate.Chrome()
		builder.JA().Chrome83()
	case "chrome87":
		imporsonate.Chrome()
		builder.JA().Chrome87()
	case "chrome96":
		imporsonate.Chrome()
		builder.JA().Chrome96()
	case "chrome100":
		imporsonate.Chrome()
		builder.JA().Chrome100()
	case "chrome102":
		imporsonate.Chrome()
		builder.JA().Chrome102()
	case "chrome106":
		imporsonate.Chrome()
		builder.JA().Chrome106()
	case "chrome120":
		imporsonate.Chrome()
		builder.JA().Chrome120()
	case "chrome145":
		imporsonate.Chrome()
		builder.JA().Chrome145()
	case "firefox":
		imporsonate.Firefox()
		builder.JA().Firefox()
	case "firefox55":
		imporsonate.Firefox()
		builder.JA().Firefox55()
	case "firefox56":
		imporsonate.Firefox()
		builder.JA().Firefox56()
	case "firefox63":
		imporsonate.Firefox()
		builder.JA().Firefox63()
	case "firefox65":
		imporsonate.Firefox()
		builder.JA().Firefox65()
	case "firefox99":
		imporsonate.Firefox()
		builder.JA().Firefox99()
	case "firefox102":
		imporsonate.Firefox()
		builder.JA().Firefox102()
	case "firefox105":
		imporsonate.Firefox()
		builder.JA().Firefox105()
	case "firefox120":
		imporsonate.Firefox()
		builder.JA().Firefox120()
	case "firefox148":
		imporsonate.Firefox()
		builder.JA().Firefox148()
	default:
		return nil, fmt.Errorf("unknown browser %s", browser)
	}

	surfClient := builder.Build().Unwrap()
	var req *surf.Request
	switch strings.ToLower(request.Method) {
	case "get":
		req = surfClient.Get(g.String(request.Url))
	case "post":
		req = surfClient.Post(g.String(request.Url))
	case "put":
		req = surfClient.Put(g.String(request.Url))
	case "options":
		req = surfClient.Options(g.String(request.Url))
	case "delete":
		req = surfClient.Delete(g.String(request.Url))
	case "head":
		req = surfClient.Head(g.String(request.Url))
	case "patch":
		req = surfClient.Patch(g.String(request.Url))
	case "connect":
		req = surfClient.Connect(g.String(request.Url))
	default:
		req = surfClient.Get(g.String(request.Url))
	}
	// Headers
	for k, v := range request.Headers {
		for i := range v {
			req.AddHeaders(k, v[i])
		}
	}
	//Cookies
	if len(request.Cookies) != 0 {
		cookies := make([]*http.Cookie, 0, len(request.Cookies))
		for _, it := range request.Cookies {
			name, value, found := strings.Cut(it, "=")
			if found {
				cookies = append(cookies, &http.Cookie{Name: name, Value: value})
			} else {
				fmt.Fprintf(os.Stderr, "wrong cookie %s", it)
			}
		}
		req.AddCookies(cookies...)
	}
	if len(request.Body) != 0 {
		req.Body(request.Body)
	}
	// Do a request
	result := req.Do()
	resp, err := result.Result()
	if err != nil {
		return nil, err
	}

	response = &Response{
		RawResponse: resp.GetResponse(),
		RawRequest:  resp.GetResponse().Request,
	}

	return
}
