package parse

import (
	// "fmt"
	"bufio"
	"io"
	"os"
	"strings"
)

type ParserInfo struct {
	line, col uint64

	errors []ParserErrorInfo
}

type locationTaggedString struct {
	val       string
	line, col uint64
}

func ParseFromFile(fn string) (Contract, []ParserErrorInfo) {
	f, error := os.Open(fn)
	if error != nil {
		return Contract{}, []ParserErrorInfo{{
			err: FileNotFound,
		}}
	}
	defer f.Close()
	r := bufio.NewReader(f)

	return ParseFromReader(r)
}

func ParseFromString(in string) (Contract, []ParserErrorInfo) {
	r := strings.NewReader(in)
	return ParseFromReader(r)
}

func ParseFromReader(reader io.RuneScanner) (Contract, []ParserErrorInfo) {
	pi := ParserInfo{
		line: 0,
		col:  0,
	}

	var c Contract
	for {
		consumeBlankLines(reader, &pi)
		ch, size, _ := reader.ReadRune()
		if size == 0 {
			break
		} else if size > 1 {
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
		case '#':
			reader.UnreadRune()
			ref := parseAgentDecl(reader, &pi)
			if ref != nil {
				c.agents = append(c.agents, *ref)
			}
		case '=':
			reader.UnreadRune()
			ref := parsePathDecl(reader, &pi)
			if ref != nil {
				c.paths = append(c.paths, *ref)
			}
		default:
			pi.addError(ExpectedOuterDecl)
			pi.col++
		}
	}

	return c, pi.errors
}

func consumeBlankLines(reader io.RuneScanner, pi *ParserInfo) {
	for {
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

func consumeLineRemainder(reader io.RuneScanner, pi *ParserInfo) {
	for {
		ch, size, _ := reader.ReadRune()
		if size == 0 {
			return
		}
		if ch == '\n' {
			pi.col = 0
			pi.line++
			return
		}
	}
}

func consumeSpaces(reader io.RuneScanner, pi *ParserInfo) {
	for {
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

func tryParseRune(reader io.RuneScanner, pi *ParserInfo, ch_goal ...rune) bool {
	ch, size, _ := reader.ReadRune()
	if size == 0 {
		return false
	}
	for i := 0; i < len(ch_goal); i++ {
		if ch == ch_goal[i] {
			pi.col += uint64(size)
			return true
		}
	}
	reader.UnreadRune()
	return false
}

func identStart(ch rune) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func identPart(ch rune) bool {
	return identStart(ch) || ch >= '0' && ch <= '9'
}

func parseIdentifier(reader io.RuneScanner, pi *ParserInfo) string {
	var ident strings.Builder
	ch, size, _ := reader.ReadRune()
	if size > 0 && identStart(ch) {
		ident.WriteRune(ch)
		pi.col += uint64(size)
	} else {
		reader.UnreadRune()
		return ident.String()
	}

	for {
		ch, size, _ := reader.ReadRune()
		if size > 0 && identPart(ch) {
			ident.WriteRune(ch)
			pi.col += uint64(size)
		} else {
			reader.UnreadRune()
			return ident.String()
		}
	}
	return ident.String()
}

func parseTags(reader io.RuneScanner, pi *ParserInfo) []locationTaggedString {
	var tags []locationTaggedString

	for {
		consumeSpaces(reader, pi)

		if !tryParseRune(reader, pi, ':') {
			return tags
		}

		consumeSpaces(reader, pi)

		old_col := pi.col
		ident := parseIdentifier(reader, pi)
		if ident == "" {
			pi.addError(ExpectedIdentifier)
			return tags
		}
		tags = append(tags, locationTaggedString{
			val:  strings.ToUpper(ident),
			line: pi.line,
			col:  old_col,
		})
	}
	return tags
}

func parseParams(reader io.RuneScanner, pi *ParserInfo) []Param {
	consumeSpaces(reader, pi)

	var params []Param

	// todo -- if we always return when we don't have an opening param we don't need this
	has_open_paren := false

	if tryParseRune(reader, pi, '(') {
		has_open_paren = true
		consumeSpaces(reader, pi)
	} else {
		return params
	}

	last_param_in := true
	for {
		ch, size, _ := reader.ReadRune()
		if size == 0 {
			break
		}
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

			if tryParseRune(reader, pi, '=') {
				consumeSpaces(reader, pi)
				if expected_in_out_error {
					pi.addError(ExpectedInOut)
				}
			} else {
				if expected_in_out_error {
					pi.addError(ExpectedDataName)
					p.data_name = in_out
					params = append(params, p)
					continue
				} else {
					pi.addError(ExpectedEquals)
				}
			}

			if !tryParseRune(reader, pi, '%') {
				pi.addError(ExpectedDataName)
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

		if tryParseRune(reader, pi, ',', ';') {
			consumeSpaces(reader, pi)
		}
	}

	if tryParseRune(reader, pi, ')') != has_open_paren {
		pi.addError(MismatchedParens)
	}

	return params
}

func parseSpaceDecl(reader io.RuneScanner, pi *ParserInfo) *SpaceDecl {
	if !tryParseRune(reader, pi, '@') {
		pi.addError(ExpectedSpaceDecl)
		return nil
	}

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
			if decl.Space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.Space_type = UI
			}
		case "IO":
			if decl.Space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.Space_type = IO
			}
		case "DATA":
			if decl.Space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.Space_type = DATA
			}
		case "CALL":
			if decl.Space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.Space_type = CALL
			}
		case "CHAT":
			if decl.Space_type != UnknownSpace {
				pi.addErrorTagged(DuplicateTag, tags[i])
			} else {
				decl.Space_type = CHAT
			}
		default:
			pi.addErrorTagged(UnknownTag, tags[i])
		}
	}

	decl.params = parseParams(reader, pi)

	consumeLineRemainder(reader, pi)

	decl.vibe_desc = parseVibeBlock(reader, pi)

