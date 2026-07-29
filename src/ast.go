package main

import "slices"

// NODE TYPES
type node interface {
	is_node()
}

type hex_node struct {
	val string
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
		if p.pos+2 < len(p.tokens) {
			switch p.tokens[p.pos+2] {
			case ":":
				return p.parse_lambda()
			case "=":
				return p.parse_define()
			}
		}
		return p.parse_apply()
	}
	if tok == "'" {
		return p.parse_hex()
	}

	p.advance()
	return var_node{label: tok}
}

func (p *parser) parse_lambda() lambda_node {
	var node lambda_node
	p.advance()
	node.param = p.advance()
	p.advance()
	node.body = p.parse_node()
	p.advance()
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
	var node apply_node
	p.advance()
	node.function = p.parse_node()
	node.arg = p.parse_node()
	p.advance()
	return node
}

func (p *parser) parse_hex() hex_node {
	var node hex_node
	p.advance()
	node.val = p.advance()
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
