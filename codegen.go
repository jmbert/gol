package main

import (
	"fmt"
	"log"

	"github.com/jmbert/golib"
)

type Code struct {
	predirs map[string]bool
	code    []string
	data    []string

	strvals map[string]string
}

func (c Code) String() string {
	var ret string
	for predir, valid := range c.predirs {
		if valid {
			ret += fmt.Sprintf(".extern %s\n", predir)
		}
	}
	for _, data := range c.data {
		ret += fmt.Sprintln(data)
	}
	for _, code := range c.code {
		ret += fmt.Sprintln(code)
	}
	return ret
}

type Macro interface {
	Expand(*Code, *Environment, *golib.Tree[ASTCore])
}

type LispMacro struct {
	node golib.Tree[ASTCore]
}

func (m LispMacro) Expand(c *Code, e *Environment, node *golib.Tree[ASTCore]) {
	Codegen(m.node, c, e)
}

type SpecialCaseMacro struct {
	Run func(e *Environment, node *golib.Tree[ASTCore])
}

func (m SpecialCaseMacro) Expand(c *Code, e *Environment, node *golib.Tree[ASTCore]) {
	m.Run(e, node)
}

type ExternalMacro struct {
	name string
}

func (m ExternalMacro) Expand(c *Code, e *Environment, node *golib.Tree[ASTCore]) {
	for _, expr := range node.Children {
		Codegen(*expr, c, e)
	}
	c.predirs[m.name] = true
	c.code = append(c.code, "mov %rsp, %rsi")
	c.code = append(c.code, fmt.Sprintf("mov $%d, %%rdi", len(node.Children)))
	c.code = append(c.code, "call "+m.name)
	c.code = append(c.code, fmt.Sprintf("add $%d, %%rsp", len(node.Children)*8))
	c.code = append(c.code, "push %rax")
}

type Environment map[string]Macro

func Codegen(node golib.Tree[ASTCore], c *Code, e *Environment) {
	switch node.Core.core.id {
	case TOKEN_VALUE:
		if len(node.Children) > 0 {
			macro, there := (*e)[node.Core.core.content]
			if !there {
				log.Fatalln("no such macro " + node.Core.core.content)
			}
			macro.Expand(c, e, &node)
		} else {
			macro, there := (*e)[node.Core.core.content]
			if there {
				macro.Expand(c, e, &node)
				return
			}
			name, there := c.strvals[node.Core.core.content]
			if !there {
				c.strvals[node.Core.core.content] = fmt.Sprintf("str%d", len(c.strvals))
				name = c.strvals[node.Core.core.content]
				c.data = append(c.data, fmt.Sprintf("%s: .asciz \"%s\"", name, node.Core.core.content))
			}
			c.code = append(c.code, fmt.Sprintf("pushq $%s", name))
		}

	case TOKEN_PROGRAM:
		for _, expr := range node.Children {
			Codegen(*expr, c, e)
		}
	}
}
