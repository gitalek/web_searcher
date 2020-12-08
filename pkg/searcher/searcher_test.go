package searcher

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSearchRealRequest(t *testing.T) {
	urls := []string{
		"http://example.com",
		"https://opennet.ru",
		"https://habr.com/ru/companies/",
		"https://example.com",
		"http://habr.com/ru/companies/",
		"http://opennet.ru",
		"https://go.dev/",
		"http://go.dev/",
		"http://go.dev/",
		"http://go.dev/",
	}
	keyword := "background"
	timeout := 2000
	var limit int

	// context
	ctx := context.Background()

	want := map[string]int{
		"http://example.com":             2,
		"http://go.dev/":                 0,
		"http://habr.com/ru/companies/":  1,
		"http://opennet.ru":              30,
		"https://example.com":            2,
		"https://go.dev/":                0,
		"https://habr.com/ru/companies/": 1,
		"https://opennet.ru":             30,
	}

	got := Search(ctx, keyword, urls, limit, timeout)

	if !cmp.Equal(want, got) {
		t.Errorf("Wrong object received, got=%s", cmp.Diff(want, got))
	}
}
