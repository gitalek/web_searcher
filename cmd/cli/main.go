package main

import (
"errors"
"flag"
"fmt"
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
	keyword := flag.String("k", "", "Search string")
	var urls urlsFlag
	flag.Var(&urls, "urls", "Comma-separated urls list")
	flag.Parse()

	fmt.Printf("Search string -> %#v\n", *keyword)
	fmt.Printf("%#v\n", urls.GetUrls())

	for i, url := range urls.GetUrls() {
		fmt.Printf("at %d: %s\n", i, url)
	}
}
