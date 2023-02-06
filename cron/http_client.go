package cron

import (
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Get(url string) ([]byte, error)
}

type CronHttpClient struct{}

func (c *CronHttpClient) Get(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
