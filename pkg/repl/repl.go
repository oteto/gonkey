package repl

import (
	"bufio"
	"fmt"
	"io"

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
		t := tokenizer.New(line)
		p := parser.New(t)
		program := p.ParseProgram()

		for _, s := range program.Statements {
			fmt.Printf("%+v\n", s.String())
		}
	}
}
