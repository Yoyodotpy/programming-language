package main

// -------- VALUES ----------

type value interface {
	is_value()
}

type hex_val struct {
	val int64
}

type closure_val struct {
	param string
	body  node
	env   *env
}

type concell_val struct {
	val1 value
	val2 value
}

type prim_val struct {
	fn func(value) value
}

func (v hex_val) is_value()     {}
func (v closure_val) is_value() {}
func (v concell_val) is_value() {}
func (v prim_val) is_value()    {}

// -------- ENVIRONMENT STUFF --------

type env struct {
	parent *env
	val    map[string]value
}

func (e *env) getvar(label string) value {
	environment := e
	var val value
	var found = false
	for !found {
		val, found = environment.val[label]
		if !found {
			if environment.parent == nil {
				panic("Variable does not exist: " + label)
			}
			environment = environment.parent
		}
	}
	return val
}

func (e *env) setvar(label string, val value) {
	e.val[label] = val
}

// ------- Evaluator --------

func (e *env) eval(node node) value {
	switch n := node.(type) {
	case hex_node:
		return hex_val{val: n.val}
	case var_node:
		return e.getvar(n.label)
	case lambda_node:
		return closure_val{env: e, param: n.param, body: n.body}
	case concell_node:
		return n.concell_eval(e)
	case define_node:
		return n.define_eval(e)
	case apply_node:
		return n.apply_eval(e)
	default:
		panic("Unknown AST node type in eval")
	}
}

func (n concell_node) concell_eval(e *env) value {
	return concell_val{
		val1: e.eval(n.val1),
		val2: e.eval(n.val2),
	}
}

func (n define_node) define_eval(e *env) value {
	val := e.eval(n.value)
	e.setvar(n.label, val)
	return val
}

func (n apply_node) apply_eval(e *env) value {
	function := e.eval(n.function)
	arg := e.eval(n.arg)
	switch v := function.(type) {
	case prim_val:
		return v.fn(arg)
	case closure_val:
		envi := &env{
			parent: v.env,
			val:    make(map[string]value),
		}
		envi.setvar(v.param, arg)
		return envi.eval(v.body)
	default:
		panic("Attempted to run a non-function")
	}
}
