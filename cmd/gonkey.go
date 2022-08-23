package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/oteto/gonkey/pkg/repl"
)

var (
	tokenizerOpt = flag.Bool("t", false, "help message for \"t\" option")
	parserOpt    = flag.Bool("p", false, "help message for \"p\" option")
	evalOpt      = flag.Bool("e", false, "help message for \"e\" option")
)

func main() {
	flag.Parse()
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello %s! This is the Gonkey programing language!\n", user.Username)
	fmt.Println("Feel free to type in commands")

	if *tokenizerOpt {
		fmt.Println("output Token.")
		repl.TokenizerStart(os.Stdin, os.Stdout)
	} else if *parserOpt {
		fmt.Println("output AST.")
		repl.ParserStart(os.Stdin, os.Stdout)
	} else if *evalOpt {
		fmt.Println("output Eval.")
		repl.EvalStart(os.Stdin, os.Stdout)
	} else {
		fmt.Println("please input option -t or -p.")
	}
}
