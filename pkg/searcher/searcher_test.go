package searcher

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
			_, err := res.Write(content)
			if err != nil {
				return
			}
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

	want := map[string]*UrlResult{
		fmt.Sprintf("%s/http_example_com.html/", ts.URL):  &(UrlResult{count: 2, err: nil}),
		fmt.Sprintf("%s/http_go_dev.html/", ts.URL):       &(UrlResult{count: 0, err: nil}),
		fmt.Sprintf("%s/http_habr.html/", ts.URL):         &(UrlResult{count: 1, err: nil}),
		fmt.Sprintf("%s/http_opennet_ru.html/", ts.URL):   &(UrlResult{count: 31, err: nil}),
		fmt.Sprintf("%s/https_example_com.html/", ts.URL): &(UrlResult{count: 2, err: nil}),
		fmt.Sprintf("%s/https_go_dev.html/", ts.URL):      &(UrlResult{count: 0, err: nil}),
		fmt.Sprintf("%s/https_habr.html/", ts.URL):        &(UrlResult{count: 1, err: nil}),
		fmt.Sprintf("%s/https_opennet_ru.html/", ts.URL):  &(UrlResult{count: 31, err: nil}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf("Wrong object received:\ngot:\t%v\nwant:\t%v", got, want)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf("Got object has no \"%s\" site\n", site)
			t.Errorf("Wrong object received:\ngot:\t%#v\nwant:\t%#v", got, want)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf("Objects don't match at site %s: got -> %d, want -> %d\n", site, gotRes.count, wantRes.count)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
		if gotRes.err != nil {
			t.Errorf("Objects ERRORS don't match at site %s: got -> %s, should be -> %v\n", site, gotRes.err.Error(), nil)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
	}
}

func TestSearchNewRequestError(t *testing.T) {
	var urls []string
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	for _, filename := range sites {
		urls = append(urls, fmt.Sprintf("%s/%s/", ts.URL, filename))
	}
	keyword := "background"
	timeout := 2000
	var limit int

	// context
	// nil context case
	var ctx context.Context

	want := map[string]*UrlResult{
		fmt.Sprintf("%s/http_example_com.html/", ts.URL):  &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/http_go_dev.html/", ts.URL):       &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/http_habr.html/", ts.URL):         &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/http_opennet_ru.html/", ts.URL):   &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/https_example_com.html/", ts.URL): &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/https_go_dev.html/", ts.URL):      &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/https_habr.html/", ts.URL):        &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
		fmt.Sprintf("%s/https_opennet_ru.html/", ts.URL):  &(UrlResult{count: 0, err: errors.New("net/http: nil Context")}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf("Wrong object received:\ngot:\t%v\nwant:\t%v", got, want)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf("Got object has no \"%s\" site\n", site)
			t.Errorf("Wrong object received:\ngot:\t%#v\nwant:\t%#v", got, want)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf("Objects don't match at site %s: got -> %d, want -> %d\n", site, gotRes.count, wantRes.count)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
		if wantRes.err == nil {
			if gotRes.err == nil {
				continue
			}
			t.Errorf("Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %v\n", site, gotRes.err, nil)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
		if wantRes.err != nil {
			if gotRes.err == nil {
				t.Errorf("Objects ERRORS don't match at site %s: got nil -> %s, should be -> %s\n", site, gotRes.err, wantRes.err)
				t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
				break
			}
			if gotRes.err.Error() == wantRes.err.Error() {
				continue
			}
			t.Errorf("Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %s\n", site, gotRes.err.Error(), wantRes.err.Error())
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
	}
}

func TestSearchDoRequestError(t *testing.T) {
	var urls []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// timeout case
		time.Sleep(2 * time.Millisecond)
		fmt.Fprintln(w, "site content")
	}))
	defer ts.Close()
	for _, filename := range sites {
		urls = append(urls, fmt.Sprintf("%s/%s/", ts.URL, filename))
	}
	keyword := "background"
	// timeout case
	timeout := 1
	var limit int
	// context
	ctx := context.Background()

	want := map[string]*UrlResult{
		helper(ts.URL, "/http_example_com.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/http_example_com.html/"))}),
		helper(ts.URL, "/http_go_dev.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/http_go_dev.html/"))}),
		helper(ts.URL, "/http_habr.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/http_habr.html/"))}),
		helper(ts.URL, "/http_opennet_ru.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/http_opennet_ru.html/"))}),
		helper(ts.URL, "/https_example_com.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_example_com.html/"))}),
		helper(ts.URL, "/https_go_dev.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_go_dev.html/"))}),
		helper(ts.URL, "/https_habr.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_habr.html/"))}),
		helper(ts.URL, "/https_opennet_ru.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_opennet_ru.html/"))}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf("Wrong object received:\ngot:\t%v\nwant:\t%v", got, want)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf("Got object has no \"%s\" site\n", site)
			t.Errorf("Wrong object received:\ngot:\t%#v\nwant:\t%#v", got, want)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf("Objects don't match at site %s: got -> %d, want -> %d\n", site, gotRes.count, wantRes.count)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
		if wantRes.err == nil {
			if gotRes.err == nil {
				continue
			}
			t.Errorf("Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %v\n", site, gotRes.err, nil)
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
		if wantRes.err != nil {
			if gotRes.err == nil {
				t.Errorf("Objects ERRORS don't match at site %s: got nil -> %s, should be -> %s\n", site, gotRes.err, wantRes.err)
				t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
				break
			}
			if gotRes.err.Error() == wantRes.err.Error() {
				continue
			}
			t.Errorf("Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %s\n", site, gotRes.err.Error(), wantRes.err.Error())
			t.Errorf("Wrong object received:\ngot:\t%v\nwant:\t%v\n", got, want)
			break
		}
	}

}

func helper(url string, suffix string) string {
	return fmt.Sprintf("%s%s", url, suffix)
}
