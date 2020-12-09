package searcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var sites = []string{
	"http_example_com.html",
	"http_go_dev.html",
	"http_habr.html",
	"http_opennet_ru.html",
	"https_example_com.html",
	"https_go_dev.html",
	"https_habr.html",
	"https_opennet_ru.html",
}

func TestSearchWithTestServer(t *testing.T) {
	var urls []string
	mux := http.NewServeMux()
	for _, filename := range sites {
		content, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", filename))
		if err != nil {
			t.Fatal(err)
		}
		mux.HandleFunc(fmt.Sprintf("/%s/", filename), func(res http.ResponseWriter, req *http.Request) {
			res.Write(content)
		})
	}
	ts := httptest.NewServer(mux)
	defer ts.Close()
	for _, filename := range sites {
		urls = append(urls, fmt.Sprintf("%s/%s/", ts.URL, filename))
	}
	keyword := "background"
	timeout := 2000
	var limit int

	// context
	ctx := context.Background()

	want := map[string]int{
		fmt.Sprintf("%s/http_example_com.html/", ts.URL):  2,
		fmt.Sprintf("%s/http_go_dev.html/", ts.URL):       0,
		fmt.Sprintf("%s/http_habr.html/", ts.URL):         1,
		fmt.Sprintf("%s/http_opennet_ru.html/", ts.URL):   31,
		fmt.Sprintf("%s/https_example_com.html/", ts.URL): 2,
		fmt.Sprintf("%s/https_go_dev.html/", ts.URL):      0,
		fmt.Sprintf("%s/https_habr.html/", ts.URL):        1,
		fmt.Sprintf("%s/https_opennet_ru.html/", ts.URL):  31,
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf("Wrong object received:\ngot:\t%v\nwant:\t%v", got, want)
	}
	for site, wantNum := range want {
		gotNum, ok := got[site]
		if !ok {
			t.Errorf("Got object has no \"%s\" site\n", site)
			t.Errorf("Wrong object received:\ngot:\t%#v\nwant:\t%#v", got, want)
			break
		}
		if gotNum != wantNum {
			t.Errorf("Objects don't match at site %s: got -> %d, want -> %d\n", site, gotNum, wantNum)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
	}
}
