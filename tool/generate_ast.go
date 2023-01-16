package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	outputDir := args[1]

	defineAst(outputDir, "Expr", []string{
		"Assign   : Name token.Token, Value Expr",
		"Binary   : Left Expr, Operator token.Token, Right Expr",
        "Call     : Callee Expr, Paren token.Token, Arguments []Expr",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
        "Logical  : Left Expr, Operator token.Token, Right Expr",
		"Unary    : Operator token.Token, Right Expr",
		"Variable : Name token.Token",
	}, "(interface{}, error)")

	defineAst(outputDir, "Stmt", []string{
		"Block      : Statements []Stmt",
		"Expression : Exp Expr",
        "Function   : Name token.Token, Params []token.Token, Body []Stmt",
        "If         : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print      : Exp Expr",
        "Return     : Keyword token.Token, Value Expr",
		"Var        : Name token.Token, Initializer Expr",
        "While      : Condition Expr, Body Stmt",
	}, "error")
}

func checkError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func defineAst(outputDir string, baseName string, types []string, visitorReturn string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"

	f, err := os.Create(path)
	checkError(err)

	defer f.Close()

	f.WriteString("package ast\n\n")
	f.WriteString("import \"glox/token\"\n\n")
	f.WriteString("type " + baseName + " interface {\n\tAccept(Visitor" + baseName + ") " + visitorReturn + "\n}\n\n")

	defineVisitor(f, baseName, types, visitorReturn)

	for _, t := range types {
		parts := strings.Split(t, ":")

		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])

		defineStruct(f, structName, baseName, fields, visitorReturn)
	}
}

func defineStruct(f *os.File, structName string, baseName string, fieldList string, visitorReturn string) {
	// Struct
	f.WriteString("type " + structName + " struct {\n")

	fields := strings.Split(fieldList, ", ")

	for _, field := range fields {
		f.WriteString("\t" + field + "\n")
	}

	f.WriteString("}\n\n")

	// Constructor

	f.WriteString("func New" + structName + "(" + fieldList + ") *" + structName + " {\n")
	f.WriteString("\t return &" + structName + "{")

	for i, field := range fields {
		fieldName := strings.Split(field, " ")[0]
		f.WriteString(fieldName + ": " + fieldName)

		if i < len(fields)-1 {
			f.WriteString(", ")
		}
	}

	f.WriteString("}\n}\n\n")

	f.WriteString("func (e *" + structName + ") Accept(v Visitor" + baseName + ") " + visitorReturn + " {\n")
	f.WriteString("\treturn v.Visit" + structName + baseName + "(e)\n")
	f.WriteString("}\n\n")

	f.WriteString("\n")
}

func defineVisitor(f *os.File, baseName string, types []string, visitorReturn string) {
	f.WriteString("type Visitor" + baseName + " interface {\n")

	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])

		f.WriteString("\tVisit" + typeName + baseName + "(" + strings.ToLower(baseName) + " *" + typeName + ") " + visitorReturn + "\n")
	}

	f.WriteString("}\n\n")
}
