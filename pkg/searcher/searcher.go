package searcher

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func worker(ctx context.Context, url, k string, wg *sync.WaitGroup, storage *MutexMap, s chan struct{}, t int) {
	defer func() {
		<-s
		wg.Done()
	}()
	// do nothing, if key exists
	if _, ok := storage.GetValue(url); ok {
		return
	}

	// create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if t > 0 { // timeout case
		// create timeout-context and add it to request
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(t))
		defer cancel()
		req = req.WithContext(ctx)
	}
	// create client and run request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	// search
	count := bytes.Count(body, []byte(k))
	// write to shared map
	storage.SetValue(url, count)
}

func Search(ctx context.Context, k string, urls []string, limit, timeout int) map[string]int {
	initStorage := make(map[string]int, len(urls))
	storage := NewStorage(initStorage)
	var wg sync.WaitGroup
	// no-limits case
	if limit < 1 {
		limit = len(urls)
	}
	semaphore := make(chan struct{}, limit)
	for _, url := range urls {
		semaphore <- struct{}{}
		wg.Add(1)
		go worker(ctx, url, k, &wg, storage, semaphore, timeout)
	}
	wg.Wait()
	return storage.storage
}
