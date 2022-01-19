package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var delim = '\n'

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s regexp [regexp...] [: path [path...]]\n", os.Args[0])
		flag.PrintDefaults()
	}
	dirs := flag.Bool("dirs", false, "only check directories")
	files := flag.Bool("files", false, "only check files")
	union := flag.Bool("union", false, "take the union of the provided expressions instead of the intersection")
	null := flag.Bool("0", false, "use a null character as a record delimiter (instead of newline)")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *null {
		delim = '\000'
	}

	var exps []*regexp.Regexp
	var paths []string
	foundDelim := false
	for _, arg := range args {
		if !foundDelim && arg == ":" {
			foundDelim = true
		} else if foundDelim {
			paths = append(paths, arg)
		} else {
			exp, err := regexp.Compile(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: error compiling regexp: %v\n", os.Args[0], err)
				os.Exit(1)
			}
			exps = append(exps, exp)
		}
	}

	if len(exps) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, p := range paths {
		filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			}

			if *files && d.IsDir() {
				return nil
			} else if *dirs && !d.IsDir() {
				return nil
			}

			if *union {
				for _, exp := range exps {
					if exp.MatchString(d.Name()) {
						fmt.Println(path)
						break
					}
				}
			} else {
				match := true
				for _, exp := range exps {
					match = match && exp.MatchString(d.Name())
					if !match {
						break
					}
				}
				if match {
					fmt.Printf("%s%c", path, delim)
				}
			}

			return nil
		})
	}
}
