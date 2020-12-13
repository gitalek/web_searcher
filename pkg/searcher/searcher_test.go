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

func TestSearchWithTestServer(t *testing.T) {
	sites := []string{
		"http_example_com.html",
		"http_go_dev.html",
		"http_habr.html",
		"http_opennet_ru.html",
		"https_example_com.html",
		"https_go_dev.html",
		"https_habr.html",
		"https_opennet_ru.html",
	}
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
		fmt.Sprintf("%s/http_example_com.html/", ts.URL):  &(UrlResult{count: 2}),
		fmt.Sprintf("%s/http_go_dev.html/", ts.URL):       &(UrlResult{count: 0}),
		fmt.Sprintf("%s/http_habr.html/", ts.URL):         &(UrlResult{count: 1}),
		fmt.Sprintf("%s/http_opennet_ru.html/", ts.URL):   &(UrlResult{count: 31}),
		fmt.Sprintf("%s/https_example_com.html/", ts.URL): &(UrlResult{count: 2}),
		fmt.Sprintf("%s/https_go_dev.html/", ts.URL):      &(UrlResult{count: 0}),
		fmt.Sprintf("%s/https_habr.html/", ts.URL):        &(UrlResult{count: 1}),
		fmt.Sprintf("%s/https_opennet_ru.html/", ts.URL):  &(UrlResult{count: 31}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf(
			"Wrong object received:\ngot:\t%v\nwant:\t%v",
			got, want,
		)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf(
				"Got object has no \"%s\" site\nWrong object received:\ngot:\t%#v\nwant:\t%#v\n",
				site, got, want,
			)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf(
				"Objects don't match at site %s: got -> %d, want -> %d\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.count, wantRes.count, got, want,
			)
			break
		}
		if gotRes.err != nil {
			t.Errorf(
				"Objects ERRORS don't match at site %s: got -> %s, should be -> %v\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err.Error(), nil, got, want,
			)
			break
		}
	}
}

func TestSearchNewRequestError(t *testing.T) {
	sites := []string{
		"http_example_com.html",
	}
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
		fmt.Sprintf("%s/http_example_com.html/", ts.URL): &(UrlResult{err: errors.New("net/http: nil Context")}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf(
			"Wrong object received:\ngot:\t%v\nwant:\t%v",
			got, want,
		)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf(
				"Got object has no \"%s\" site\nWrong object received:\ngot:\t%#v\nwant:\t%#v\n",
				site, got, want,
			)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf(
				"Objects don't match at site %s: got -> %d, want -> %d\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.count, wantRes.count, got, want,
			)
			break
		}
		if wantRes.err == nil {
			if gotRes.err == nil {
				continue
			}
			t.Errorf(
				"Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %v\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err, nil, got, want,
			)
			break
		}
		if wantRes.err != nil {
			if gotRes.err == nil {
				t.Errorf(
					"Objects ERRORS don't match at site %s: got nil -> %s, should be -> %s\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
					site, gotRes.err, wantRes.err, got, want,
				)
				break
			}
			if gotRes.err.Error() == wantRes.err.Error() {
				continue
			}
			t.Errorf(
				"Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %s\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err.Error(), wantRes.err.Error(), got, want,
			)
			break
		}
	}
}

func TestSearchDoRequestError(t *testing.T) {
	sites := []string{
		"https_example_com.html",
		"https_go_dev.html",
		"https_habr.html",
		"https_opennet_ru.html",
	}
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
		helper(ts.URL, "/https_example_com.html/"): &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_example_com.html/"))}),
		helper(ts.URL, "/https_go_dev.html/"):      &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_go_dev.html/"))}),
		helper(ts.URL, "/https_habr.html/"):        &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_habr.html/"))}),
		helper(ts.URL, "/https_opennet_ru.html/"):  &(UrlResult{count: 0, err: fmt.Errorf("Get \"%s\": context deadline exceeded", helper(ts.URL, "/https_opennet_ru.html/"))}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf(
			"Wrong object received:\ngot:\t%v\nwant:\t%v",
			got, want,
		)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf(
				"Got object has no \"%s\" site\nWrong object received:\ngot:\t%#v\nwant:\t%#v",
				site, got, want,
			)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf(
				"Objects don't match at site %s: got -> %d, want -> %d\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.count, wantRes.count, got, want,
			)
			break
		}
		if wantRes.err == nil {
			if gotRes.err == nil {
				continue
			}
			t.Errorf(
				"Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %v\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err, nil, got, want,
			)
			break
		}
		if wantRes.err != nil {
			if gotRes.err == nil {
				t.Errorf(
					"Objects ERRORS don't match at site %s: got nil -> %s, should be -> %s\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
					site, gotRes.err, wantRes.err, got, want,
				)
				break
			}
			if gotRes.err.Error() == wantRes.err.Error() {
				continue
			}
			t.Errorf(
				"Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %s\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err.Error(), wantRes.err.Error(), got, want,
			)
			break
		}
	}
}

func helper(url string, suffix string) string {
	return fmt.Sprintf("%s%s", url, suffix)
}

func TestSearchReadBodyError(t *testing.T) {
	sites := []string{
		"http_example_com.html",
	}
	var urls []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer ts.Close()
	for _, filename := range sites {
		urls = append(urls, fmt.Sprintf("%s/%s/", ts.URL, filename))
	}
	keyword := "background"
	// timeout case
	timeout := 2000
	var limit int
	// context
	ctx := context.Background()

	wantErr := errors.New("unexpected EOF")
	want := map[string]*UrlResult{
		helper(ts.URL, "/http_example_com.html/"): &(UrlResult{count: 0, err: wantErr}),
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if len(got) != len(want) {
		t.Fatalf(
			"Wrong object received:\ngot:\t%v\nwant:\t%v",
			got, want,
		)
	}
	for site, wantRes := range want {
		gotRes, ok := got[site]
		if !ok {
			t.Errorf(
				"Got object has no \"%s\" site\nWrong object received:\ngot:\t%#v\nwant:\t%#v",
				site, got, want,
			)
			break
		}
		if gotRes.count != wantRes.count {
			t.Errorf(
				"Objects don't match at site %s: got -> %d, want -> %d\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.count, wantRes.count, got, want,
			)
			break
		}
		if wantRes.err == nil {
			if gotRes.err == nil {
				continue
			}
			t.Errorf(
				"Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %v\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err, nil, got, want,
			)
			break
		}
		if wantRes.err != nil {
			if gotRes.err == nil {
				t.Errorf(
					"Objects ERRORS don't match at site %s: got nil -> %s, should be -> %s\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
					site, gotRes.err, wantRes.err, got, want,
				)
				break
			}
			if gotRes.err.Error() == wantRes.err.Error() {
				continue
			}
			t.Errorf(
				"Objects ERRORS don't match at site %s: got ERROR -> %s, should be -> %s\nWrong object received:\ngot:\t%v\nwant:\t%v\n",
				site, gotRes.err.Error(), wantRes.err.Error(), got, want,
			)
			break
		}
	}
}
