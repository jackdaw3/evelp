package net

import (
	"evelp/log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	3 * time.Second,
}

func GetWithRetries(client *http.Client, request string) (*http.Response, error) {
	var resp *http.Response
	var err error

	for _, backoff := range backoffSchedule {
		resp, err = Get(client, request)

		if err == nil {
			code := resp.StatusCode
			if code == http.StatusOK {
				break
			}
			err = errors.Errorf("http request %s error status code %d", request, code)
		}

		log.Warnf("http request %s failed: %+v \nretrying in %v", request, err, backoff)
		time.Sleep(backoff)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "http all request retries failed")
	}

	return resp, nil
}

func Get(client *http.Client, request string) (*http.Response, error) {
	resp, err := client.Get(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
