package searcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Search(k string, urls []string) map[string]int {
	results := make(map[string]int, len(urls))

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		err = resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}

		count := bytes.Count(body, []byte(k))
		if _, ok := results[url]; !ok {
			results[url] = count
		}
	}

	return results
}
