package esi

import (
	"net/http"
	"sync"
)

var (
	wg     sync.WaitGroup
	client = &http.Client{}
)
