package main

/*
#include <stdlib.h>
*/

import "C"
import (
	"fmt"
	"unsafe"
)

type value struct {
	tag     int8
	payload uint64
}

type concell struct {
	val1 *value
	val2 *value
}

const (
	hex_tag     int8 = 0
	con_tag     int8 = 1
	closure_tag int8 = 2
)

// HELPER FUNCTIONS
func alloc_val(tag int8, payload uint64) *value {
	v := (*value)(C.malloc(C.size_t(unsafe.Sizeof(value{}))))
	v.tag = tag
	v.payload = payload
	return v
}

func main() {}

// CONCELLS

//export head
func head(v *value) *value {
	concell := (*concell)(unsafe.Pointer(&v.payload))
	return concell.val1
}

//export tail
func tail(v *value) *value {
	concell := (*concell)(unsafe.Pointer(&v.payload))
	return concell.val2
}

// BASIC MATHS

//export go_add
func go_add(v *value) *value {
	concell := (*concell)(unsafe.Pointer(&v.payload))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return alloc_val(0, uint64(x+y))
}

//export go_mult
func go_mult(v *value) *value {
	concell := (*concell)(unsafe.Pointer(&v.payload))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return alloc_val(0, uint64(x*y))
}

//export go_min
func go_min(v *value) *value {
	concell := (*concell)(unsafe.Pointer(&v.payload))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return alloc_val(0, uint64(x-y))
}

//export go_div
func go_div(v *value) *value {
	concell := (*concell)(unsafe.Pointer(&v.payload))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return alloc_val(0, uint64(x/y))
}

//BOOLEAN OPERATORS
//EXTRAS

//export print
func print(v *value) *value {
	if v == nil {
		return nil
	}

	switch v.tag {
	case con_tag:
		cur := v
		var list []int64

		for {
			if cur == nil || cur.tag != con_tag {
				if cur != nil && cur.tag == hex_tag {
					list = append(list, int64(cur.payload))
				}
				break
			}

			concell := (*concell)(unsafe.Pointer(&cur.payload))

			if concell.val1 != nil && concell.val1.tag == hex_tag {
				list = append(list, int64(concell.val1.payload))
			} else {
				break
			}

			cur = concell.val2
		}

		if len(list) > 0 {
			rlist := make([]rune, len(list))
			for i, v := range list {
				rlist[i] = rune(v)
			}
			fmt.Println(string(rlist))
		} else {
			fmt.Println("<concell>")
		}
	case hex_tag:
		fmt.Printf("'%x'\n", v.payload)
	default:
		fmt.Println("<closure or unknown>")
	}
	return v
}
