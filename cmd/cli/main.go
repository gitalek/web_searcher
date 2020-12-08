package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gitalek/web_searcher/pkg/searcher"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type urlsFlag struct {
	Urls []string
}

func (f *urlsFlag) GetUrls() []string {
	return f.Urls
}

func (f *urlsFlag) String() string {
	return fmt.Sprint(f.Urls)
}

func (f *urlsFlag) Set(v string) error {
	if len(f.Urls) > 0 {
		return errors.New("urls flag has already been set")
	}
	urls := strings.Split(v, ",")
	for _, item := range urls {
		f.Urls = append(f.Urls, item)
	}
	return nil
}

func main() {
	// setup flags
	keyword := flag.String("k", "", "Search string")
	limit := flag.Int("l", 0, "goroutine limitation")
	timeout := flag.Int("t", 0, "request timeout in milliseconds")
	var urls urlsFlag
	flag.Var(&urls, "urls", "Comma-separated urls list")
	flag.Parse()

	// context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
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

	// make search
	results := searcher.Search(ctx, *keyword, urls.GetUrls(), *limit, *timeout)

	// print results
	fmt.Printf("keyword: %#v\n", *keyword)
	fmt.Printf("%#v\n", results)
	for url, count := range results {
		fmt.Printf("at %s -> %d\n", url, count)
	}
}
