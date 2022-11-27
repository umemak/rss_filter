package main

import (
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

func run(name string) error {
	res, err := rssfilter.GetByName(name)
	if err != nil {
		return fmt.Errorf("rssfilter.GetByName: %w", err)
	}
	fmt.Println(res)
	return nil
}