InnerDeclLoop:
	for {
		consumeSpaces(reader, pi)

		ch, size, _ := reader.ReadRune()
		if size == 0 {
			break InnerDeclLoop
		}
		reader.UnreadRune()
		switch ch {
		case '#':
			ref := parseAgentDecl(reader, pi)
			if ref != nil {
				decl.Agents = append(decl.Agents, *ref)
			}
		case '$':
			ref := parseTaskDecl(reader, pi)
			if ref != nil {
				decl.Tasks = append(decl.Tasks, *ref)
			}
		case '@':
			pi.addError(IllegalDeclarationInsideSpaceScope)
			break InnerDeclLoop
		case '=':
			pi.addError(IllegalDeclarationInsideSpaceScope)
			break InnerDeclLoop
		case '\n':
			break InnerDeclLoop
		case '\r':
			if !tryParseRune(reader, pi, '\n') {
				pi.addError(ExpectedInnerDecl)
			}
			break InnerDeclLoop
		default:
			pi.addError(ExpectedInnerDecl)
			break InnerDeclLoop
		}
	}

	decl.line_end = pi.line
	return &decl
}

func parseAgentDecl(reader io.RuneScanner, pi *ParserInfo) *AgentDecl {
	if !tryParseRune(reader, pi, '#') {
		pi.addError(ExpectedAgentDecl)
		return nil
	}

	var agent AgentDecl
	agent.line_start = pi.line

	agent.ident = parseIdentifier(reader, pi)
	if agent.ident == "" {
		pi.addError(ExpectedIdentifier)
		return nil
	}

	consumeSpaces(reader, pi)

	tags := parseTags(reader, pi)

	for i := 0; i < len(tags); i++ {
		switch strings.ToUpper(tags[i].val) {
		case "DF":
			if agent.Agent_type != UnknownAgent {
				pi.addErrorTagged(DuplicateTag, tags[i])
			}
			agent.Agent_type = DF
		case "AF":
			if agent.Agent_type != UnknownAgent {
				pi.addErrorTagged(DuplicateTag, tags[i])
			}
			agent.Agent_type = AF
		default:
			pi.addErrorTagged(UnknownTag, tags[i])
		}
	}
	if agent.Agent_type == UnknownAgent {
		pi.addError(MissingRequiredTag)
	}

	agent.params = parseParams(reader, pi)

	consumeLineRemainder(reader, pi)

	agent.vibe_desc = parseVibeBlock(reader, pi)

	agent.line_end = pi.line
	return &agent
}

