package parser

import (
	"fmt"
	"glox/ast"
	"glox/errors"
	"glox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []ast.Stmt {
	statements := []ast.Stmt{}

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() ast.Stmt {
	var stmt ast.Stmt
	var err error

	if p.match(token.VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil
	}

	return stmt
}

func (p *Parser) statement() (ast.Stmt, error) {
    if p.match(token.IF) {
        return p.ifStatement()
    }

    if p.match(token.WHILE) {
        return p.whileStatement()
    }

    if p.match(token.FOR) {
        return p.forStatement()
    }

	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.LEFT_BRACE) {
		block, err := p.block()

		if err != nil {
			return nil, err
		}

		return ast.NewBlock(block), nil
	}

	return p.expressionStatement()
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
    if _, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
        return nil, err
    }

    condition, err := p.expression()

    if err != nil {
        return nil, err
    }

    if _, err := p.consume(token.RIGHT_PAREN, "Expect ')' after if condition."); err != nil {
        return nil, err
    }

    thenBranch, err := p.statement()

    if err != nil {
        return nil, err
    }
    
    var elseBranch ast.Stmt

    if p.match(token.ELSE) {
        b, err := p.statement()
        if err != nil {
            return nil, err
        }

        elseBranch = b
    }

    return ast.NewIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
    _, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")

    if err != nil {
        return nil, err
    }

    condition, err := p.expression()

    if err != nil {
        return nil, err
    }

    _, err = p.consume(token.RIGHT_PAREN, "Expect ')' after while condition.")

    if err != nil {
        return nil, err
    }

    body, err := p.statement()

    if err != nil {
        return nil, err
    }

    return ast.NewWhile(condition, body), nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
    if _, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
        return nil, err
    }

    var initializer ast.Stmt
    if p.match(token.SEMICOLON) {
        initializer = nil
    } else if p.match(token.VAR) {
        i, err := p.varDeclaration()

        if err != nil {
            return nil, err
        }

        initializer = i
    } else {
        i, err := p.expressionStatement()

        if err != nil {
            return nil, err
        }

        initializer = i
    }

    var condition ast.Expr
    
    if !p.check(token.SEMICOLON) {
        c, err := p.expression()

        if err != nil {
            return nil, err
        }

        condition = c
    }

    _, err := p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

    if err != nil {
        return nil, err
    }

    var increment ast.Expr
    if !p.check(token.RIGHT_PAREN) {
        i, err := p.expression()

        if err != nil {
            return nil, err
        }

        increment = i
    }

    p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

    body, err := p.statement()

    if err != nil {
        return nil, err
    }

    if increment != nil {
        body = ast.NewBlock([]ast.Stmt{body, ast.NewExpression(increment)})
    }

    if condition == nil {
        condition = ast.NewLiteral(true)
    }

    body = ast.NewWhile(condition, body)

    if initializer != nil {
        body = ast.NewBlock([]ast.Stmt{initializer, body})
    }

    return body, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(token.SEMICOLON, "Expected ';' after value."); err != nil {
		return nil, err
	}

	return ast.NewPrint(val), nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	if _, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")

	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")

	if err != nil {
		return nil, err
	}

	return ast.NewVar(name, initializer), nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err = p.consume(token.SEMICOLON, "Expected ';' after value."); err != nil {
		return nil, err
	}

	return ast.NewExpression(expr), nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignement()
}

func (p *Parser) assignement() (ast.Expr, error) {
	expr, err := p.or()

	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignement()

		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*ast.Variable); ok {
			name := variable.Name
			return ast.NewAssign(name, value), nil
		}

		errors.Error(equals, "Invalid assignement target.")
	}

	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
    expr, err := p.and()

    if err != nil {
        return nil, err
    }

    for p.match(token.OR) {
        operator := p.previous()
        right, err := p.and()

        if err != nil {
            return nil, err
        }

        expr = ast.NewLogical(expr, operator, right)
    }

    return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
    expr, err := p.equality()

    if err != nil {
        return nil, err
    }

    for p.match(token.AND) {
        operator := p.previous()
        right, err := p.equality()

        if err != nil {
            return nil, err
        }

        expr = ast.NewLogical(expr, operator, right)
    }

    return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()

	if err != nil {
		return nil, err
	}

	for p.match(token.EQUAL_EQUAL, token.BANG_EQUAL) {
		op := p.previous()
		right, err := p.comparison()

		if err != nil {
			return nil, err
		}

		expr = ast.NewBinary(expr, op, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()

	if err != nil {
		return nil, err
	}

	for p.match(token.LESS, token.LESS_EQUAL, token.GREATER, token.GREATER_EQUAL) {
		op := p.previous()
		right, err := p.term()

		if err != nil {
			return nil, err
		}

		expr = ast.NewBinary(expr, op, right)
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()

	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		op := p.previous()
		right, err := p.factor()

		if err != nil {
			return nil, err
		}

		expr = ast.NewBinary(expr, op, right)
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		op := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = ast.NewBinary(expr, op, right)
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		expr, err := p.unary()
		if err != nil {
			return nil, err
		}

		return ast.NewUnary(p.advance(), expr), nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.TRUE) {
		return ast.NewLiteral(true), nil
	}

	if p.match(token.FALSE) {
		return ast.NewLiteral(false), nil
	}

	if p.match(token.STRING, token.NUMBER) {
		return ast.NewLiteral(p.previous().Literal()), nil
	}

	if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")

		return ast.NewGrouping(expr), nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(tokenType token.TokenType, message string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	err := p.error(p.peek(), message)

	return token.Token{}, err
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type() == token.SEMICOLON {
			return
		}

		switch p.previous().Type() {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) error(t token.Token, message string) error {
	errors.Error(t, message)
	return fmt.Errorf("")
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.tokens[p.current].Type() == t
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().Type() == token.EOF
}

func (p Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p Parser) previous() token.Token {
	return p.tokens[p.current-1]
}
