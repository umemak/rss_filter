package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	rssfilter "github.com/umemak/rss_filter"
)

func main() {
	flag.Parse()
	err := run(flag.Arg(0))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run(url string) error {
	buf, err := rssfilter.Fetch(url)
	if err != nil {
		return fmt.Errorf("rssfilter.Fetch: %w", err)
	}
	rss, err := rssfilter.Parse(buf)
	if err != nil {
		return fmt.Errorf("rssfilter.Parse: %w", err)
	}
	res, err := xml.MarshalIndent(rss, "", "  ")
	fmt.Println(string(res))
	return nil
}
