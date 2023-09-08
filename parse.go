package main

import (
	"github.com/jmbert/golib"
)

type ASTCore struct {
	core Token
}

func (a ASTCore) String() string {
	return a.core.String()
}

func Parse(tokens []Token) (golib.Tree[ASTCore], int, error) {
	var node golib.Tree[ASTCore]

	current_token := tokens[0]
	var read = 1

	switch current_token.id {
	case TOKEN_PROGRAM:
		node.Core.core = current_token
		for i := read; i < len(tokens); {
			child, size, err := Parse(tokens[(i):])
			if err != nil {
				return node, read, err
			}
			i += size
			read += size
			node.Children = append(node.Children, &child)
		}
	case TOKEN_EXPRESSION:
		node.Core.core = tokens[1]
		read = 2
		for i := read; i < len(tokens); {
			if tokens[i].id == TOKEN_CLOSEPAREN {
				read++
				break
			} else {
				child, size, err := Parse(tokens[(i):])
				if err != nil {
					return node, read, err
				}
				i += size
				read += size
				node.Children = append(node.Children, &child)
			}
		}
	case TOKEN_VALUE:
		node.Core.core = current_token
	}

	return node, read, nil
}
