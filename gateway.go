package feather

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"reflect"
	"strings"
)

const (
	contentType     string = "application/x-www-form-urlencoded"
	defaultProtocol string = "https"
	defaultHost     string = "api.feather.id"
	defaultPort     string = "443"
	defaultBasePath string = "/v1"
)

type gateway struct {
	apiKey string
	config Config
}

func (g gateway) sendRequest(method string, path string, data interface{}, writeTo interface{}) error {
	req, err := g.buildRequest(method, path, data)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return parseResponse(resp, writeTo)
}

func (g gateway) buildRequest(method string, path string, data interface{}) (*http.Request, error) {
	url := buildRequestURL(method, path, data, g.config)
	var body io.Reader
	if method == http.MethodPost {
		body = buildRequestBody(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(g.apiKey, "")
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func buildRequestURL(method string, path string, data interface{}, cfg Config) string {
	protocol := defaultProtocol
	if cfg.Protocol != "" {
		protocol = cfg.Protocol
	}
	host := defaultHost
	if cfg.Host != "" {
		host = cfg.Host
	}
	port := defaultPort
	if cfg.Port != "" {
		port = cfg.Port
	}
	basePath := defaultBasePath
	if cfg.BasePath != "" {
		basePath = cfg.BasePath
	}
	query := ""
	if method == http.MethodGet {
		query = "?" + urlEncodeData(data)
	}
	return fmt.Sprintf("%v://%v:%v%v%v%v", protocol, host, port, basePath, path, query)
}

func buildRequestBody(data interface{}) io.Reader {
	encData := urlEncodeData(data)
	var body io.Reader
	if encData != "" {
		body = strings.NewReader(encData)
	}
	return body
}

func urlEncodeData(data interface{}) string {
	joinEncData := func(a string, b string) string {
		if a == "" {
			return b
		} else if b == "" {
			return a
		}
		return a + "&" + b
	}
	encData := ""
	if data != nil {
		rVal := reflect.ValueOf(data)
		rType := reflect.TypeOf(data)
		urlVals := netUrl.Values{}
		for i := 0; i < rVal.NumField(); i++ {
			if rVal.Field(i).Kind() != reflect.Ptr || !rVal.Field(i).IsNil() {
				if rVal.Field(i).Kind() == reflect.Struct {
					// Recursively flatten the request
					encData = joinEncData(encData, urlEncodeData(rVal.Field(i).Interface()))
				} else if tag := rType.Field(i).Tag.Get("json"); tag != "" {
					var f = rVal.Field(i)
					if f.Kind() == reflect.Ptr {
						f = f.Elem()
					}
					urlVals.Add(tag, fmt.Sprintf("%v", f.Interface()))
				}
			}
		}
		encData = joinEncData(encData, urlVals.Encode())
	}
	return encData
}

func parseResponse(resp *http.Response, into interface{}) error {
	type object struct {
		Object string `json:"object"`
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var obj object
	if err = json.Unmarshal(bytes, &obj); err != nil {
		return err
	}
	const objectError = "error"
	if obj.Object == objectError {
		var ferr Error
		if err = json.Unmarshal(bytes, &ferr); err != nil {
			return err
		}
		return ferr
	}
	return json.Unmarshal(bytes, into)
}
