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
		{
			`{"name": "monkey"}[fn(){}]`,
			UNUSABLE_HASH_KEY + "FUNCTION",
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
		{`len([])`, 0},
		{`len([1, "a"])`, 2},
		{`first([1, "a"])`, 1},
		{`first([])`, nil},
		{`last([1,2,3])`, 3},
		{`last([])`, nil},
		{`rest([1,2,3])`, []int{2, 3}},
		{`rest([1])`, []int{}},
		{`rest([])`, nil},
		{`rest(rest([1,2,3]))`, []int{3}},
		{`push([], 1)`, []int{1}},
		{`push([1,2], 3)`, []int{1, 2, 3}},
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
		case []int:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
			}
			if len(expect) != len(arr.Elements) {
				t.Fatalf("wrong length of array. want=%d, got=%d", len(expect), len(arr.Elements))
			}
			for i, a := range arr.Elements {
				testIntegerObject(t, a, int64(expect[i]))
			}
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, "hello"]`
	evaluated := testEval(input)
	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not *object.Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) is not 3. got=%d", len(array.Elements))
	}
	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 4)
	testStringObject(t, array.Elements[2], "hello")
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect any
	}{
		{`[1,2,3][0]`, 1},
		{`["1",2,3][0]`, "1"},
		{`let a = [1,2,3]; a[1]`, 2},
		{`[1,2,3][1+1]`, 3},
		{`[1,2,3][3]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expect := tt.expect.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expect))
		case string:
			testStringObject(t, evaluated, expect)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{
	"one": 10 - 9,
	two: 1 + 1,
	"thr"+"ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6,
}
`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d, want=%d", len(result.Pairs), len(expected))
	}
	for k, v := range expected {
		pair, ok := result.Pairs[k]
		if !ok {
			t.Fatalf("no pair for given key in Pairs.")
		}
		testIntegerObject(t, pair.Value, v)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect any
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let k = 1;{1: 5}[k]`, 5},
		{`{}["foo"]`, nil},
		{`{true: 5}[true]`, 5},
		{`let k = false;{false: 5}[k]`, 5},
		{`let k = "bar";{false: 5, "bar": 1}[k]`, 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expect := tt.expect.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expect))
		case nil:
			testNullObject(t, evaluated)
		default:
			t.Fatalf("no expected type")
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
