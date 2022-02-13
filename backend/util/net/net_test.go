package net

import (
	"evelp/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	request string = "http://127.0.0.1:9090/hello"
	hello   string = "hello"
	count   int    = 0
	mu      sync.Mutex
)

func setUpServer() {
	log.Init()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		defer mu.Unlock()
		mu.Lock()
		if count == 0 {
			time.Sleep(1500 * time.Millisecond)
		}
		count++
		fmt.Fprintln(w, hello)
	})
	http.ListenAndServe(":9090", nil)
}

func TestGetWithRetries(t *testing.T) {
	go setUpServer()

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := GetWithRetries(client, request)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, hello, strings.ReplaceAll(string(body), "\n", ""))
	assert.NoError(t, err)
}
