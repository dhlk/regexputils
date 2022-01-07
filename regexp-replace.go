package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [option...] regexp replacement [file...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	exp, err := regexp.Compile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error compiling regexp: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	replacement := []byte(args[1])

	var inputs []io.Reader
	if len(args) == 2 {
		inputs = []io.Reader{os.Stdin}
	} else {
		inputs = make([]io.Reader, len(args)-2)
		for i, path := range args[2:] {
			inputs[i], err = os.Open(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], path, err)
				os.Exit(1)
			}
		}
	}

	for _, input := range inputs {
		r := bufio.NewReader(input)
		for line, err := r.ReadBytes('\n'); err == nil || err == io.EOF; line, err = r.ReadBytes('\n') {
			fmt.Printf("%s", exp.ReplaceAll(line, replacement))
			if err == io.EOF {
				break
			}
		}
	}

	os.Exit(0)
}
