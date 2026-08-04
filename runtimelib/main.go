package main

/*
#include <stdlib.h>
#include <stdint.h>

typedef struct Value {
	int8_t tag;
	uint64_t payload;
} Value;

typedef Value* (*lambda_func_t)(void* env, Value* arg);

typedef struct Closure {
	lambda_func_t fn;
	void* env;
} Closure;

static Value* call_closure(Value* closure_val, Value* arg) {
	if (!closure_val || closure_val->tag != 2) return arg;
	Closure* cl = (Closure*)(closure_val->payload);
	return cl->fn(cl->env, arg);
}
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
	hex_tag int8 = 0
	con_tag int8 = 1

// closure_tag int8 = 2
)

func alloc_val(tag int8, payload uint64) *value {
	v := (*value)(C.malloc(C.size_t(unsafe.Sizeof(value{}))))
	v.tag = tag
	v.payload = payload
	return v
}

func main() {}

// CONCELLS

//export head
func head(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	concell := (*concell)(unsafe.Pointer(uintptr(v.payload)))
	return unsafe.Pointer(concell.val1)
}

//export tail
func tail(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	concell := (*concell)(unsafe.Pointer(uintptr(v.payload)))
	return unsafe.Pointer(concell.val2)
}

// BASIC MATHS

//export go_add
func go_add(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	concell := (*concell)(unsafe.Pointer(uintptr(v.payload)))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return unsafe.Pointer(alloc_val(0, uint64(x+y)))
}

//export go_mult
func go_mult(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	concell := (*concell)(unsafe.Pointer(uintptr(v.payload)))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return unsafe.Pointer(alloc_val(0, uint64(x*y)))
}

//export go_min
func go_min(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	concell := (*concell)(unsafe.Pointer(uintptr(v.payload)))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return unsafe.Pointer(alloc_val(0, uint64(x-y)))
}

//export go_div
func go_div(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	concell := (*concell)(unsafe.Pointer(uintptr(v.payload)))

	x := int64(concell.val1.payload)
	y := int64(concell.val2.payload)

	return unsafe.Pointer(alloc_val(0, uint64(x/y)))
}

//EXTRAS

//export print
func print(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
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

			concell := (*concell)(unsafe.Pointer(uintptr(cur.payload)))

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
	return ptr
}

//export go_repeat
func go_repeat(env unsafe.Pointer, ptr unsafe.Pointer) unsafe.Pointer {
	v := (*value)(ptr)
	if v == nil {
		return nil
	}

	c := (*concell)(unsafe.Pointer(uintptr(v.payload)))

	if c.val1 == nil || c.val1.tag != hex_tag {
		return nil
	}

	count := int64(c.val1.payload)
	fn_val := c.val2
	var last_res *value = nil
	for range int64(count) {
		dummy := alloc_val(0, 0)
		res, _ := C.call_closure(
			(*C.Value)(unsafe.Pointer(fn_val)),
			(*C.Value)(unsafe.Pointer(dummy)),
		)
		last_res = (*value)(unsafe.Pointer(res))
	}

	return unsafe.Pointer(last_res)
}
