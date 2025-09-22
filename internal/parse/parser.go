package parse

import (
	"fmt"
	//	"io"
	"strings"
)

type ParserInfo struct {
	line, col uint64

	// todo errorset
}

func ParseFromReader(reader *strings.Reader) {
	pi := ParserInfo{
		line: 0,
		col: 0,
	}
	for reader.Len() > 0 {
		consumeBlankLines(reader, &pi)
		ch, size, _ := reader.ReadRune()
		if size != 1 {
			// todo error
			break
		}
		switch ch {
		case '@':
			reader.UnreadRune()
			parseSpaceDecl(reader, &pi)
		default:
			// todo error
			break
		}
	}
}

func consumeBlankLines(reader *strings.Reader, pi *ParserInfo) {
	for reader.Len() > 0 {
		ch, size, _ := reader.ReadRune()
		if size != 1 {
			reader.UnreadRune()
			return
		}
		switch ch {
		case '\n':
			pi.line++
			pi.col = 0
		case ' ', '\t':
			pi.col++
		default:
			reader.UnreadRune()
			return
		}
	}
}

func parseSpaceDecl(reader *strings.Reader, pi *ParserInfo) {
	fmt.Println("ok guys i guess it is time to parse this struct over here !")
}
