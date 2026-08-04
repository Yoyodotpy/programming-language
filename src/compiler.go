package main

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

const (
	hex_tag     int64 = 0
	con_tag     int64 = 1
	closure_tag int64 = 2
)

type compiler struct {
	mod        *ir.Module
	funcdef    *ir.Func
	block      *ir.Block
	vars       map[string]value.Value
	valtype    types.Type
	contype    types.Type
	lambtype   types.Type
	lambda_num int
}

func new_compiler() *compiler {
	mod := ir.NewModule()

	valtype := mod.NewTypeDef("Value", types.NewStruct(types.I8, types.I64))

	val_ptr_type := types.NewPointer(valtype)
	contype := mod.NewTypeDef("Pair", types.NewStruct(val_ptr_type, val_ptr_type))

	env_ptr_type := types.NewPointer(val_ptr_type)
	func_sig := types.NewFunc(val_ptr_type, env_ptr_type, val_ptr_type)
	func_ptr_type := types.NewPointer(func_sig)

	lambtype := mod.NewTypeDef("Closure", types.NewStruct(func_ptr_type, env_ptr_type))

	c := &compiler{
		mod:      mod,
		vars:     make(map[string]value.Value),
		valtype:  valtype,
		contype:  contype,
		lambtype: lambtype,
	}

	return c
}

func (c *compiler) compile_node(n node) value.Value {
	switch n := n.(type) {
	case hex_node:
		return c.put_value(hex_tag, irint(n.val, types.I64))
	case var_node:
		alloc_ptr, ok := c.vars[n.label]
		if !ok {
			panic("Variable does not exist: " + n.label)
		}
		return c.block.NewLoad(types.NewPointer(c.valtype), alloc_ptr)
	case define_node:
		return c.compile_define(n)
	case apply_node:
		return c.compile_apply(n)
	case lambda_node:
		return c.compile_lambda(n)
	case concell_node:
		return c.compile_concell(n)
	default:
		panic("Compiler does not recognize node")
	}
}

func (c *compiler) compile_define(n define_node) value.Value {
	val := c.compile_node(n.value)
	c.bind_local(n.label, val)
	return val
}
func (c *compiler) compile_apply(n apply_node) value.Value {
	func_val := c.compile_node(n.function)
	arg_val := c.compile_node(n.arg)

	closure_int := c.get_value(func_val)

	closure_ptr_type := types.NewPointer(c.lambtype)
	closure_ptr := c.block.NewIntToPtr(closure_int, closure_ptr_type)

	val_ptr_type := types.NewPointer(c.valtype)
	env_ptr_type := types.NewPointer(val_ptr_type)
	func_sig := types.NewPointer(types.NewFunc(val_ptr_type, env_ptr_type, val_ptr_type))

	func_ptr_gep := c.get_struct_field(c.lambtype, closure_ptr, 0)
	func_ptr := c.block.NewLoad(func_sig, func_ptr_gep)

	env_ptr_gep := c.get_struct_field(c.lambtype, closure_ptr, 1)
	env_ptr := c.block.NewLoad(env_ptr_type, env_ptr_gep)

	return c.block.NewCall(func_ptr, env_ptr, arg_val)

}
func (c *compiler) compile_lambda(n lambda_node) value.Value {
	val_ptr_type := types.NewPointer(c.valtype)
	env_ptr_type := types.NewPointer(val_ptr_type)

	label := "lambda_" + strconv.Itoa(c.lambda_num)
	c.lambda_num++

	env_param := ir.NewParam("env", env_ptr_type)
	lambda_param := ir.NewParam(n.param, val_ptr_type)
	lambda_func := c.mod.NewFunc(label, val_ptr_type, env_param, lambda_param)

	//REMOVES DUPLICATES
	raw_free_vars := n.free_vars()
	var free_vars []string
	seen := make(map[string]bool)
	for _, v := range raw_free_vars {
		if !seen[v] {
			seen[v] = true
			free_vars = append(free_vars, v)
		}
	}

	closure_ptr := c.block.NewAlloca(c.lambtype)

	func_ptr_gep := c.get_struct_field(c.lambtype, closure_ptr, 0)
	c.block.NewStore(lambda_func, func_ptr_gep)

	var env_array_ptr value.Value = constant.NewNull(env_ptr_type)
	if len(free_vars) > 0 {
		env_array_type := types.NewArray(uint64(len(free_vars)), val_ptr_type)
		env_alloc := c.block.NewAlloca(env_array_type)
		env_array_ptr = c.block.NewBitCast(env_alloc, env_ptr_type)

		for i, fv := range free_vars {
			val_ptr, ok := c.vars[fv]
			if !ok {
				panic("free varibale not found: " + fv)
			}
			val := c.block.NewLoad(val_ptr_type, val_ptr)

			idx_gep := c.get_struct_field(env_array_type, env_alloc, int64(i))
			c.block.NewStore(val, idx_gep)
		}
	}

	env_ptr_gep := c.get_struct_field(c.lambtype, closure_ptr, 1)
	c.block.NewStore(env_array_ptr, env_ptr_gep)

	closure_int := c.block.NewPtrToInt(closure_ptr, types.I64)

	prev_func := c.funcdef
	prev_block := c.block
	prev_vars := c.vars

	c.funcdef = lambda_func
	c.block = lambda_func.NewBlock("entry")
	c.vars = make(map[string]value.Value)

	if len(free_vars) > 0 {
		for i, fv := range free_vars {
			idx_gep := c.block.NewGetElementPtr(val_ptr_type, env_param, irint(int64(i), types.I32))
			captured_val := c.block.NewLoad(val_ptr_type, idx_gep)

			c.bind_local(fv, captured_val)
		}
	}

	c.bind_local(n.param, lambda_param)

	body_val := c.compile_node(n.body)
	c.block.NewRet(body_val)

	c.funcdef = prev_func
	c.block = prev_block
	c.vars = prev_vars

	return c.put_value(closure_tag, closure_int)

}
func (c *compiler) compile_concell(n concell_node) value.Value {
	val1 := c.compile_node(n.val1)
	val2 := c.compile_node(n.val2)

	con_ptr := c.block.NewAlloca(c.contype)

	val1_ptr := c.block.NewGetElementPtr(c.contype, con_ptr, irint(0, types.I32), irint(0, types.I32))
	c.block.NewStore(val1, val1_ptr)

	val2_ptr := c.block.NewGetElementPtr(c.contype, con_ptr, irint(0, types.I32), irint(0, types.I32))
	c.block.NewStore(val2, val2_ptr)

	con_int := c.block.NewPtrToInt(con_ptr, types.I64)

	return c.put_value(con_tag, con_int)
}

