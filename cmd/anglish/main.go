package main

import (
	"fmt"
	// "io"
	"os"

	"github.com/anotherLostKitten/Anglish/internal/compile"
	"github.com/anotherLostKitten/Anglish/internal/parse"
)

func main() {
	fn := "examples/chatbot.ang"
	if len(os.Args) == 2 {
		fn = os.Args[1]
	} else {
		fmt.Printf("No filename specified, using default: `%s`\n", fn)
	}

	c, errors := parse.ParseFromFile(fn)

	fmt.Printf("error number : %d\n", len(errors))
	for _, err := range errors {
		parse.PrintErrorInfo(err)
	}
	fmt.Printf("%+v\n", c)

	po := parse.GetParseOrder(&c)

	_ = compile.Compile(&po)
}
