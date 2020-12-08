package searcher

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"os"
	"os/signal"
	"syscall"
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
	ctx, cancel := context.WithCancel(ctx)

	want := map[string]int{
		"http://example.com": 2,
		"http://go.dev/": 0,
		"http://habr.com/ru/companies/": 1,
		"http://opennet.ru": 31,
		"https://example.com": 2,
		"https://go.dev/": 0,
		"https://habr.com/ru/companies/": 1,
		"https://opennet.ru": 31,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		<-c
		cancel()
	}()

	got := Search(ctx, keyword, urls, limit, timeout)

	if !cmp.Equal(want, got) {
		t.Errorf("Wrong object received, got=%s", cmp.Diff(want, got))
	}
}
