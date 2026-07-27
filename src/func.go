package main

import (
	"fmt"
	"io"
	"os"
)

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func read_file(path string) string {
	file, err := os.Open(path)
	checkerr(err)
	defer file.Close()

	text, err := io.ReadAll(file)
	checkerr(err)

	return string(text)
}
