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

func worker(url string, k string, wg *sync.WaitGroup, storage *MutexMap, s chan int, t int) {
	//defer wg.Done()
	defer func() {
		<-s
		wg.Done()
	}()
	// do nothing, if key exists
	if _, ok := storage.GetValue(url); ok {
		return
	}

	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ctx context.Context
	// no-timeout case
	if t < 1 {
		ctx = context.Background()
	} else { // timeout case
		// create timeout-context and add it to request
		cancelCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(t))
		ctx = cancelCtx
		defer cancel()
	}
	req = req.WithContext(ctx)
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
		return
	}

	// search
	count := bytes.Count(body, []byte(k))
	// write to shared map
	storage.SetValue(url, count)
}

func Search(k string, urls []string, limit int, timeout int) map[string]int {
	initStorage := make(map[string]int, len(urls))
	storage := NewStorage(initStorage)
	var wg sync.WaitGroup
	// no-limits case
	if limit < 1 {
		limit = len(urls)
	}
	semaphore := make(chan int, limit)
	for _, url := range urls {
		semaphore <- 1
		wg.Add(1)
		go worker(url, k, &wg, storage, semaphore, timeout)
	}
	wg.Wait()
	return storage.storage
}
