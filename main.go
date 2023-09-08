package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	input := get_input()
	source, err := io.ReadAll(input.source)
	if err != nil {
		log.Fatalln(err)
	}

	tokens, err := Lex(string(source))
	if err != nil {
		log.Fatalln(err)
	}

	ast, _, err := Parse(tokens)
	if err != nil {
		log.Fatalln(err)
	}

	var code Code

	code.data = append(code.data, ".data")
	code.code = append(code.code, ".text", ".globl _start", "_start:")
	code.predirs = make(map[string]bool)

	code.strvals = make(map[string]string)

	env := make(Environment)

	initEnv(&env)

	Codegen(ast, &code, &env)

	code.code = append(code.code, "mov $0x3C, %rax")
	code.code = append(code.code, "pop %rdi")
	code.code = append(code.code, "syscall")

	assemble := exec.Command("/usr/bin/as", "-o/tmp/golc.out.o")

	assemble.Stdin = strings.NewReader(code.String())

	assemble.Stdout = os.Stdout
	assemble.Stderr = os.Stderr

	err = assemble.Run()
	if err != nil {
		log.Fatalf("Assembler: %v\n", err)
	}

	link := exec.Command("/usr/bin/ld", "--dynamic-linker", "/lib64/ld-linux-x86-64.so.2", "/tmp/golc.out.o", "-lgolstd", "-o"+input.out)
	link.Stdout = os.Stdout
	link.Stderr = os.Stderr

	err = link.Run()
	if err != nil {
		log.Fatalf("Linker:%v\n", err)
	}

	if err != nil {
		exiterr, _ := err.(*exec.ExitError)
		log.Fatalln(exiterr)

	}
}
