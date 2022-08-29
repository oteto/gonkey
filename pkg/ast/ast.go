package ast

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/oteto/gonkey/pkg/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type       string
		Statements []Statement
	}{
		Type:       "RootNode",
		Statements: p.Statements,
	})
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// Statement interface を満たすので Program.Statements に格納できる
func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type       string
		Expression Expression
	}{
		Type:       "ExpressionStatementNode",
		Expression: es.Expression,
	})
}

type LetStatement struct {
	Token token.Token // LET
	Name  *Identifer
	Value Expression
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type       string
		Identifier *Identifer
		Value      Expression
	}{
		Type:       "LetStatementNode",
		Identifier: l.Name,
		Value:      l.Value,
	})
}

type ReturnStatement struct {
	Token       token.Token // RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string
		Value Expression
	}{
		Type:  "ReturnStatementNode",
		Value: rs.ReturnValue,
	})
}

type Identifer struct {
	Token token.Token // IDENT
	Value string
}

func (i *Identifer) expressionNode() {}

func (i *Identifer) String() string {
	return i.Value
}

func (i *Identifer) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string
		Name string
	}{
		Type: "IdentifierExpressionNode",
		Name: i.Value,
	})
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string
		Value int64
	}{
		Type:  "IntegerLiteralExpressionNode",
		Value: i.Value,
	})
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type     string
		Operator string
		Right    Expression
	}{
		Type:     "PrefixExpressionNode",
		Operator: pe.Operator,
		Right:    pe.Right,
	})
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type     string
		Operator string
		Left     Expression
		Right    Expression
	}{
		Type:     "InfixExpressionNode",
		Operator: ie.Operator,
		Left:     ie.Left,
		Right:    ie.Right,
	})
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string
		Value bool
	}{
		Type:  "BooleanExpressionNode",
		Value: b.Value,
	})
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type        string
		Condition   Expression
		Consequence *BlockStatement
		Alternative *BlockStatement
	}{
		Type:        "IfExpressionNode",
		Condition:   ie.Condition,
		Consequence: ie.Consequence,
		Alternative: ie.Alternative,
	})
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type       string
		Statements []Statement
	}{
		Type:       "BlockStatementNode",
		Statements: bs.Statements,
	})
}

type FunctionLiteral struct {
	Token      token.Token
	Patameters []*Identifer
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Patameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type        string
		Paramenters []*Identifer
		Body        *BlockStatement
	}{
		Type:        "FunctionLiteralExpressionNode",
		Paramenters: f.Patameters,
		Body:        f.Body,
	})
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

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

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type      string
		Function  Expression
		Arguments []Expression
	}{
		Type:      "FunctionCallExpressionNode",
		Function:  ce.Function,
		Arguments: ce.Arguments,
	})
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string
		Value string
	}{
		Type:  "StringLiteralExpressionNode",
		Value: sl.Value,
	})
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := make([]string, len(al.Elements))
	for i, el := range al.Elements {
		elements[i] = el.String()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type     string
		Elements []Expression
	}{
		Type:     "ArrayLiteralExpressionNode",
		Elements: al.Elements,
	})
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string
		Left  Expression
		Index Expression
	}{
		Type:  "IndexExpressionNode",
		Left:  ie.Left,
		Index: ie.Index,
	})
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := make([]string, len(hl.Pairs))
	i := 0
	for k, v := range hl.Pairs {
		pairs[i] = k.String() + ":" + v.String()
		i++
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) MarshalJSON() ([]byte, error) {
	type Pair struct {
		Key   Expression
		Value Expression
	}
	pairs := make([]Pair, len(hl.Pairs))
	i := 0
	for k, v := range hl.Pairs {
		pairs[i] = Pair{Key: k, Value: v}
		i++
	}

	return json.Marshal(&struct {
		Type  string
		Pairs []Pair
	}{
		Type:  "HashLiteralExpressionNode",
		Pairs: pairs,
	})
}
