package parser

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/operators"
	"github.com/chrispyles/slow/types"
)

// -------------------------------------------------------------------------------------------------
// Assignment node

type assignmentTarget struct {
	variable  string
	attribute *AttributeNode
	index     *IndexNode
}

type AssignmentNode struct {
	Left  assignmentTarget
	Right execute.Expression
}

func (n *AssignmentNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Right.Execute(e)
	if err != nil {
		return nil, err
	}
	if n := n.Left.variable; n != "" {
		return e.Set(n, expr)
	}
	if an := n.Left.attribute; an != nil {
		val, err := an.Left.Execute(e)
		if err != nil {
			return nil, err
		}
		err = val.SetAttribute(an.Right, expr)
		if err != nil {
			return nil, err
		}
		return expr, nil
	}
	if in := n.Left.index; in != nil {
		val, err := in.Container.Execute(e)
		if err != nil {
			return nil, err
		}
		idx, err := in.Index.Execute(e)
		if err != nil {
			return nil, err
		}
		if err := val.SetIndex(idx, expr); err != nil {
			return nil, err
		}
		return val, nil
	}
	panic("unhandled target case in AssignmentNode.Execute")
}

// -------------------------------------------------------------------------------------------------
// Attribute node

type AttributeNode struct {
	Left  execute.Expression
	Right string
}

func (n *AttributeNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Left.Execute(e)
	if err != nil {
		return nil, err
	}
	return expr.GetAttribute(n.Right)
}

// -------------------------------------------------------------------------------------------------
// Binary operator node

type BinaryOpNode struct {
	Op    *operators.BinaryOperator
	Left  execute.Expression
	Right execute.Expression
}

func (n *BinaryOpNode) Execute(e *execute.Environment) (execute.Value, error) {
	le, err := n.Left.Execute(e)
	if err != nil {
		return nil, err
	}
	re, err := n.Right.Execute(e)
	if err != nil {
		return nil, err
	}
	val, err := n.Op.Value(le, re)
	if err != nil {
		return nil, err
	}
	if n.Op.IsReassignmentOperator() {
		switch left := n.Left.(type) {
		case *VariableNode:
			return e.Set(left.Name, val)
		case *AttributeNode:
			expr, err := left.Left.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := expr.SetAttribute(left.Right, val); err != nil {
				return nil, err
			}
			return val, nil
		case *IndexNode:
			expr, err := left.Container.Execute(e)
			if err != nil {
				return nil, err
			}
			idx, err := left.Index.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := expr.SetIndex(idx, val); err != nil {
				return nil, err
			}
			return val, nil
		default:
			panic("unexpected node type in reassignment operator")
		}
	}
	return val, nil
}

// -------------------------------------------------------------------------------------------------
// Break node

type breakError struct{}

func (*breakError) Error() string { return "" }

type BreakNode struct{}

func (*BreakNode) Execute(e *execute.Environment) (execute.Value, error) {
	return types.Null, &breakError{}
}

// -------------------------------------------------------------------------------------------------
// Call node

type CallNode struct {
	Func execute.Expression
	Args []execute.Expression
}

