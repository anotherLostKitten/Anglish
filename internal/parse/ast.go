package parse

import (
	"strings"
)

type MetaType byte
const (
	SPACE MetaType = iota
	AGENT
	TASK
	PATH
)
type Ident struct {
	t MetaType
	n string
}

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

func (me *SpaceDecl) GetName() Ident {
	return Ident{
		t: SPACE,
		n: me.ident,
	}
}

func (me *SpaceDecl) GetChildren() []ParseUnit {
	a_len := len(me.agents)
	children := make([]ParseUnit, a_len + len(me.tasks))
	for i, a := range me.agents {
		children[i] = a
	}
	for i, t := range me.tasks {
		children[i + a_len] = t
	}
	// DatumDecls?
	return children
}

func (me *SpaceDecl) GetDeps(deps *map[uint64]bool, scope *Scope) bool {
	valid_idents := me.vibe_desc.getDeps(deps, scope)
	if !valid_idents {
		return false
	}

	for _, c := range me.GetChildren() {
		id := c.GetName()
		if !scope.tryAddDep(id, deps) {
			return false
		}
	}
	return true
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


func (me *AgentDecl) GetName() Ident {
	return Ident{
		t: AGENT,
		n: me.ident,
	}
}

func (me *AgentDecl) GetChildren() []ParseUnit {
	return []ParseUnit{}
}

func (me *AgentDecl) GetDeps(deps *map[uint64]bool, scope *Scope) bool {
	return me.vibe_desc.getDeps(deps, scope)
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
	space_source Ident
	space_dest Ident
	vibe_desc VibeBlock

	line_start, line_end uint64
}

func (me *PathDecl) GetName() Ident {
	return Ident{
		t: PATH,
		n: me.ident,
	}
}

func (me *PathDecl) GetChildren() []ParseUnit {
	return []ParseUnit{}
}

func (me *PathDecl) GetDeps(deps *map[uint64]bool, scope *Scope) bool {
	if !scope.tryAddDep(me.space_source, deps) {
		return false
	}
	if !scope.tryAddDep(me.space_dest, deps) {
		return false
	}
	return me.vibe_desc.getDeps(deps, scope)
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


func (me *TaskDecl) GetName() Ident {
	return Ident{
		t: TASK,
		n: me.ident,
	}
}

func (me *TaskDecl) GetChildren() []ParseUnit {
	return []ParseUnit{}
}

func (me *TaskDecl) GetDeps(deps *map[uint64]bool, scope *Scope) bool {
	return me.vibe_desc.getDeps(deps, scope)
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
