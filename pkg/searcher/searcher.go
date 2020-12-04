package searcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func worker(url string, k string, wg *sync.WaitGroup, storage *MutexMap, s chan int) {
	//defer wg.Done()
	defer func() {
		<-s
		wg.Done()
	}()
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

func Search(k string, urls []string, limit int) map[string]int {
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
		go worker(url, k, &wg, storage, semaphore)
	}
	wg.Wait()
	return storage.storage
}