func (c *compiler) compile_program(ast []node) string {
	c.funcdef = c.mod.NewFunc("main", types.I32)
	c.block = c.funcdef.NewBlock("entry")

	c.bind_prim("go_add")
	c.bind_prim("go_mult")
	c.bind_prim("go_min")
	c.bind_prim("go_div")
	c.bind_prim("go_repeat")
	c.bind_prim("head")
	c.bind_prim("tail")
	c.bind_prim("print")

	for _, n := range ast {
		c.compile_node(n)
	}

	c.block.NewRet(constant.NewInt(types.I32, 0))

	return c.mod.String()
}

func (c *compiler) put_value(tag int64, payload value.Value) value.Value {
	val_ptr := c.block.NewAlloca(c.valtype)

	tag_ptr := c.block.NewGetElementPtr(c.valtype, val_ptr, irint(0, types.I32), irint(0, types.I32))
	c.block.NewStore(irint(tag, types.I8), tag_ptr)

	payload_ptr := c.block.NewGetElementPtr(c.valtype, val_ptr, irint(0, types.I32), irint(1, types.I32))
	c.block.NewStore(payload, payload_ptr)

	return val_ptr
}
func (c *compiler) get_value(ptr value.Value) value.Value {
	payload_ptr := c.block.NewGetElementPtr(c.valtype, ptr, irint(0, types.I32), irint(1, types.I32))
	return c.block.NewLoad(types.I64, payload_ptr)
}
func (c *compiler) get_struct_field(struct_type types.Type, struct_ptr value.Value, index int64) value.Value {
	return c.block.NewGetElementPtr(struct_type, struct_ptr, irint(0, types.I32), irint(index, types.I32))
}
func (c *compiler) bind_local(name string, val value.Value) {
	val_ptr_type := types.NewPointer(c.valtype)
	var_ptr := c.block.NewAlloca(val_ptr_type)
	c.block.NewStore(val, var_ptr)
	c.vars[name] = var_ptr
}
func (c *compiler) bind_prim(label string) {
	val_ptr_type := types.NewPointer(c.valtype)
	env_ptr_type := types.NewPointer(val_ptr_type)

	ext_func := c.mod.NewFunc(label, val_ptr_type, ir.NewParam("arg", val_ptr_type))

	wrapper_name := label + "_wrapper"
	wrapper_func := c.mod.NewFunc(wrapper_name, val_ptr_type, ir.NewParam("env", env_ptr_type), ir.NewParam("arg", val_ptr_type))
	wrapper_block := wrapper_func.NewBlock("entry")

	call_res := wrapper_block.NewCall(ext_func, wrapper_func.Params[1])
	wrapper_block.NewRet(call_res)

	closure_ptr := c.block.NewAlloca(c.lambtype)
	c.block.NewStore(wrapper_func, c.get_struct_field(c.lambtype, closure_ptr, 0))
	c.block.NewStore(constant.NewNull(env_ptr_type), c.get_struct_field(c.lambtype, closure_ptr, 1))

	closure_int := c.block.NewPtrToInt(closure_ptr, types.I64)
	boxed_closure := c.put_value(closure_tag, closure_int)
	c.bind_local(label, boxed_closure)
}
func irint(i int64, t *types.IntType) value.Value {
	return constant.NewInt(t, i)
}