func parseTaskDecl(reader io.RuneScanner, pi *ParserInfo) *TaskDecl {
	if !tryParseRune(reader, pi, '$') {
		pi.addError(ExpectedTaskDecl)
		return nil
	}

	var task TaskDecl
	task.line_start = pi.line

	task.ident = parseIdentifier(reader, pi)
	if task.ident == "" {
		pi.addError(ExpectedIdentifier)
		return nil
	}

	// todo check if identifier is reserved -- viz., $use

	task.Params = parseParams(reader, pi)

	consumeLineRemainder(reader, pi)

	task.Vibe_desc = parseVibeBlock(reader, pi)

	task.line_end = pi.line
	return &task
}

func parseSpaceParams(reader io.RuneScanner, pi *ParserInfo) []locationTaggedString {
	consumeSpaces(reader, pi)

	var spaces []locationTaggedString

	if !tryParseRune(reader, pi, '(') {
		return spaces
	}

	for {
		consumeSpaces(reader, pi)

		if tryParseRune(reader, pi, ')') {
			return spaces
		}

		if !tryParseRune(reader, pi, '@') {
			pi.addError(ExpectedSpaceName)
			break
		}

		next_space := locationTaggedString{
			line: pi.line,
			col:  pi.col,
		}
		next_space.val = parseIdentifier(reader, pi)

		if next_space.val == "" {
			pi.addError(ExpectedSpaceName)
		} else {
			spaces = append(spaces, next_space)
			consumeSpaces(reader, pi)
		}

		if tryParseRune(reader, pi, ',', ';') {
			consumeSpaces(reader, pi)
		}
	}

	if !tryParseRune(reader, pi, ')') {
		pi.addError(MismatchedParens)
	}

	return spaces
}

func parsePathDecl(reader io.RuneScanner, pi *ParserInfo) *PathDecl {
	if !tryParseRune(reader, pi, '=') {
		pi.addError(ExpectedPathDecl)
		return nil
	}

	var path PathDecl
	path.line_start = pi.line

	path.ident = parseIdentifier(reader, pi)
	if path.ident == "" {
		pi.addError(ExpectedIdentifier)
		return nil
	}

	tags := parseTags(reader, pi)

	for i := 0; i < len(tags); i++ {
		switch strings.ToUpper(tags[i].val) {
		case "INVOKE":
			if path.Path_type != UnknownPath {
				pi.addErrorTagged(DuplicateTag, tags[i])
			}
			path.Path_type = INVOKE
		case "ATTEND":
			if path.Path_type != UnknownPath {
				pi.addErrorTagged(DuplicateTag, tags[i])
			}
			path.Path_type = ATTEND
		default:
			pi.addErrorTagged(UnknownTag, tags[i])
		}
	}
	if path.Path_type == UnknownPath {
		pi.addError(MissingRequiredTag)
	}

	path_spaces := parseSpaceParams(reader, pi)

	if len(path_spaces) != 2 {
		pi.addError(IncorrectNumberPathSpaces)
		return nil
	}

	path.space_source = Ident{
		t: SPACE,
		N: path_spaces[0].val,
	}
	path.space_dest = Ident{
		t: SPACE,
		N: path_spaces[1].val,
	}

	consumeLineRemainder(reader, pi)

	path.vibe_desc = parseVibeBlock(reader, pi)

	path.line_end = pi.line
	return &path
}

