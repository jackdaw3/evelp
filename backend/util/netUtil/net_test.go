package netUtil

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	request string = "http://127.0.0.1:9090/hello"
	name    string = "hello"
)

func setUpServer() {
	count := 0
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			time.Sleep(1500 * time.Millisecond)
		}
		count++
		fmt.Fprintln(w, name)

	})
	http.ListenAndServe(":9090", nil)
}

func TestGetWithRetries(t *testing.T) {
	go setUpServer()

	client := &http.Client{Timeout: 1 * time.Second}

	body, err := GetWithRetries(client, request)
	str := strings.ReplaceAll(string(body), "\n", "")
	assert.Equal(t, name, str)
	assert.Nil(t, err)
}
