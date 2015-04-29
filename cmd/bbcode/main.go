package main

import (
	"bufio"
	"fmt"
	"github.com/frustra/bbcode"
	"os"
)

func main() {
	compiler := bbcode.NewCompiler(true, true)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(compiler.Compile(scanner.Text()))
	}
}
