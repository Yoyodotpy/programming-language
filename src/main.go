package main

import (
	"fmt"
	"os"
)

func main() {
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

func print_ast(n node) string {
	switch v := n.(type) {
	case var_node:
		return v.label
	case hex_node:
		return fmt.Sprintf("'%x'", v.val)
	case lambda_node:
		return fmt.Sprintf("(%s : %s)", v.param, print_ast(v.body))
	case apply_node:
		return fmt.Sprintf("(%s %s)", print_ast(v.function), print_ast(v.arg))
	case concell_node:
		return fmt.Sprintf("[%s %s]", print_ast(v.val1), print_ast(v.val2))
	case define_node:
		return fmt.Sprintf("(%s = %s)", v.label, print_ast(v.value))
	default:
		return "UNKNOWN_NODE"
	}
}
