package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gitalek/web_searcher/pkg/searcher"
	"strings"
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
	var urls urlsFlag
	flag.Var(&urls, "urls", "Comma-separated urls list")
	flag.Parse()

	// make search
	results := searcher.Search(*keyword, urls.GetUrls())
	// print results
	fmt.Printf("keyword: %#v\n", *keyword)
	for url, count := range results {
		fmt.Printf("at %s -> %d\n", url, count)
	}
}
