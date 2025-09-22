package parse

import (
	// "fmt"
	// "io"
	"strings"
)

type ParserInfo struct {
	line, col uint64

	errors []ParserErrorInfo
}

type locationTaggedString struct {
	val string
	line, col uint64
}

func ParseFromReader(reader *strings.Reader) (Contract, []ParserErrorInfo) {
	pi := ParserInfo{
		line: 0,
		col: 0,
	}

	var c Contract
	for reader.Len() > 0 {
		consumeBlankLines(reader, &pi)
		ch, size, _ := reader.ReadRune()
		if size != 1 {
			pi.addError(NonAsciiChar)
			pi.col += uint64(size)
		}
		switch ch {
		case '@':
			reader.UnreadRune()
			spacey := parseSpaceDecl(reader, &pi)
			if spacey != nil {
				c.spaces = append(c.spaces, *spacey)
			}
		default:
			pi.addError(ExpectedOuterDecl)
			pi.col++
		}
	}

	return c, pi.errors
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

func consumeLineRemainder(reader *strings.Reader, pi *ParserInfo) {
	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if ch == '\n' {
			pi.col = 0
			pi.line++
			return
		}
	}
}

func consumeSpaces(reader *strings.Reader, pi *ParserInfo) {
	for reader.Len() > 0 {
		ch, size, _ := reader.ReadRune()
		if size != 1 {
			reader.UnreadRune()
			return
		}
		switch ch {
		case ' ', '\t':
			pi.col++
		default:
			reader.UnreadRune()
			return
		}
	}
}

func identStart(ch rune) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func identPart(ch rune) bool {
	return identStart(ch) || ch >= '0' && ch <= '9'
}

func parseIdentifier(reader *strings.Reader, pi *ParserInfo) string {
	var ident strings.Builder
	ch, size, _ := reader.ReadRune()
	if identStart(ch) {
		ident.WriteRune(ch)
		pi.col += uint64(size)
	} else {
		reader.UnreadRune()
		return ident.String()
	}

	for reader.Len() > 0 {
		ch, size, _ := reader.ReadRune()
		if identPart(ch) {
			ident.WriteRune(ch)
			pi.col += uint64(size)
		} else {
			reader.UnreadRune()
			return ident.String()
		}
	}
	return ident.String()
}

func parseTags(reader *strings.Reader, pi *ParserInfo) []locationTaggedString {
	var tags []locationTaggedString

	for reader.Len() > 0 {
		consumeSpaces(reader, pi)

		ch, size, _ := reader.ReadRune()
		if ch != ':' {
			reader.UnreadRune()
			return tags
		}
		pi.col += uint64(size)

		consumeSpaces(reader, pi)

		old_col := pi.col
		ident := parseIdentifier(reader, pi)
		if ident == "" {
			pi.addError(ExpectedIdentifier)
			return tags
		}
		tags = append(tags, locationTaggedString{
			val: strings.ToUpper(ident),
			line: pi.line,
			col: old_col,
		})
	}
	return tags
}

func parseParams(reader *strings.Reader, pi *ParserInfo) []Param {
	consumeSpaces(reader, pi)

	var params []Param

	// todo -- if we always return when we don't have an opening param we don't need this
	has_open_paren := false
	ch, size, _ := reader.ReadRune()
	if ch == '(' {
		pi.col += uint64(size)
		has_open_paren = true
		consumeSpaces(reader, pi)
	} else {
		reader.UnreadRune()
		return params
	}

	last_param_in := true
	for reader.Len() > 0 {
		ch, size, _ := reader.ReadRune()
		if ch == '\n' || ch == ')' {
			reader.UnreadRune()
			break
		}

		var p Param
		p.line = pi.line
		p.col = pi.col
		if ch != '%' {
			reader.UnreadRune()
			in_out := parseIdentifier(reader, pi)
			if in_out == "" {
				pi.addError(ExpectedIdentifier)
				return params
			}

			expected_in_out_error := false

			switch strings.ToLower(in_out) {
			case "in":
				last_param_in = true
				p.in_param = true
			case "out":
				last_param_in = false
				p.in_param = false
			default:
				expected_in_out_error = true
				p.in_param = last_param_in
			}
			consumeSpaces(reader, pi)
			ch, size, _ := reader.ReadRune()
			if ch != '=' {
				reader.UnreadRune()

				if expected_in_out_error {
					pi.addError(ExpectedDataName)
					p.data_name = in_out
					params = append(params, p)
					continue
				} else {
					pi.addError(ExpectedEquals)
				}
			} else {
				pi.col += uint64(size)
				consumeSpaces(reader, pi)

				if expected_in_out_error {
					pi.addError(ExpectedInOut)
				}
			}

			ch, size, _ = reader.ReadRune()
			if ch != '%' {
				reader.UnreadRune()
				pi.addError(ExpectedDataName)
			} else {
				pi.col += uint64(size)
			}
		} else {
			p.in_param = last_param_in
			pi.col += uint64(size)
		}
		p.data_name = parseIdentifier(reader, pi)
		if p.data_name == "" {
			pi.addError(ExpectedIdentifier)
			return params
		}

		consumeSpaces(reader, pi)

		params = append(params, p)

		ch, size, _ = reader.ReadRune()
		if ch == ',' || ch == ';' {
			pi.col += uint64(size)
			consumeSpaces(reader, pi)
		} else {
			reader.UnreadRune()
		}
	}

	ch, size, _ = reader.ReadRune()
	if ch == ')' {
		if !has_open_paren {
			pi.addError(MismatchedParens)
		}
		pi.col += uint64(size)
	} else {
		if has_open_paren {
			pi.addError(MismatchedParens)
		}
		reader.UnreadRune()
	}
	return params
}

func parseSpaceDecl(reader *strings.Reader, pi *ParserInfo) *SpaceDecl {
	ch, size, _ := reader.ReadRune()
	if ch != '@' {
		reader.UnreadRune()
		pi.addError(ExpectedSpaceDecl)
		return nil
	}
	pi.col += uint64(size)

	var decl SpaceDecl
	decl.line_start = pi.line

	decl.ident = parseIdentifier(reader, pi)
	if decl.ident == "" {
		pi.addError(ExpectedIdentifier)
		consumeLineRemainder(reader, pi)
		return nil
	}

	tags := parseTags(reader, pi)
	for i := 0; i < len(tags); i++ {
		switch tags[i].val {
		case "REPLICABLE":
			if decl.replicable == true {
				pi.addErrorTagged(DuplicateTag, tags[i])
			}
			decl.replicable = true
		case "UI":
			if decl.space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.space_type = UI
			}
		case "IO":
			if decl.space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.space_type = IO
			}
		case "DATA":
			if decl.space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.space_type = DATA
			}
		case "FUNC":
			if decl.space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.space_type = FUNC
			}
		case "AGENTIC":
			if decl.space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.space_type = AGENTIC
			}
		default:
			pi.addErrorTagged(UnknownTag, tags[i])
		}
	}

	decl.params = parseParams(reader, pi)

	consumeLineRemainder(reader, pi)

	// todo parse vibe blocks
	// todo parse rest of scope

	decl.line_end = pi.line
	return &decl
}
