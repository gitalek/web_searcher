package searcher

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Search(k string, urls []string) (map[string]int, error) {
	results := make(map[string]int)

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return results, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return results, err
		}

		err = resp.Body.Close()
		if err != nil {
			return results, err
		}

		count := bytes.Count(body, []byte(k))
		if _, ok := results[url]; !ok {
			results[url] = count
		}
	}

	return results, nil
}
