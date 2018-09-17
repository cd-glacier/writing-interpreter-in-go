package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/g-hyoga/writing-interpreter-in-go/src/lexer"
	"github.com/g-hyoga/writing-interpreter-in-go/src/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Println("bye")
			return
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		fmt.Printf("%+v\n", program.String())
	}
}
