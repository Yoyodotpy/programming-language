package main

import (
	"fmt"
	"os"
)

var file_text_unprocessed = ""

func main() {
	file_text_unprocessed = read_file(os.Args[1])
	fmt.Println(file_text_unprocessed)
}
