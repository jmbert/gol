package main

import "github.com/jmbert/golib"

func define_f(e *Environment, node *golib.Tree[ASTCore]) {
	newMacro := node.Children[0].Core.core.content
	newReplace := node.Children[1]
	(*e)[newMacro] = LispMacro{*newReplace}
}

var defineMacro = SpecialCaseMacro{define_f}

// External macros
var printMacro = ExternalMacro{"_print_lisp"}

var addMacro = ExternalMacro{"_add_lisp"}
var subMacro = ExternalMacro{"_sub_lisp"}
var mulMacro = ExternalMacro{"_mul_lisp"}
var divMacro = ExternalMacro{"_div_lisp"}

func initEnv(e *Environment) {
	(*e)["define"] = defineMacro
	(*e)["+"] = addMacro
	(*e)["-"] = subMacro
	(*e)["*"] = mulMacro
	(*e)["/"] = divMacro
	(*e)["print"] = printMacro
}
