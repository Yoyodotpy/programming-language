package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
	}

	file_text_unprocessed := read_file(os.Args[1])
	lexed := lexer(file_text_unprocessed)
	p := parser{
		tokens: lexed,
		pos:    0,
	}
	ast := p.parse()

	global_env := &env{
		parent: nil,
		val:    make(map[string]value),
	}

	init_primitives(global_env)

	for _, n := range ast {
		global_env.eval(n)
	}
}

func help() {
	fmt.Println("You need to input a program file for this to work!")
	fmt.Println("You can find some examples in the git repo:")
	fmt.Println("https://github.com/Yoyodotpy/programming-language/tree/main/examples")
	fmt.Println("\nCommand format:")
	fmt.Printf("gfpl /path/to/file\n\n")
	os.Exit(1)
}
