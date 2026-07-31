package main

import (
	"slices"
	"strconv"
)

// NODE TYPES
type node interface {
	is_node()
}

type hex_node struct {
	val int64
}

type var_node struct {
	label string
}

type lambda_node struct {
	param string
	body  node
}

type apply_node struct {
	function node
	arg      node
}

type concell_node struct {
	val1 node
	val2 node
}

type define_node struct {
	label string
	value node
}

func (n hex_node) is_node()     {}
func (n var_node) is_node()     {}
func (n lambda_node) is_node()  {}
func (n apply_node) is_node()   {}
func (n concell_node) is_node() {}
func (n define_node) is_node()  {}

// ACTUAL AST CODE

type parser struct {
	tokens []string
	pos    int
}

func (p *parser) current() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *parser) advance() string {
	tok := p.current()
	p.pos++
	return tok
}

func (p *parser) parse_node() node {
	tok := p.current()

	if tok == "[" {
		return p.parse_concell()
	}
	if tok == "(" {
		for _, curtok := range p.tokens[p.pos+1:] {
			if curtok == "(" || curtok == ")" {
				break
			}
			switch curtok {
			case "=":
				return p.parse_define()
			case ":":
				return p.parse_lambda()
			case "do":
				return p.parse_do()
			}
		}
		return p.parse_apply()
	}
	if tok == "'" {
		return p.parse_hex()
	}
	if len(tok) >= 1 && tok[0] == '"' {
		return p.parse_string()
	}
	if str_is_int(tok) {
		return p.parse_int()
	}

	p.advance()
	return var_node{label: tok}
}

func (p *parser) parse_lambda() lambda_node {
	var params []string
	var node lambda_node

	p.advance()
	for p.current() != ":" && p.current() != "" {
		tok := p.advance()
		if tok != "," {
			params = append(params, tok)
		}
	}
	p.advance()
	node.body = p.parse_node()
	p.advance()

	if len(params) >= 1 {
		slices.Reverse(params)
		node.param = params[0]
		params = params[1:]
		for i := range len(params) {
			node = lambda_node{
				param: params[i],
				body:  node,
			}
		}
	} else {
		node.param = "_"
	}

	return node
}

func (p *parser) parse_define() define_node {
	var node define_node
	p.advance()
	node.label = p.advance()
	p.advance()
	node.value = p.parse_node()
	p.advance()
	return node
}

func (p *parser) parse_apply() apply_node {
	var applynode apply_node
	var args []node

	p.advance()
	applynode.function = p.parse_node()
	for p.current() != ")" && p.current() != "" {
		if p.current() == "," {
			p.advance()
			continue
		}
		args = append(args, p.parse_node())
	}
	p.advance()

	if len(args) >= 1 {
		applynode.arg = args[0]
		args = args[1:]
		for i := range len(args) {
			applynode = apply_node{
				function: applynode,
				arg:      args[i],
			}
		}
	} else {
		applynode.arg = var_node{label: "_"}
	}

	return applynode
}

func (p *parser) parse_do() node {
	var args []node

	p.advance()
	p.advance()
	for p.current() != ")" && p.current() != "" {
		if p.current() == "," {
			p.advance()
			continue
		}
		args = append(args, p.parse_node())
	}
	p.advance()

	if len(args) < 1 {
		return var_node{label: "nil"}
	}

	slices.Reverse(args)

	node := args[0]
	args = args[1:]

	for _, cur := range args {
		def_node, is_def := cur.(define_node)
		if is_def {
			node = apply_node{
				function: lambda_node{
					param: def_node.label,
					body:  node,
				},
				arg: def_node.value,
			}
		} else {
			node = apply_node{
				function: lambda_node{
					param: "_",
					body:  node,
				},
				arg: cur,
			}
		}
	}
	return node
}

func (p *parser) parse_int() node {
	str := p.advance()
	numb, err := strconv.Atoi(str)
	checkerr(err)

	return hex_node{
		val: int64(numb),
	}
}

func (p *parser) parse_string() node {
	var list []node
	str := p.advance()
	str = str[1 : len(str)-1]

	if len(str) < 1 {
		return var_node{label: "nil"}
	}

	for _, r := range str {
		node_val := hex_node{val: int64(r)}
		list = append(list, node_val)
	}

	slices.Reverse(list)

	var node node = list[0]
	list = list[1:]

	for i := range len(list) {
		node = concell_node{
			val1: list[i],
			val2: node,
		}
	}

	return node
}

func (p *parser) parse_hex() hex_node {
	var node hex_node
	p.advance()
	var err error
	node.val, err = strconv.ParseInt(p.advance(), 16, 64)
	checkerr(err)
	p.advance()
	return node
}

func (p *parser) parse_concell() node {
	var list []node

	p.advance()
	for p.current() != "]" && p.current() != "" {
		list = append(list, p.parse_node())
	}
	p.advance()

	switch len(list) {
	case 0:
		return var_node{label: "nil"}
	case 1:
		return list[0]
	}

	slices.Reverse(list)

	var node node = list[0]
	list = list[1:]

	for i := range len(list) {
		node = concell_node{
			val1: list[i],
			val2: node,
		}
	}

	return node
}

func (p *parser) parse() []node {
	var ast []node
	for p.current() != "" {
		node := p.parse_node()
		ast = append(ast, node)
	}
	return ast
}
