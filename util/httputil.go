package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	client *http.Client
)

func NewHttpClient() *http.Client {
	client = &http.Client{}
	return client
}

func PerformGetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func PerformPostRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func PerformPutRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func PerformDeleteRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func PerformPutRequestWithQueryParams(url string, body map[string]string, queryParams map[string]string) (*http.Response, error) {

	parsedBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(parsedBody))
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	return client.Do(req)
}
