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
	lambda_num int
}

func new_compiler() *compiler {
	mod := ir.NewModule()

	valtype := mod.NewTypeDef("Value", types.NewStruct(types.I8, types.I64))

	val_ptr_type := types.NewPointer(valtype)
	contype := mod.NewTypeDef("Pair", types.NewStruct(val_ptr_type, val_ptr_type))

	return &compiler{
		mod:     mod,
		vars:    make(map[string]value.Value),
		valtype: valtype,
		contype: contype,
	}
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
	var_ptr := c.block.NewAlloca(types.NewPointer(c.valtype))
	c.block.NewStore(val, var_ptr)
	c.vars[n.label] = var_ptr
	return val
}
func (c *compiler) compile_apply(n apply_node) value.Value {
	func_val := c.compile_node(n.function)
	arg_val := c.compile_node(n.arg)

	func_int := c.get_value(func_val)

	val_ptr_type := types.NewPointer(c.valtype)
	func_sig := types.NewPointer(types.NewFunc(val_ptr_type, val_ptr_type))
	func_ptr := c.block.NewIntToPtr(func_int, func_sig)

	return c.block.NewCall(func_ptr, arg_val)
}
func (c *compiler) compile_lambda(n lambda_node) value.Value {
	val_ptr_type := types.NewPointer(c.valtype)

	func_name := "lambda_" + strconv.Itoa(c.lambda_num)
	c.lambda_num++
	lambda_param := ir.NewParam(n.param, val_ptr_type)
	lambda_func := c.mod.NewFunc(func_name, val_ptr_type, lambda_param)

	prev_func := c.funcdef
	prev_block := c.block
	prev_vars := make(map[string]value.Value)
	for k, v := range c.vars {
		prev_vars[k] = v
	}

	c.funcdef = lambda_func
	c.block = lambda_func.NewBlock("entry")

	param_ptr := c.block.NewAlloca(val_ptr_type)
	c.block.NewStore(lambda_param, param_ptr)
	c.vars[n.param] = param_ptr

	body_val := c.compile_node(n.body)
	c.block.NewRet(body_val)

	c.funcdef = prev_func
	c.block = prev_block
	c.vars = prev_vars

	func_int := c.block.NewPtrToInt(lambda_func, types.I64)
	return c.put_value(closure_tag, func_int)
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
func irint(i int64, t *types.IntType) value.Value {
	return constant.NewInt(t, i)
}
