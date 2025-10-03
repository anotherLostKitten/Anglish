package parse

import (
	"strings"
)

type Contract struct {
	spaces []SpaceDecl
	agents []AgentDecl
	paths []PathDecl
}

type SpaceDecl struct {
	ident string
	space_type SpaceType
	replicable bool
	params []Param
	vibe_desc VibeBlock

	// inner decls
	agents []AgentDecl
	tasks []TaskDecl
	// data []DatumDecl

	line_start, line_end uint64
}

type SpaceType byte
const (
	UnknownSpace SpaceType = iota
	UI
	IO
	DATA
	FUNC
	AGENTIC
)

type AgentDecl struct {
	ident string
	agent_type AgentType
	params []Param
	vibe_desc VibeBlock

	line_start, line_end uint64
}

type AgentType byte
const (
	UnknownAgent AgentType = iota
	AF
	DF
)

type PathDecl struct {
	ident string
	path_type PathType
	space_source string
	space_dest string
	vibe_desc VibeBlock

	line_start, line_end uint64
}

type PathType byte
const (
	UnknownPath PathType = iota
	INVOKE
	ATTEND
)

type TaskDecl struct {
	ident string
	params []Param
	vibe_desc VibeBlock

	line_start, line_end uint64
}

// type DatumDecl struct {
// 	ident string
// 	line uint64
// }

type Param struct {
	in_param bool
	data_name string

	line, col uint64
}

func (p *Param) ToStr() string {
	if p.in_param {
		return "in=%" + p.data_name
	} else {
		return "out=%" + p.data_name
	}
}

type VibeBlock struct {
	vibe_prose []string
	meta_refs []MetaRef

	line_start, line_end uint64
}

// can do meta_ref.(type) to get type
type MetaRef interface {
	// Line() int
	// Col() int
	ToStr() string
}

type MetaRefData struct {
	ident string

	line, col uint64
}

func (mr *MetaRefData) ToStr() string {
	return "%" + mr.ident
}

type UseImportType byte
const (
	UseImportSpace UseImportType = iota
	UseImportAgent
)

type MetaRefUseImport struct {
	imported string

	import_type UseImportType

	line, col uint64
}

func (mr *MetaRefUseImport) ToStr() string {
	var joiner string
	switch mr.import_type {
	case UseImportSpace: joiner = "@"
	case UseImportAgent: joiner = "#"
	default: panic(-1)
	}
	return "$use(" + joiner + mr.imported + ")"
}

type MetaRefTask struct {
	ident string

	line, col uint64

	args []Param
}

func (mr *MetaRefTask) ToStr() string {
	arg_strs := make([]string, len(mr.args), len(mr.args))
	for i := 0; i < len(mr.args); i++ {
		arg_strs[i] = mr.args[i].ToStr()
	}
	return "$" + mr.ident + "(" + strings.Join(arg_strs, ", ") + ")"
}

type MetaRefPath struct {
	ident string

	line, col uint64
}

func (mr *MetaRefPath) ToStr() string {
	return "=" + mr.ident
}
