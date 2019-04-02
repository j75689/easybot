package plugin

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/j75689/easybot/pkg/util"
	"go.uber.org/zap"
)

type CurlPluginConfig struct {
	URL         string            `json:"url"`
	PostData    string            `json:"data"`
	Method      string            `json:"method"`
	Header      map[string]string `json:"header"`
	ProxyURL    string            `json:"proxy"`
	Timeout     int               `json:"timeout"`
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	Insecure    bool              `json:"insecure"`
	NoKeepalive bool              `json:"no-keepalive"`
	Location    bool              `json:"location"`

	Output *CurlPluginOutputConfig `json:"output"`
}

type CurlPluginOutputConfig struct {
	ParseDriver string            `json:"parser"`
	Variables   map[string]string `json:"variables"`
}

func Curl(input interface{}, variables map[string]interface{}, logger *zap.SugaredLogger) (map[string]interface{}, bool, error) {
	logger.Info("Excute Graphql Plugin")
	var (
		config = CurlPluginConfig{
			Method:      http.MethodGet,
			Header:      make(map[string]string),
			Timeout:     60,
			Insecure:    false,
			NoKeepalive: true,
			Location:    true,
			Output:      new(CurlPluginOutputConfig),
		}
		next = true
		err  error
	)

	err = json.Unmarshal(util.GetJSONBytes(input), &config)
	if err != nil {
		logger.Error(err)
		return nil, next, err
	}

	// Client
	var (
		proxy           func(*http.Request) (*url.URL, error)
		noCheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	)
	logger.Debug("timeout: ", config.Timeout)
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}
	if !config.Location {
		client.CheckRedirect = noCheckRedirect
	}

	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		logger.Debug("proxyURL: ", proxyURL)
		if err != nil {
			logger.Error("parse ProxyURL error: ", err)
		}
		proxy = http.ProxyURL(proxyURL)

	} else {
		proxy = http.ProxyFromEnvironment
	}
	client.Transport = &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
		DisableKeepAlives: config.NoKeepalive,
	}

	// Request
	req, err := http.NewRequest(config.Method, config.URL, strings.NewReader(config.PostData))
	if err != nil {
		logger.Error("NewRequest", err)
	}
	for k, v := range config.Header {
		req.Header.Set(k, v)
	}
	if config.PostData != "" {
		if strings.HasPrefix(config.PostData, "{") {
			req.Header.Set("Content-Type", "application/json")
		} else if strings.HasPrefix(config.PostData, "<") {
			req.Header.Set("Content-Type", "application/xml")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	logger.Debug(req)

	response, err := client.Do(req)
	if err != nil {
		logger.Error("request error: ", err)
		return variables, next, err
	}
	defer response.Body.Close()

	var body []byte
	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch strings.ToUpper(response.Header.Get("Content-Encoding")) {
	case "GZIP":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err = ioutil.ReadAll(reader)
	if err != nil {
		logger.Error("client error: ", err)
		return variables, next, err
	}

	switch strings.ToUpper(config.Output.ParseDriver) {
	case "JSON":
		var data map[string]interface{}
		if err = json.Unmarshal(body, &data); err == nil {
			for k, v := range config.Output.Variables {
				variables[k] = util.GetJSONValue(v, data)
			}
		}
	case "TEXT":
		data := string(body)
		for k, v := range config.Output.Variables {
			r := regexp.MustCompile(v)
			variables[k] = r.FindStringSubmatch(data)[1]
		}
	}

	return variables, next, err
}
