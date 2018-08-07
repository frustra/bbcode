package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/frustra/bbcode"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	compiler := bbcode.NewCompiler(true, true)
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		//stdin is from a pipe. Compile it all at once.
		b, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(compiler.Compile(string(b)))
	} else {
		//stdin isn't from a pipe; compile it line-by-line.
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(compiler.Compile(scanner.Text()))
		}
	}
}
