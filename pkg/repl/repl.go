package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/oteto/gonkey/pkg/tokenizer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		t := tokenizer.New(line)

		for token := t.NextToken(); token.Type != tokenizer.EOF; token = t.NextToken() {
			fmt.Printf("%+v\n", token)
		}
	}
}
