package main

import (
	"log"
	"os"
)

type Input struct {
	source *os.File
	out    string
}

func get_input() *Input {
	var structured Input

	for i, arg := range os.Args {
		if i == 1 {
			f, err := os.Open(arg)
			if err != nil {
				log.Fatalln(err)
			}
			structured.source = f
		} else if i == 2 {
			structured.out = arg
		}
	}

	if structured.source == nil {
		structured.source = os.Stdin
	}
	if structured.out == "" {
		structured.out = "a.out"
	}

	return &structured
}
