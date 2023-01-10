package main

import (
	"bufio"
	"fmt"
	"glox/errors"
	"glox/parser"
	"glox/scanner"
	"glox/tree"
	"io"
	"os"
)

func main() {
    args := os.Args

    if len(args) > 2 {
        fmt.Println("Usage: glox [script]")
        os.Exit(64)
    } else if len(args) == 2 {
        runFile(args[1])
    } else {
        runPrompt()
    }
}

func runFile(path string) error {
    b, err := os.ReadFile(path)

    if err != nil {
        return err
    }

    run(string(b))

    if errors.HadError {
        os.Exit(65)
    }

    return nil
}

func runPrompt() error {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("> ")
        line, err := reader.ReadString('\n')

        if err == io.EOF {
            break
        }

        if err != nil {
            return err
        }

        run(line)

        errors.HadError = false
    }
    
    fmt.Print("\n")
    return nil
}

func run(source string) {
    s := scanner.NewScanner(source)
    tokens := s.ScanTokens()

    p := parser.NewParser(tokens)
    expr := p.Parse()

    if errors.HadError {
        return
    }

    fmt.Println(tree.AstPrinter{}.Print(expr))

}

