package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s literal\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Println(regexp.QuoteMeta(os.Args[1]))
}
