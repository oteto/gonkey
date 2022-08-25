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
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!false;", "!", false},
		{"!true;", "!", true},
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

		switch v := tt.value.(type) {
		case int64:
			testIntegerLiteral(t, exp.Right, v)
		case int:
			testIntegerLiteral(t, exp.Right, int64(v))
		case bool:
			testBooleanLiteral(t, exp.Right, v)
		default:
			t.Fatalf("no type")
		}

	}
}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
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

		testInfixExpression(t, stmt.Expression, tt.left, tt.operator, tt.right)
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
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b + c) + d",
			"((a + add((b + c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
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

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
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
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	tn := tokenizer.New(input)
	p := New(tn)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpresionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")
	checkStatementLength(t, exp.Consequence.Statements, 1)

	con, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpresionStatement. got=%T", exp.Consequence.Statements[0])
	}

	testIdentifier(t, con.Expression, "x")

	if exp.Alternative != nil {
		t.Fatalf("exp.Alternative is not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	tn := tokenizer.New(input)
	p := New(tn)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpresionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")
	checkStatementLength(t, exp.Consequence.Statements, 1)

	con, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpresionStatement. got=%T", exp.Consequence.Statements[0])
	}

	testIdentifier(t, con.Expression, "x")

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not *ast.ExpresionStatement. got=%T", exp.Consequence.Statements[0])
	}
	testIdentifier(t, alt.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	tn := tokenizer.New(input)
	p := New(tn)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpresionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Patameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got+%d", len(function.Patameters))
	}

	testLiteralExpression(t, function.Patameters[0], "x")
	testLiteralExpression(t, function.Patameters[1], "y")
	checkStatementLength(t, function.Body.Statements, 1)
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not *ast.ExpresionStatement. got=%T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input  string
		expect []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		p := New(tokenizer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Patameters) != len(tt.expect) {
			t.Errorf("parameter length wrong. want=%d, got=%d", len(tt.expect), len(function.Patameters))
		}

		for i, ident := range tt.expect {
			testLiteralExpression(t, function.Patameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	p := New(tokenizer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpresionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
	}

	testIdentifier(t, exp.Function, "add")
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong args length. want=3, got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input            string
		expectIdentifier string
		expectValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foo = y;", "foo", "y"},
	}

	for _, tt := range tests {
		p := New(tokenizer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)
		checkStatementLength(t, program.Statements, 1)

		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectIdentifier)
		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(t, val, tt.expectValue)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"Hello World"`
	p := New(tokenizer.New(input))
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)

	stmt := testExpressoinStatement(t, program)
	testStringLiteral(t, stmt.Expression, "Hello World")
}

func TestArrayLiteralExpression(t *testing.T) {
	input := `[1, 2 * 2, "hello"]`
	p := New(tokenizer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)
	stmt := testExpressoinStatement(t, program)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *asy.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) is not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testStringLiteral(t, array.Elements[2], "hello")
}

func TestParsingIndexExpression(t *testing.T) {
	input := "array[1 + 1]"
	p := New(tokenizer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)
	stmt := testExpressoinStatement(t, program)
	idx, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast,IndexExpression. got=%T", stmt.Expression)
	}
	testIdentifier(t, idx.Left, "array")
	testInfixExpression(t, idx.Index, 1, "+", 1)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	p := New(tokenizer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)
	stmt := testExpressoinStatement(t, program)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast,HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. want=3, got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for k, v := range hash.Pairs {
		key, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("key is not string literal. got=%T", k)
		}
		testIntegerLiteral(t, v, expected[key.Value])
	}
}

func TestParsingEmrtyHashLiteral(t *testing.T) {
	input := "{}"
	p := New(tokenizer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)
	stmt := testExpressoinStatement(t, program)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast,HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Fatalf("hash.Pairs has wrong length. want=0, got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	p := New(tokenizer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLength(t, program.Statements, 1)
	stmt := testExpressoinStatement(t, program)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast,HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. want=3, got=%d", len(hash.Pairs))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("key is not string literal. got=%T", k)
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Fatalf("No test function for key %q found.", literal.String())
		}
		testFunc(v)
	}
}

func testExpressoinStatement(t *testing.T, program *ast.Program) *ast.ExpressionStatement {
	t.Helper()
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpresionStatement. got=%T", program.Statements[0])
	}
	return stmt
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) {
	t.Helper()

	literal, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.StringLiteral. got=%T", exp)
	}
	if literal.Value != value {
		t.Fatalf("literal.Value is not %q. got=%q", value, literal.Value)
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	t.Helper()

	ident, ok := exp.(*ast.Identifer)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", exp)
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not %s. got=%s", value, ident.Value)
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s. got=%s", value, ident.TokenLiteral())
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	t.Helper()

	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp is not *ast.Boolean. got=%T", exp)
	}

	if b.Value != value {
		t.Errorf("b.value is not %t. got=%t", value, b.Value)
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral() is not %t. got=%s", value, b.TokenLiteral())
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) {
	t.Helper()

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not *ast.InfixExpression. got=%T", exp)
	}
	testLiteralExpression(t, opExp.Left, left)
	if opExp.Operator != operator {
		t.Fatalf("exp.Operator is not %s. got=%s", operator, opExp.Operator)
	}
	testLiteralExpression(t, opExp.Right, right)
}
