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

	// create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		storage.SetValue(url, 0, err)
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
		storage.SetValue(url, 0, err)
		return
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		storage.SetValue(url, 0, err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	// search
	count := bytes.Count(body, []byte(k))
	// write to shared map
	storage.SetValue(url, count, nil)
}

func Search(ctx context.Context, k string, urls []string, limit, timeout int) map[string]UrlResult {
	initStorage := make(map[string]UrlResult, len(urls))
	storage := NewStorage(initStorage)
	var wg sync.WaitGroup
	// no-limits case
	if limit < 1 {
		limit = len(urls)
	}
	// remove duplicating items
	urls = sliceUnique(urls)
	semaphore := make(chan struct{}, limit)
	for _, url := range urls {
		semaphore <- struct{}{}
		wg.Add(1)
		go worker(ctx, url, k, &wg, storage, semaphore, timeout)
	}
	wg.Wait()
	return storage.storage
}
