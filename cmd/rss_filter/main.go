package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	err := run(flag.Arg(0))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run(fname string) error {
	return nil
}
