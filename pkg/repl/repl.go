package repl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/oteto/gonkey/pkg/evaluator"
	"github.com/oteto/gonkey/pkg/object"
	"github.com/oteto/gonkey/pkg/parser"
	"github.com/oteto/gonkey/pkg/token"
	"github.com/oteto/gonkey/pkg/tokenizer"
)

const PROMPT = ">> "

func TokenizerStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		t := tokenizer.New(line)

		for tkn := t.NextToken(); tkn.Type != token.EOF; tkn = t.NextToken() {
			fmt.Printf("%+v\n", tkn)
		}
	}
}

func ParserStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		p := parser.New(tokenizer.New(line))
		program := p.ParseProgram()
		programJson, err := json.Marshal(program)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		var buf bytes.Buffer
		err = json.Indent(&buf, programJson, "", "  ")
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		fmt.Println(buf.String())

	}
}

func printParserError(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func EvalStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		p := parser.New(tokenizer.New(line))
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserError(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
