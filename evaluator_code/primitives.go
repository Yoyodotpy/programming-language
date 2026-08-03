package main

import "fmt"

func (e *env) define_primitive(name string, fn func(value) value) {
	e.setvar(name, prim_val{fn: fn})
}

func init_primitives(e *env) {
	e.define_primitive("print", func(arg value) value {
		switch v := arg.(type) {
		case concell_val:
			cur := value(v)
			var list []int64

			for {
				cell, ok := cur.(concell_val)
				if !ok {
					if h, ok := cur.(hex_val); ok {
						list = append(list, h.val)
					}
					break
				}

				if h, ok := cell.val1.(hex_val); ok {
					list = append(list, h.val)
				} else {
					break
				}

				cur = cell.val2
			}

			if len(list) > 0 {
				rlist := make([]rune, len(list))
				for i, v := range list {
					rlist[i] = rune(v)
				}
				fmt.Println(string(rlist))
			} else {
				fmt.Println(v)
			}
		case hex_val:
			fmt.Printf("'%x' \n", v.val)
		default:
			fmt.Println(v)
		}

		return arg
	})

	//BOOLEANS
	e.setvar("true", closure_val{
		env:   e,
		param: "x",
		body: lambda_node{
			param: "y",
			body:  var_node{label: "x"},
		},
	})
	e.setvar("false", closure_val{
		env:   e,
		param: "x",
		body: lambda_node{
			param: "y",
			body:  var_node{label: "y"},
		},
	})
	e.setvar("not", closure_val{
		env:   e,
		param: "x",
		body: apply_node{
			arg: var_node{label: "true"},
			function: apply_node{
				arg:      var_node{label: "false"},
				function: var_node{label: "x"},
			},
		},
	})
	e.setvar("and", closure_val{
		env:   e,
		param: "x",
		body: lambda_node{
			param: "y",
			body: apply_node{
				arg: var_node{label: "true"},
				function: apply_node{
					arg:      var_node{label: "y"},
					function: var_node{label: "x"},
				},
			},
		},
	})
	e.setvar("or", closure_val{
		env:   e,
		param: "x",
		body: lambda_node{
			param: "y",
			body: apply_node{
				arg: var_node{label: "y"},
				function: apply_node{
					arg:      var_node{label: "false"},
					function: var_node{label: "x"},
				},
			},
		},
	})

	//BASIC MATHS
	e.define_primitive("add", func(x value) value {
		return prim_val{
			fn: func(y value) value {
				x2, ok := x.(hex_val)
				y2, ok2 := y.(hex_val)

				if !(ok && ok2) {
					panic("can not add non-number datatypes (yet?)")
				}

				val := x2.val + y2.val

				return hex_val{val: val}
			},
		}
	})
	e.define_primitive("minus", func(x value) value {
		return prim_val{
			fn: func(y value) value {
				x2, ok := x.(hex_val)
				y2, ok2 := y.(hex_val)

				if !(ok && ok2) {
					panic("can not minus non-number datatypes (yet?)")
				}

				val := x2.val - y2.val

				return hex_val{val: val}
			},
		}
	})
	e.define_primitive("mult", func(x value) value {
		return prim_val{
			fn: func(y value) value {
				x2, ok := x.(hex_val)
				y2, ok2 := y.(hex_val)

				if !(ok && ok2) {
					panic("can not multiply non-number datatypes (yet?)")
				}

				val := x2.val * y2.val

				return hex_val{val: val}
			},
		}
	})
	e.define_primitive("divi", func(x value) value {
		return prim_val{
			fn: func(y value) value {
				x2, ok := x.(hex_val)
				y2, ok2 := y.(hex_val)

				if !(ok && ok2) {
					panic("can not divide non-number datatypes (yet?)")
				}

				val := x2.val / y2.val

				return hex_val{val: val}
			},
		}
	})

	// ------- CONCELL STUFF ---------
	e.define_primitive("head", func(v value) value {
		c, ok := v.(concell_val)
		if !ok {
			panic("Head function only takes a concell value")
		}
		return c.val1
	})
	e.define_primitive("tail", func(v value) value {
		c, ok := v.(concell_val)
		if !ok {
			panic("Head function only takes a concell value")
		}
		return c.val2
	})

	// ------- EXTRAS --------
	e.define_primitive("c", func(v value) value {
		//CONVERTS HEX TO CHURCH NUMBERS
		n, ok := v.(hex_val)
		if !ok {
			panic("cannot \"Church-ify\" non-hex values (yet).")
		}

		var body node = var_node{label: "x"}
		for _ = range n.val {
			body = apply_node{
				function: var_node{label: "f"},
				arg:      body,
			}
		}

		return closure_val{
			env:   e,
			param: "f",
			body: lambda_node{
				param: "x",
				body:  body,
			},
		}
	})
	e.setvar("nil", closure_val{
		param: "x",
		body:  var_node{label: "x"},
		env:   e,
	})
}
