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
		fmt.Fprintf(os.Stderr, "usage: %s [option...] regexp [file...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	lines := flag.Bool("lines", false, "print only matching lines")
	paths := flag.Bool("paths", false, "print only matching paths")
	quiet := flag.Bool("quiet", false, "exit after finding a match")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	exp, err := regexp.Compile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error compiling regexp: %v\n", os.Args[0], err)
		os.Exit(1)
	}

	var inputs []io.Reader
	if len(args) == 1 {
		inputs = []io.Reader{os.Stdin}
	} else {
		inputs = make([]io.Reader, len(args)-1)
		for i, path := range args[1:] {
			inputs[i], err = os.Open(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], path, err)
				os.Exit(1)
			}
		}
	}

	*lines = *lines || len(inputs) == 1
	matchFound := false
	for i, input := range inputs {
		r := bufio.NewReader(input)
		for line, err := r.ReadBytes('\n'); err == nil || err == io.EOF; line, err = r.ReadBytes('\n') {
			if exp.Match(line) {
				matchFound = true
				if *quiet {
					os.Exit(0)
				} else if *paths {
					fmt.Println(args[i+1])
					break
				} else if *lines {
					fmt.Printf("%s", line)
				} else {
					fmt.Printf("%s:%s", args[i+1], line)
				}
			}

			if err == io.EOF {
				break
			}
		}
	}

	if matchFound {
		os.Exit(0)
	}
	os.Exit(1)
}
