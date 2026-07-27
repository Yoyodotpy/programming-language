package main

type node interface {
	is_node()
}

type hex_node struct {
	val string
}

func (n hex_node) is_node() {}

type var_node struct {
	label string
}

func (n var_node) is_node() {}

type lambda_node struct {
	params []string
	body   node
}

func (n lambda_node) is_node() {}

type apply_node struct {
	function node
	arg      node
}

func (n apply_node) is_node() {}
