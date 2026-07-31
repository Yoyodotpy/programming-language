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
			fmt.Printf("'%x'", v.val)
		default:
			fmt.Println(v)
		}

		return arg
	})

	//BOOLEANS
	e.define_primitive("true", func(x value) value {
		return prim_val{
			fn: func(y value) value {
				return x
			},
		}
	})
	e.define_primitive("false", func(x value) value {
		return prim_val{
			fn: func(y value) value {
				return y
			},
		}
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
}
