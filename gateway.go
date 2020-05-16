package feather

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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
	client *http.Client
}

func (g gateway) sendRequest(method string, path string, data interface{}, writeTo interface{}) error {
	req, err := g.buildRequest(method, path, data)
	if err != nil {
		return err
	}
	resp, err := g.getClient().Do(req)
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

func (g gateway) getClient() *http.Client {
	if g.config.HTTPClient != nil {
		return g.config.HTTPClient
	}
	return http.DefaultClient
}

func buildRequestURL(method string, path string, data interface{}, cfg Config) string {
	protocol := defaultProtocol
	if cfg.Protocol != nil {
		protocol = *cfg.Protocol
	}
	host := defaultHost
	if cfg.Host != nil {
		host = *cfg.Host
	}
	port := defaultPort
	if cfg.Port != nil {
		port = *cfg.Port
	}
	basePath := defaultBasePath
	if cfg.BasePath != nil {
		basePath = *cfg.BasePath
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
	if data == nil {
		return ""
	}

	// Joins two strings of url-encoded data together
	joinEncData := func(a string, b string) string {
		if a == "" {
			return b
		} else if b == "" {
			return a
		}
		return a + "&" + b
	}

	// Assumes data is a struct
	// Iterate through the fields of the struct and url-encodes the field/value pairs
	encData := ""
	rData := reflect.ValueOf(data)
	rDataType := reflect.TypeOf(data)
	urlVals := url.Values{}
	for i := 0; i < rData.NumField(); i++ {

		// Only encode non-nil values
		if rData.Field(i).Kind() != reflect.Ptr || !rData.Field(i).IsNil() {
			if rData.Field(i).Kind() == reflect.Struct {

				// If the value is a struct, recursively flatten the request
				encData = joinEncData(encData, urlEncodeData(rData.Field(i).Interface()))

			} else if tag := rDataType.Field(i).Tag.Get("json"); tag != "" {
				// Otherwise, get the json tag for this value...
				var rDataField = rData.Field(i)

				// Dereference pointer values
				if rDataField.Kind() == reflect.Ptr {
					rDataField = rDataField.Elem()
				}

				if rDataField.Kind() == reflect.Map {
					// If the value is a map, create a special key
					iter := rDataField.MapRange()
					for iter.Next() {
						k := iter.Key().Interface()
						v := iter.Value().Interface()
						urlVals.Add(fmt.Sprintf("%v[%v]", tag, k), fmt.Sprintf("%v", v))
					}
				} else {
					// Otherwise, just cast the value to a string
					urlVals.Add(tag, fmt.Sprintf("%v", rDataField.Interface()))
				}
			}
		}
	}

	// Encode the collected key/vale pairs
	encData = joinEncData(encData, urlVals.Encode())
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
