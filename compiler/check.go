package compiler

import "github.com/dcaiafa/xml2csv/parser"

func (c *Compiler) check(order parser.Order, ast parser.AST) {
	if order == parser.Prefix {
		c.createSymbols(ast)
	} else {
		c.checkTypes(ast)
	}
}

func (c *Compiler) createSymbols(ast parser.AST) {
	switch ast := ast.(type) {
	case *parser.Program:
		c.globalNFA.AddStates(1)

	case *parser.Transform:
		c.curTransform = &transform{}
		c.curTransform.NFA.AddStates(1)
		c.curTransform.GlobalPath = c.createPath(&c.globalNFA, ast.Foreach)
		c.transforms = append(c.transforms, c.curTransform)

	case *parser.PathExpr:
		c.createPath(&c.curTransform.NFA, ast.Path)
		c.exprs[ast] = &expr{typ: exprPath}

	case *parser.BinaryExpr:
		c.exprs[ast] = &expr{typ: exprUnkn}

	case *parser.CallExpr:
		panic("not implemented")

	case *parser.LiteralExpr:
		var typ exprType
		switch ast.Val.(type) {
		case int64:
			typ = exprNum

		case string:
			typ = exprStr

		default:
			panic("not reached")
		}
		c.exprs[ast] = &expr{typ: typ, val: ast.Val}
	}
}

func (c *Compiler) checkTypes(ast parser.AST) {
S:
	switch ast := ast.(type) {
	case *parser.Program:

	case *parser.Transform:

	case *parser.PathExpr:

	case *parser.BinaryExpr:
		expr := c.exprs[ast]
		exprLHS := c.exprs[ast.LHS]
		exprRHS := c.exprs[ast.RHS]
		if exprLHS.typ == exprUnkn || exprRHS.typ == exprUnkn {
			// Type-checking has already failed for a child expression and an error
			// has already been logged.
			break
		}

		switch ast.Op {
		case parser.OpAdd, parser.OpSub, parser.OpMul, parser.OpDiv, parser.OpMod:
			operandTyp, exprTyp := c.checkArithExpr(exprLHS.typ, exprRHS.typ)
			if exprTyp == exprUnkn {
				break S
			}
			ast.LHS = c.castExpr(ast.LHS, operandTyp)
			ast.RHS = c.castExpr(ast.RHS, operandTyp)
			expr.typ = exprTyp

		case parser.OpEq, parser.OpLT, parser.OpLE, parser.OpGT, parser.OpGE:
			operandTyp := c.checkEqExpr(ast.Op, exprLHS.typ, exprRHS.typ)
			if operandTyp == exprUnkn {
				break S
			}
			ast.LHS = c.castExpr(ast.LHS, operandTyp)
			ast.RHS = c.castExpr(ast.RHS, operandTyp)
			expr.typ = exprBool

		case parser.OpAnd, parser.OpOr:
			operandTyp := c.checkLogExpr(ast.Op, exprLHS.typ, exprRHS.typ)
			if operandTyp == exprUnkn {
				break S
			}
			ast.LHS = c.castExpr(ast.LHS, operandTyp)
			ast.RHS = c.castExpr(ast.RHS, operandTyp)
			expr.typ = exprBool
		}

	case *parser.CallExpr:
		panic("not implemented")

	case *parser.LiteralExpr:
		var typ exprType
		switch ast.Val.(type) {
		case int64:
			typ = exprNum

		case string:
			typ = exprStr

		default:
			panic("not reached")
		}
		c.exprs[ast] = &expr{typ: typ, val: ast.Val}

	case *parser.CastExpr:
		// CastExprs are generated in this pass. They are not allowed to exist
		// prior to this.
		panic("not reached")
	}
}

func (c *Compiler) createPath(nfa *NFA, components []string) *path {
	startStateID := 0
	for _, component := range components {
		stateID := nfa.AddStates(1)
		baseSymID := c.names.IDFromName(component) * 10
		nfa.AddTransition(startStateID, stateID, baseSymID+1)
		nfa.AddTransition(stateID, startStateID, baseSymID+2)
		startStateID = stateID
	}
	nfa.SetAccepting(startStateID)
	return &path{nfaStateID: startStateID}
}

func (c *Compiler) checkArithExpr(
	lhs exprType,
	rhs exprType) (exprType, exprType) {

	if lhs != rhs || lhs != exprNum {
		c.fail("only numbers allowed in binary arithmetic expression")
		return exprUnkn, exprUnkn
	}

	return exprNum, exprNum
}

func (c *Compiler) checkEqExpr(
	op parser.Op,
	lhs exprType,
	rhs exprType) exprType {

	if lhs > rhs {
		lhs, rhs = rhs, lhs
	}

	operandType := exprUnkn

	switch lhs {
	case exprUnkn:
		panic("not reached")

	case exprBool:
		break

	case exprNum:
		if rhs == exprNum || rhs == exprPath {
			operandType = exprNum
		}

	case exprStr:
		if op != parser.OpEq {
			break
		}
		if rhs == exprStr || rhs == exprPath {
			operandType = exprStr
		}

	case exprPath:
		if rhs == exprPath {
			operandType = exprStr
		}

	default:
		panic("not reached")
	}

	if operandType == exprUnkn {
		c.fail("incompatible types in binary equality expression")
	}

	return operandType
}

func (c *Compiler) checkLogExpr(
	op parser.Op,
	lhs exprType,
	rhs exprType) exprType {

	if lhs != rhs || lhs != exprBool {
		c.fail("incompatible types in logical binary expression")
		return exprUnkn
	}

	return exprBool
}

func (c *Compiler) castExpr(exprAST parser.AST, typ exprType) parser.AST {
	e := c.exprs[exprAST]
	if e == nil {
		panic("assert failed")
	}
	if e.typ != typ {
		if e.val != nil {
			panic("not implemented")
		}
		exprAST = &parser.CastExpr{Expr: exprAST}
		c.exprs[exprAST] = &expr{typ: typ}
	}
	return exprAST
}
