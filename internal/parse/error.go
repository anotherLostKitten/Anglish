package parse

import (
	"fmt"
)

type ParserError uint64
const (
	UnexpectedMetachar ParserError = iota
	NonAsciiChar
	ExpectedOuterDecl
	ExpectedInnerDecl
	ExpectedSpaceDecl
	ExpectedAgentDecl
	ExpectedTaskDecl
	ExpectedPathDecl
	ExpectedDataName
	ExpectedIdentifier
	ExpectedInOut
	ExpectedEquals
	MissingRequiredTag
	DuplicateTag
	UnknownTag
	MismatchedParens
	UseMissingImport
	UseUnsupportedImport
	IllegalDeclarationInsideSpaceScope
)

type ParserErrorInfo struct {
	err ParserError
	line, col uint64
}

func (pi *ParserInfo) addError(errno ParserError) {
	pi.errors = append(pi.errors, ParserErrorInfo{
		err: errno,
		line: pi.line,
		col: pi.col,
	})
}

func (pi *ParserInfo) addErrorTagged(errno ParserError, location locationTaggedString) {
	pi.errors = append(pi.errors, ParserErrorInfo{
		err: errno,
		line: location.line,
		col: location.col,
	})
}

func PrintErrorInfo(errinf ParserErrorInfo) {
	fmt.Printf("Error at (%d, %d): ", errinf.line, errinf.col)
	switch errinf.err {
	case UnexpectedMetachar: fmt.Printf("Unexpected meta-character")
	case NonAsciiChar: fmt.Printf("Unexpected non-ASCII character")
	case ExpectedOuterDecl: fmt.Printf("Expected Declaration: @space, #agent, $task, =path")
	case ExpectedInnerDecl: fmt.Printf("Expected Declaration inside @space scope: #agent, $task")
	case ExpectedSpaceDecl: fmt.Printf("Expected Space Declaration: @space")
	case ExpectedAgentDecl: fmt.Printf("Expected Agent Declaration: #agent")
	case ExpectedTaskDecl: fmt.Printf("Expected Task Declaration: $task")
	case ExpectedPathDecl: fmt.Printf("Expected Path Declaration: =path")
	case ExpectedDataName: fmt.Printf("Expected Data Name: %%data")
	case ExpectedIdentifier: fmt.Printf("Expected Identifier: ident")
	case ExpectedInOut: fmt.Printf("Expected in or out")
	case ExpectedEquals: fmt.Printf("Expected =")
	case MissingRequiredTag: fmt.Printf("Missing Required Tag Definition")
	case DuplicateTag: fmt.Printf("Duplicate or Contradictory Tag Definition")
	case UnknownTag: fmt.Printf("Unknown Tag Name")
	case MismatchedParens: fmt.Printf("Mismatched Parentheses")
	case UseMissingImport: fmt.Printf("Missing import for $use expression: should take the form $use(element), where element is a @space or #agent.")
	case UseUnsupportedImport: fmt.Printf("Cannot import this element. Expression should take the form $use(element), where element is a @space or #agent.")
	case IllegalDeclarationInsideSpaceScope: fmt.Printf("Illegal declaration inside @space scope. Should be: #agent, $task")
	default: fmt.Printf("???")
	}
	fmt.Printf("\n")

}
