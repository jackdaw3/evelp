package net

import (
	"evelp/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	request string = "http://127.0.0.1:9090/hello"
	hello   string = "hello"
	mu      sync.Mutex
)

func setUpServer() {
	log.Init()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		defer mu.Unlock()
		mu.Lock()
		fmt.Fprintln(w, hello)
	})
	http.ListenAndServe(":9090", nil)
}

func TestGetWithRetries(t *testing.T) {
	go setUpServer()

	resp, err := Get(request)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, hello, strings.ReplaceAll(string(body), "\n", ""))
	assert.NoError(t, err)
}
