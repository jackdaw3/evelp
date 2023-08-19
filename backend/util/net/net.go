package net

import (
	"evelp/log"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type RetryStrategy struct {
	MaxRetries int
	Interval   time.Duration
	Factor     float64
}

type RetryableHttpClient struct {
	Client  *http.Client
	Retryer *RetryStrategy
}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 3,
	},
	Timeout: 100 * time.Second,
}

var retryer = &RetryStrategy{
	MaxRetries: 5,
	Interval:   time.Second,
	Factor:     3,
}

var retryableClient = &RetryableHttpClient{
	Client:  client,
	Retryer: retryer,
}

func Get(req string) (*http.Response, error) {
	retries := 0
	interval := retryableClient.Retryer.Interval

	for {
		resp, err := retryableClient.Client.Get(req)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		if err == nil && resp.StatusCode >= 500 && resp.StatusCode < 600 {
			resp.Body.Close()

			if retries >= retryableClient.Retryer.MaxRetries {
				return nil, fmt.Errorf("max retries reached: %d", retries)
			}

			log.Infof("retrying request after %v, retry count: %d\n", interval, retries)
			time.Sleep(interval)
			retries++
			interval = time.Duration(float64(interval) * retryableClient.Retryer.Factor)

			continue
		}

		return resp, err
	}
}