func parseVibeBlock(reader io.RuneScanner, pi *ParserInfo) VibeBlock {
	var vb VibeBlock
	vb.line_start = pi.line
BlockLoop:
	for {
		consumeSpaces(reader, pi)

		if !tryParseRune(reader, pi, '>') {
			break BlockLoop
		}

		consumeSpaces(reader, pi)

		if tryParseRune(reader, pi, '\n') {
			pi.line++
			pi.col = 0
			continue BlockLoop
		}

		var vl strings.Builder
	LineLoop:
		for {
			ch, size, _ := reader.ReadRune()
			if size == 0 {
				break LineLoop
			}
			pi.col += uint64(size)
			switch ch {
			case ' ', '\t': // normalize any amount of whitespace into a single space
				consumeSpaces(reader, pi)
				ch2, _, _ := reader.ReadRune()
				if ch2 != '\n' { // trims trailing whitespaces
					vl.WriteRune(' ')
				}
				reader.UnreadRune()
			case '\n':
				pi.line++
				pi.col = 0
				break LineLoop
			case '%':
				ref := parseMetaRefData(reader, pi)
				if ref != nil {
					vl.WriteString(ref.ToStr())
					vb.meta_refs = append(vb.meta_refs, ref)
				} else {
					vl.WriteRune(ch)
				}
			case '$':
				ref := parseMetaRefTask(reader, pi)
				if ref != nil {
					vl.WriteString(ref.ToStr())
					vb.meta_refs = append(vb.meta_refs, ref)
				} else {
					vl.WriteRune(ch)
				}

				// because we can consume spaces while parsing the task params in a non-hygenic way,
				// we make sure there is a space afterwards before remainder of vibe line
				consumeSpaces(reader, pi)
				ch2, _, _ := reader.ReadRune()
				reader.UnreadRune()
				if ch2 != '\n' { // trims trailing whitespaces
					vl.WriteRune(' ')
				}
			case '=':
				ref := parseMetaRefData(reader, pi)
				if ref != nil {
					vl.WriteString(ref.ToStr())
					vb.meta_refs = append(vb.meta_refs, ref)
				} else {
					vl.WriteRune(ch)
				}
			default:
				vl.WriteRune(ch)
			}
		}

		vb.Vibe_prose = append(vb.Vibe_prose, vl.String())
	}
	vb.line_end = pi.line
	return vb
}

func parseMetaRefData(reader io.RuneScanner, pi *ParserInfo) *MetaRefData {
	tryParseRune(reader, pi, '%')

	var mrd MetaRefData
	mrd.line = pi.line
	mrd.col = pi.col

	mrd.ident = parseIdentifier(reader, pi)
	if mrd.ident == "" {
		return nil
	}
	return &mrd
}

func parseMetaRefTask(reader io.RuneScanner, pi *ParserInfo) MetaRef {
	tryParseRune(reader, pi, '$')

	// todo this gives us line & col in the source file -- do we want line / col in the vibe block ?
	col := pi.col
	line := pi.line

	ident := parseIdentifier(reader, pi)
	if ident == "" {
		return nil
	}
	if ident == "use" {
		consumeSpaces(reader, pi)

		if !tryParseRune(reader, pi, '(') {
			pi.addError(UseMissingImport)
			// todo? do we want a way to ROLL BACK if we don't have a valid "use" defn. ?
			return nil
		}

		consumeSpaces(reader, pi)
		mru := MetaRefUseImport{
			line: line,
			col:  col,
		}
		ch, size, _ := reader.ReadRune()
		switch ch {
		case '@':
			mru.import_type = UseImportSpace
		case '#':
			mru.import_type = UseImportAgent
		default:
			reader.UnreadRune()
			// todo maybe consume until end parenthesis ... ?
			// or roll back... :/
			pi.addError(UseUnsupportedImport)
			return nil
		}
		pi.col += uint64(size)

		mru.imported = parseIdentifier(reader, pi)

		consumeSpaces(reader, pi)

		if !tryParseRune(reader, pi, ')') {
			pi.addError(MismatchedParens)
		}

		if mru.imported == "" {
			pi.addError(UseUnsupportedImport)
			return nil
		}

		return &mru
	} else {
		consumeSpaces(reader, pi)
		mrt := MetaRefTask{
			ident: ident,
			line:  line,
			col:   col,
			args:  parseParams(reader, pi),
		}
		return &mrt
	}
}

func parseMetaRefPath(reader io.RuneScanner, pi *ParserInfo) *MetaRefPath {
	tryParseRune(reader, pi, '=')

	var mrp MetaRefPath
	mrp.line = pi.line
	mrp.col = pi.col

	mrp.ident = parseIdentifier(reader, pi)
	if mrp.ident == "" {
		return nil
	}
	return &mrp
}
