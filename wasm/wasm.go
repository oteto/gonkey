package main

import (
	"bytes"
	"encoding/json"
	"syscall/js"

	"github.com/oteto/gonkey/pkg/evaluator"
	"github.com/oteto/gonkey/pkg/object"
	"github.com/oteto/gonkey/pkg/parser"
	"github.com/oteto/gonkey/pkg/tokenizer"
)

func tokenize(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	t := tokenizer.New(input)
	tokenJson, err := json.Marshal(t)
	if err != nil {
		return `{"error": "error json.Marshal."}`
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, tokenJson, "", "\t")
	if err != nil {
		return `{"error": "error json.Ident."}`
	}
	return buf.String()
}

func parse(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	p := parser.New(tokenizer.New(input))
	program := p.ParseProgram()
	parseJson, err := json.Marshal(program)
	if err != nil {
		return `{"error": "error json.Marshal."}`
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, parseJson, "", "\t")
	if err != nil {
		return `{"error": "error json.Ident."}`
	}
	return buf.String()
}

func eval(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	p := parser.New(tokenizer.New(input))
	program := p.ParseProgram()
	env := object.NewEnvironment()
	var buf bytes.Buffer
	evaluator.SetWriter(&buf)
	evaluator.Eval(program, env)

	return buf.String()
}

func main() {
	c := make(chan struct{})
	js.Global().Set("tokenize", js.FuncOf(tokenize))
	js.Global().Set("parse", js.FuncOf(parse))
	js.Global().Set("eval", js.FuncOf(eval))
	<-c
}
