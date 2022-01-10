package util

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	3 * time.Second,
}

func GetWithRetries(client *http.Client, request string) ([]byte, error) {
	var body []byte
	var err error

	for _, backoff := range backoffSchedule {
		body, err = get(client, request)

		if err == nil {
			break
		}

		log.Warnf("Request %s error: %+v\n", request, err)
		log.Warnf("Request %s retrying in %v\n", request, backoff)
		time.Sleep(backoff)
	}

	if err != nil {
		log.Errorf("Request %s error: %+v\n", request, err)
		return nil, err
	}

	return body, nil
}

func get(client *http.Client, request string) ([]byte, error) {
	resp, err := client.Get(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
