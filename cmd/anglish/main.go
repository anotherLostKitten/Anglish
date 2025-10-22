package main

import (
	"fmt"
	// "io"
	"strings"

	"github.com/anotherLostKitten/Anglish/internal/parse"
)

func main() {
	r := strings.NewReader("@spacey:CALL ( in=%argy )\n> i'm vibing %big time 100%\n>\n>$tasky(in=%a, out=%b)")
	c, errors := parse.ParseFromReader(r)
	fmt.Printf("error number : %d\n", len(errors))
	for i := 0; i < len(errors); i++ {
		parse.PrintErrorInfo(errors[i])
	}
	fmt.Printf("%+v\n", c)
}
