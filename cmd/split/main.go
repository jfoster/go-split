package main

import (
	"flag"
	"log"

	"github.com/jfoster/go-split"
)

var flags struct {
	size    uint64
	unsplit bool
}

func main() {
	flag.Uint64Var(&flags.size, "size", 1000000, "size of part files")
	flag.BoolVar(&flags.unsplit, "unsplit", false, "split or unsplit")
	flag.Parse()

	path := flag.Args()[0]

	if flags.unsplit {
		err := split.Unsplit(path)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := split.Split(path, flags.size)
		if err != nil {
			log.Fatal(err)
		}
	}
}
