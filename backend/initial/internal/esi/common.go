package esi

import (
	"context"
	"evelp/config/global"
	"net/http"
	"sync"

	"golang.org/x/sync/semaphore"
)

var (
	limit  int64     = 3
	weigth int64     = 1
	langs  [6]string = [6]string{global.DE, global.EN, global.FR, global.JA, global.RU, global.ZH}
)

var sem = semaphore.NewWeighted(limit)
var wg sync.WaitGroup
var channel = make(chan offersWrapper)

var client = &http.Client{}

func acquireSem(weight int64) error {
	if err := sem.Acquire(context.Background(), weigth); err != nil {
		return err
	}
	return nil
}
