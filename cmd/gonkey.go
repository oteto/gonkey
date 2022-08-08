package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/oteto/gonkey/pkg/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello %s! This is the Gonkey programing language!\n", user.Username)
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
