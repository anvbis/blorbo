package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"blorbo/pkg/lexer"
	"blorbo/pkg/parser"
)

const version = "0.1.0"

func runScript(src string) error {
	l := lexer.New(src)
	tokens, err := l.Scan()
	if err != nil {
		return errors.New("error: lexer error(s) occurred")
	}

	p := parser.New(tokens)
	stmt, err := p.Parse()
	if err != nil {
		return err
	}

	tree, _ := json.MarshalIndent(stmt, "", "  ")
	fmt.Println(string(tree))

	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Blorbo %s\n", version)

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			src, _ := reader.ReadString('\n')

			if err := runScript(src); err != nil {
				fmt.Println(err)
			}
		}
	} else if len(os.Args) == 2 {
		src, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := runScript(string(src)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("usage: blorbo [script]")
		os.Exit(1)
	}
}
