package main

import (
	"fmt"
	// "io"
	"strings"

	"github.com/anotherLostKitten/Anglish/internal/parse"
)

func main() {
	r := strings.NewReader("@spacey:FUNC ( in=%argy )")
	parse.ParseFromReader(r)
	fmt.Println("Hello World!")
}
