package main

import (
	"bufio"
	"fmt"
	"glox/errors"
	"glox/interpreter"
	"glox/parser"
	"glox/resolver"
	"glox/scanner"
	"io"
	"os"
)

var interp = interpreter.NewInterpreter()

func main() {
    args := os.Args

    if len(args) > 2 {
        fmt.Println("Usage: glox [script]")
        os.Exit(64)
    } else if len(args) == 2 {
        err := runFile(args[1])
        if err != nil {
            panic(err)
        }
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

    if errors.HadRuntimeError {
        os.Exit(70)
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
    statements := p.Parse()

    if errors.HadError {
        return
    }

    res := resolver.NewResolver(&interp)
    res.Resolve(statements)

    if errors.HadError {
        return
    }

    interp.Interpret(statements)
}
