package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

//Node an interface for all ast Nodes
type Node interface {
	TokenLiteral() string
	String() string
}

//Statement : an interface for all Statements
// statements do not produce values
type Statement interface {
	Node
	statementNode()
}

//Expression an interface for all expressions
// expressions produce values
type Expression interface {
	Node
	expressionNode()
}

//Program the root Node for any program
type Program struct {
	Statements []Statement // list of all top level statements
}

//TokenLiteral : A representation of the first token as string
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

//String : returns string representation of Node
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//Identifier : node for all identifiers
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

//expressionNode : implementation of the Expression interface
func (i *Identifier) expressionNode() {}

//TokenLiteral a  string representation of the token.IDENT token
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

//String : returns string representation of Node
func (i *Identifier) String() string { return i.Value }

//LetStatement : Node for let statements
// eg. let a = 12;
type LetStatement struct {
	Token token.Token // the token.LET token
	//iName : this is a statement because identifiers in other
	//parts of the language produce value
	Name  *Identifier
	Value Expression
}

//statementNode : implementer of Statement interface
func (ls *LetStatement) statementNode() {}

//TokenLiteral : a string representation of the token.LET statement
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

//String : returns string representation of Node
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())

	if ls.Value != nil {
		out.WriteString(" = ")
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

//ReturnStatement : statement Node to handle return statements
// eg. return 5;
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

//statementNode implementation of Node interface
func (rs *ReturnStatement) statementNode() {}

//TokenLiteral : returns 'return' string from token
// a string representation of token.RETURN token
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

//String : returns string representation of Node
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

//ExpressionStatement : A node that holds an expression
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

//statementNode : implementation of the Node interface
func (es *ExpressionStatement) statementNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

//String : returns string representation of Node
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

//IntegerLiteral : Node for holding integers
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

//expressionNode interface implementation for Expression Interface
func (il *IntegerLiteral) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

//String : returns string representation of Node
func (il *IntegerLiteral) String() string { return il.Token.Literal }

//PrefixExpression Node for all prefix expressions
// <prefix><expression>
// eg: !true, -5, etc
type PrefixExpression struct {
	Token    token.Token // the prefix token eg. !, -
	Operator string      // the operator
	Right    Expression
}

//expressionNode interface implementation for Expression Interface
func (pe *PrefixExpression) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

//String : returns string representation of Node
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

//InfixExpression Node for all infix operations
// eg 5 + 5, 5 - 5, etc
type InfixExpression struct {
	Token    token.Token // the Operator token e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

//expressionNode interface implementation for Expression Interface
func (oe *InfixExpression) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

//String : returns string representation of Node
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

//Boolean Node for true and false
type Boolean struct {
	Token token.Token
	Value bool
}

//expressionNode interface implementation for Expression Interface
func (b *Boolean) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

//String : returns string representation of Node
func (b *Boolean) String() string { return b.Token.Literal }

//BlockStatement node for handling a group of statements in a block
// eg. If-else-block or function-block
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

//statementNode interface implementation for Statement Interface
func (bs *BlockStatement) statementNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

//String : returns string representation of Node
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//IfExpression Node for hoding if and if-else statements block
type IfExpression struct {
	Token       token.Token //The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

//expressionNode interface implementation for Expression Interface
func (ie *IfExpression) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

//String : returns string representation of Node
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString("(")
	out.WriteString(ie.Condition.String())
	out.WriteString(")")
	out.WriteString(" ")
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

//FunctionLiteral Node for holding functions
type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

//expressionNode interface implementation for Expression Interface
func (fl *FunctionLiteral) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

//String : returns string representation of Node
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

//CallExpression Node to handle function calls
type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

//expressionNode interface implementation for Expression Interface
func (ce *CallExpression) expressionNode() {}

//TokenLiteral : a string representation of the expressionstatement node
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

//String : returns string representation of Node
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

//StringLiteral node to hold strings
type StringLiteral struct {
	Token token.Token
	Value string
}

//expressionNode implementation of the Expression interface
func (sl *StringLiteral) expressionNode() {}

//TokenLiteral returns a literal string representation of the node
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

//String returns a string form of the node
func (sl *StringLiteral) String() string { return sl.Token.Literal }

//ArrayLiteral node to hold arrays
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

//expressionNode implementation of the Expression interface
func (al *ArrayLiteral) expressionNode() {}

//TokenLiteral returns a literal string representation of the node
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

//String returns a string form of the node
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

//ArrayLiteral node to hold arrays
type IndexExpression struct {
	Token token.Token // the '[' token
	Left  Expression
	Index Expression
}

//expressionNode implementation of the Expression interface
func (ie *IndexExpression) expressionNode() {}

//TokenLiteral returns a literal string representation of the node
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

//String returns a string form of the node
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

//HashLiteral node to hold arrays
type HashLiteral struct {
	Token token.Token // the '[' token
	Pairs map[Expression]Expression
}

//expressionNode implementation of the Expression interface
func (hl *HashLiteral) expressionNode() {}

//TokenLiteral returns a literal string representation of the node
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

//String returns a string form of the node
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

//ClassStatement : Node for class statements
// eg. let a = 12;
type ClassStatement struct {
	Token token.Token // the token.CLASS token
	// the name of the class
	Name    *Identifier
	Parents []*Identifier
	Body    *BlockStatement
}

//statementNode : implementer of Statement interface
func (Cs *ClassStatement) statementNode() {}

//expressionNode implementer
func (Cs *ClassStatement) expressionNode() {}

//TokenLiteral : a string representation of the token.LET statement
func (Cs *ClassStatement) TokenLiteral() string { return Cs.Token.Literal }

//String : returns string representation of Node
func (Cs *ClassStatement) String() string {
	var out bytes.Buffer
	out.WriteString(Cs.TokenLiteral() + " ")
	out.WriteString(Cs.Name.String())
	parents := []string{}
	for _, i := range Cs.Parents {
		parents = append(parents, i.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(parents, ", "))
	out.WriteString(") {")
	out.WriteString(Cs.Body.String())
	out.WriteString("}")
	return out.String()

}

//ImportStatement : statement Node to handle import statements
// eg. return 5;
type ImportStatement struct {
	Token token.Token // the token.IMPORT token
	Value *StringLiteral
}

//expressionNode implementation of Node interface
func (Is *ImportStatement) expressionNode() {}

//TokenLiteral : returns 'return' string from token
// a string representation of token.RETURN token
func (Is *ImportStatement) TokenLiteral() string { return Is.Token.Literal }

//String : returns string representation of Node
func (Is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString(Is.TokenLiteral() + " ")
	out.WriteString(Is.Value.String())
	out.WriteString(";")
	return out.String()
}
