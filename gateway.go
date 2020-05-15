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

var apiGateway = gateway{
	config: Config{},
}

type gateway struct {
	config Config
}

func (g gateway) sendRequest(method string, path string, data interface{}, writeTo interface{}) error {

	// Build request
	client := http.Client{}
	url := buildRequestURL(method, path, data, g.config)
	var body io.Reader
	if method == http.MethodPost {
		body = buildRequestBody(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(APIKey, "")
	req.Header.Set("Content-Type", contentType)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Parse response
	return parseResponse(resp, writeTo)
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
		query = urlEncodeData(data)
	}

	url := fmt.Sprintf("%v://%v:%v%v%v%v", protocol, host, port, basePath, path, query)
	return url
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
	encData := ""
	if data != nil {
		rVal := reflect.ValueOf(data)
		rType := reflect.TypeOf(data)
		urlVals := netUrl.Values{}
		for i := 0; i < rVal.NumField(); i++ {
			if rVal.Field(i).Kind() != reflect.Ptr || !rVal.Field(i).IsNil() {
				if tag := rType.Field(i).Tag.Get("json"); tag != "" {
					var f = rVal.Field(i)
					if f.Kind() == reflect.Ptr {
						f = f.Elem()
					}
					urlVals.Add(tag, fmt.Sprintf("%v", f.Interface()))
				}
			}
		}
		encData = urlVals.Encode()
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

	if obj.Object == "error" {
		var apiError Error
		if err = json.Unmarshal(bytes, &apiError); err != nil {
			return err
		}
		return apiError
	}

	return json.Unmarshal(bytes, into)
}
