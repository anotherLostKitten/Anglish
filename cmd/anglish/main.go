package main

import (
	"fmt"
	// "io"
	"strings"

	"github.com/anotherLostKitten/Anglish/internal/parse"
)

func main() {
	r := strings.NewReader("@spacey:FUNC ( in=%argy )")
	c, errors := parse.ParseFromReader(r)
	fmt.Printf("error number : %d\n", len(errors))
	fmt.Printf("%+v\n", c)
	fmt.Println("Hello World!")
}