func (n *CallNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Func.Execute(e)
	if err != nil {
		return nil, err
	}
	callable, err := expr.ToCallable()
	if err != nil {
		return nil, err
	}
	var args []execute.Value
	for _, a := range n.Args {
		v, err := a.Execute(e)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	return callable.Call(e.NewFrame(), args...)
}

// -------------------------------------------------------------------------------------------------
// Constant node

type ConstantNode struct {
	Value execute.Value
}

func (n *ConstantNode) Execute(e *execute.Environment) (execute.Value, error) {
	return n.Value, nil
}

// -------------------------------------------------------------------------------------------------
// Continue node

type continueError struct{}

func (*continueError) Error() string { return "" }

type ContinueNode struct{}

func (*ContinueNode) Execute(e *execute.Environment) (execute.Value, error) {
	return types.Null, &continueError{}
}

// -------------------------------------------------------------------------------------------------
// Fallthrough node

type fallthroughError struct{}

func (*fallthroughError) Error() string { return "" }

type FallthroughNode struct{}

func (*FallthroughNode) Execute(e *execute.Environment) (execute.Value, error) {
	return types.Null, &fallthroughError{}
}

// -------------------------------------------------------------------------------------------------
// For node

type ForNode struct {
	IterName string
	Iter     execute.Expression
	Body     execute.Block
}

func (n *ForNode) Execute(e *execute.Environment) (execute.Value, error) {
	var val execute.Value
	expr, err := n.Iter.Execute(e)
	if err != nil {
		return nil, err
	}
	iter, err := expr.ToIterator()
	if err != nil {
		return nil, err
	}
	for iter.HasNext() {
		frame := e.NewFrame()
		if err := frame.Declare(n.IterName); err != nil {
			return nil, err
		}
		expr, err := iter.Next()
		if err != nil {
			return nil, err
		}
		if _, err := frame.Set(n.IterName, expr); err != nil {
			return nil, err
		}
		val, err = n.Body.Execute(frame)
		if _, ok := err.(*breakError); ok {
			break
		} else if _, ok := err.(*continueError); ok {
			continue
		} else if err != nil {
			return nil, err
		}
	}
	return val, nil
}

// -------------------------------------------------------------------------------------------------
// Func node

type FuncNode struct {
	Name     string
	ArgNames []string
	Body     execute.Block
}

func (n *FuncNode) Execute(e *execute.Environment) (execute.Value, error) {
	if err := e.Declare(n.Name); err != nil {
		return nil, err
	}
	ft := types.NewFunc(n.Name, n.ArgNames, n.Body)
	return e.Set(n.Name, ft)
}

// -------------------------------------------------------------------------------------------------
// If node

type IfNode struct {
	Cond     execute.Expression
	Body     execute.Block
	ElseBody execute.Block
}

func (n *IfNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Cond.Execute(e)
	if err != nil {
		return nil, err
	}
	frame := e.NewFrame()
	if expr.ToBool() {
		return n.Body.Execute(frame)
	} else {
		return n.ElseBody.Execute(frame)
	}
}

// -------------------------------------------------------------------------------------------------
// Index node

type IndexNode struct {
	Container execute.Expression
	Index     execute.Expression
}

func (n *IndexNode) Execute(e *execute.Environment) (execute.Value, error) {
	container, err := n.Container.Execute(e)
	if err != nil {
		return nil, err
	}
	index, err := n.Index.Execute(e)
	if err != nil {
		return nil, err
	}
	return container.GetIndex(index)
}

// -------------------------------------------------------------------------------------------------
// List node

type ListNode struct {
	Values []execute.Expression
}

func (n *ListNode) Execute(e *execute.Environment) (execute.Value, error) {
	vs := make([]execute.Value, len(n.Values))
	for i, expr := range n.Values {
		v, err := expr.Execute(e)
		if err != nil {
			return nil, err
		}
		vs[i] = v
	}
	return types.NewList(vs), nil
}

// -------------------------------------------------------------------------------------------------
// Map node

type MapNode struct {
	Values [][]execute.Expression
}

func (n *MapNode) Execute(e *execute.Environment) (execute.Value, error) {
	m := types.NewMap()
	for _, kv := range n.Values {
		k, err := kv[0].Execute(e)
		if err != nil {
			return nil, err
		}
		v, err := kv[1].Execute(e)
		if err != nil {
			return nil, err
		}
		if _, err := m.Set(k, v); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// -------------------------------------------------------------------------------------------------
// Return node

type ReturnNode struct {
	Value execute.Expression
}

func (n *ReturnNode) Execute(e *execute.Environment) (execute.Value, error) {
	value, err := n.Value.Execute(e)
	if err != nil {
		return nil, err
	}
	return nil, &types.ReturnError{Value: value}
}

// -------------------------------------------------------------------------------------------------
// Switch node

type switchCase struct {
	CaseExpr execute.Expression
	Body     execute.Block
}

type SwitchNode struct {
	Value       execute.Expression
	Cases       []*switchCase
	DefaultCase execute.Block
}

func (n *SwitchNode) Execute(e *execute.Environment) (execute.Value, error) {
	valueExpr, err := n.Value.Execute(e)
	if err != nil {
		return nil, err
	}
	frame := e.NewFrame()
	var fallThrough bool
	for _, c := range n.Cases {
		caseExpr, err := c.CaseExpr.Execute(e)
		if err != nil {
			return nil, err
		}
		if valueExpr.Equals(caseExpr) || fallThrough {
			fallThrough = false
			val, err := c.Body.Execute(frame)
			if _, ok := err.(*fallthroughError); ok {
				fallThrough = true
				continue
			} else {
				return val, err
			}
		}
	}
	return n.DefaultCase.Execute(frame)
}

// -------------------------------------------------------------------------------------------------
// Unary operator node

type UnaryOpNode struct {
	Op   *operators.UnaryOperator
	Expr execute.Expression
}

func (n *UnaryOpNode) Execute(e *execute.Environment) (execute.Value, error) {
	operand, err := n.Expr.Execute(e)
	if err != nil {
		return nil, err
	}
	val, err := n.Op.Value(operand)
	if err != nil {
		return nil, err
	}
	if n.Op.IsReassignmentOperator() {
		// Reassignment operators return the value of the operand BEFORE the operation, but update its
		// value in the environment/object.
		switch expr := n.Expr.(type) {
		case *VariableNode:
			_, err := e.Set(expr.Name, val)
			if err != nil {
				return nil, err
			}
			return operand, nil
		case *AttributeNode:
			lVal, err := expr.Left.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := lVal.SetAttribute(expr.Right, val); err != nil {
				return nil, err
			}
			return operand, nil
		case *IndexNode:
			lVal, err := expr.Container.Execute(e)
			if err != nil {
				return nil, err
			}
			idx, err := expr.Index.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := lVal.SetIndex(idx, val); err != nil {
				return nil, err
			}
			return operand, nil
		default:
			panic("unexpected node type in reassignment operator")
		}
	}
	return val, nil
}

// -------------------------------------------------------------------------------------------------
// Var node

type VarNode struct {
	Name    string
	IsConst bool
	// Value of the expression if it is assigned; may be nil
	Value execute.Expression
}

func (n *VarNode) Execute(e *execute.Environment) (execute.Value, error) {
	var val execute.Value
	if n.Value != nil {
		var err error
		val, err = n.Value.Execute(e)
		if err != nil {
			return nil, err
		}
	}
	if n.IsConst {
		if val == nil {
			panic("encountered const node with no value")
		}
		return e.DeclareConst(n.Name, val)
	}
	if err := e.Declare(n.Name); err != nil {
		return nil, err
	}
	if val != nil {
		return e.Set(n.Name, val)
	}
	return types.Null, nil
}

// -------------------------------------------------------------------------------------------------
// Variable node (represents a variable access, not a declaration)

type VariableNode struct {
	Name string
}

func (n *VariableNode) Execute(e *execute.Environment) (execute.Value, error) {
	return e.Get(n.Name)
}

// -------------------------------------------------------------------------------------------------
// While node

type WhileNode struct {
	Cond execute.Expression
	Body execute.Block
}

func (n *WhileNode) Execute(e *execute.Environment) (execute.Value, error) {
	frame := e.NewFrame()
	var val execute.Value
	for {
		expr, err := n.Cond.Execute(e)
		if err != nil {
			return nil, err
		}
		if expr.ToBool() {
			val, err = n.Body.Execute(frame)
			if _, ok := err.(*breakError); ok {
				break
			} else if _, ok := err.(*continueError); ok {
				continue
			} else if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return val, nil
}
