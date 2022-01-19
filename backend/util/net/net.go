package net

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
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
			} else {
				err = fmt.Errorf("request %s error status code %d", request, code)
			}
		}

		log.Warn(err)
		log.Warnf("request %s retrying in %v", request, backoff)
		time.Sleep(backoff)
	}

	if err != nil {
		return nil, fmt.Errorf("all request retries failed: %v", err)
	}

	return resp, nil
}

func Get(client *http.Client, request string) (*http.Response, error) {
	resp, err := client.Get(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
