package parser

import (
	"fmt"
	"testing"

	"github.com/oteto/gonkey/pkg/ast"
	"github.com/oteto/gonkey/pkg/tokenizer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo = 12345;
`
	tn := tokenizer.New(input)
	p := New(tn)

	program := p.ParseProgram()

	checkParserErrors(t, p)
	notNilProgram(t, program)
	checkStatementLength(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expectedIdentifer)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 12345;
`
	tn := tokenizer.New(input)
	p := New(tn)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	notNilProgram(t, program)
	checkStatementLength(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ReturnStatement. got: %T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral() not 'return'. got: %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foo;"

	tn := tokenizer.New(input)
	p := New(tn)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	notNilProgram(t, program)
	checkStatementLength(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifer)
	if !ok {
		t.Fatalf("exp not ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foo" {
		t.Errorf("ident.Value not %s. got=%s", "foo", ident.Value)
	}
	if ident.TokenLiteral() != "foo" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", "foo", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	tn := tokenizer.New(input)
	p := New(tn)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	notNilProgram(t, program)
	checkStatementLength(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range tests {
		tn := tokenizer.New(tt.input)
		p := New(tn)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		notNilProgram(t, program)
		checkStatementLength(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
		}
		testIntegerLiteral(t, exp.Right, tt.value)
	}
}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range tests {
		tn := tokenizer.New(tt.input)
		p := New(tn)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		notNilProgram(t, program)
		checkStatementLength(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp not ast.InfixExpression. got=%T", stmt.Expression)
		}
		testIntegerLiteral(t, exp.Left, tt.left)

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
		}

		testIntegerLiteral(t, exp.Right, tt.right)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		tn := tokenizer.New(tt.input)
		p := New(tn)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	t.Helper()

	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not ast.IntegerLiteral. got=%T", il)
	}
	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d. got=%s", value, integer.TokenLiteral())
	}

}

func testLetStatement(t *testing.T, stmt ast.Statement, identifer string) {
	t.Helper()

	if stmt.TokenLiteral() != "let" {
		t.Fatalf("stmt.TokenLiteral() not 'let'. got: %q", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.LetStatement. got: %T", stmt)
	}

	if letStmt.Name.Value != identifer {
		t.Fatalf("letStmt.Name.Value is not %q. got: %q", identifer, letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != identifer {
		t.Fatalf("letStmt.Name.TokenLiteral() is not %q. got: %q", identifer, letStmt.Name.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors.", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func notNilProgram(t *testing.T, program *ast.Program) {
	t.Helper()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
}

func checkStatementLength(t *testing.T, stmts []ast.Statement, l int) {
	t.Helper()
	if len(stmts) != l {
		t.Fatalf("program Statement does not contain %d statements. got: %d", l, len(stmts))
	}
}
