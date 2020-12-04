package searcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func worker(url string, k string, wg *sync.WaitGroup, storage *MutexMap) {
	defer wg.Done()
	// do nothing, if key exists
	if _, ok := storage.GetValue(url); ok {
		return
	}

	// make request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
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

func Search(k string, urls []string) map[string]int {
	initStorage := make(map[string]int, len(urls))
	storage := NewStorage(initStorage)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go worker(url, k, &wg, storage)
	}
	wg.Wait()
	return storage.storage
}
