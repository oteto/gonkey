package evaluator

import (
	"fmt"
	"testing"

	"github.com/oteto/gonkey/pkg/object"
	"github.com/oteto/gonkey/pkg/parser"
	"github.com/oteto/gonkey/pkg/tokenizer"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"5;", 5},
		{"10", 10},
		{"-5;", -5},
		{"-10", -10},
		{"5 + 5 + 5 - 10", 5},
		{"2 * 2 * 2 * 2", 16},
		{"-5 + 10 + -5", 0},
		{"5 + 2 * 10", 25},
		{"3 * (1 + 2) + 5", 14},
		{"(2 - 5) * ((5 / 5) + 3)", -12},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expect)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true;", true},
		{"false;", false},
		{"1 < 2;", true},
		{"1 > 2;", false},
		{"1 < 1;", false},
		{"1 > 1;", false},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"1 == 2;", false},
		{"1 != 2;", true},
		{"true == true;", true},
		{"false == false;", true},
		{"true == false;", false},
		{"true != true;", false},
		{"true != false;", true},
		{"(1 < 2) == true;", true},
		{"(1 > 2) == true;", false},
		{"(1 > 2) != true;", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expect)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expect)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{"if (true) { 10; }", 10},
		{"if (false) { 10; }", nil},
		{"if (1) { 10; }", 10},
		{"if (1 < 2) { 10; }", 10},
		{"if (1 > 2) { 10; }", nil},
		{"if (true) { 10; } else { 20 }", 10},
		{"if (false) { 10; } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expect.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"return 10", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 7 + 9;", 10},
		{"9; return 1 + 1; 10;", 2},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expect)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			"5 + true",
			TYPE_MISMATCH_ERROR_PREFIX + "INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			TYPE_MISMATCH_ERROR_PREFIX + "INTEGER + BOOLEAN",
		},
		{
			"-true",
			UNKOWN_OPERATOR_ERROR_PREFIX + "-BOOLEAN",
		},
		{
			"false + true",
			UNKOWN_OPERATOR_ERROR_PREFIX + "BOOLEAN + BOOLEAN",
		},
		{
			"5; false + true; 5;",
			UNKOWN_OPERATOR_ERROR_PREFIX + "BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			UNKOWN_OPERATOR_ERROR_PREFIX + "BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	reutrn 1;
}
			`,
			UNKOWN_OPERATOR_ERROR_PREFIX + "BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			IDENTIFIER_NOT_FOUND_ERROR_PREFIX + "foobar",
		},
		{
			`"" - ""`,
			UNKOWN_OPERATOR_ERROR_PREFIX + "STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("no error object returned. got=%T (%+v)", evaluated, evaluated)
		}
		if errObj.Message != tt.expect {
			t.Fatalf("wrong error message. want=%q, got=%q", tt.expect, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let x = 5; x;", 5},
		{"let x = 5 * 5; x;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = b + a; c;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expect)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0].String())
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 4);", 9},
		{"let add = fn(x, y) { x + y; }; add(5, add(3, 4));", 12},
		{"fn(x) { x; }(5);", 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expect)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World"`
	evaluated := testEval(input)
	testStringObject(t, evaluated, "Hello World")
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World"`
	evaluated := testEval(input)
	testStringObject(t, evaluated, "Hello World")
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input  string
		expect any
	}{
		{`len("")`, 0},
		{`len("aiueo")`, 5},
		{`len("hello world")`, 11},
		{`len(1)`, fmt.Sprintf(BUILTIN_ARGUMENT_TYPE_ERRROR, "len", object.INTEGER_OBJECT)},
		{`len("", "a")`, fmt.Sprintf(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, 2, 1)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expect := tt.expect.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expect))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != expect {
				t.Fatalf("wrong error message. want=%q, got=%q", expect, errObj.Message)
			}
		}
	}
}

func testEval(input string) object.Object {
	p := parser.New(tokenizer.New(input))
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testStringObject(t *testing.T, obj object.Object, value string) {
	t.Helper()
	str, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", obj, obj)
	}
	if str.Value != value {
		t.Fatalf("str.Value is not %q. got=%q", value, str.Value)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, value int64) {
	t.Helper()

	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", obj, obj)
	}

	if result.Value != value {
		t.Fatalf("object has wrong value. want=%d, got=%d", value, result.Value)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, value bool) {
	t.Helper()

	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+v)", obj, obj)
	}

	if result.Value != value {
		t.Fatalf("object has wrong value. want=%t, got=%t", value, result.Value)
	}
}

func testNullObject(t *testing.T, obj object.Object) {
	t.Helper()

	if obj != NULL {
		t.Fatalf("object is not NULL. got=%T (%+v)", obj, obj)
	}
}
