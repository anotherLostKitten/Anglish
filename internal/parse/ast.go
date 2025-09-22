package parse

type Contract struct {
	spaces []SpaceDecl
	agents []AgentDecl
	paths []PathDecl
}

type SpaceDecl struct {
	ident string
	space_type SpaceType
	params []Param
	vibe_desc VibeBlock
	parent *SpaceDecl

	// todo inner decls

	line_start, line_end uint64
}

type SpaceType byte
const (
	Unknown SpaceType = iota
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
	Unknown AgentType = iota
	AF
	DF
)

type PathDecl struct {
	ident string
	path_type PathType
	space_source *SpaceDecl
	space_dest *SpaceDecl
	vibe_desc VibeBlock

	line uint64
}

type PathType byte
const (
	Unknown PathType = iota
	READ
	WRITE
)

type TaskDecl struct {
	ident string
	params []Param
	vibe_desc VibeBlock

	line_start, line_end uint64
}

type DataDecl struct {
	ident string

	line uint64
}

type Param struct {
	name, binds []byte
	mem_source *string

	line, col uint64
}

type VibeBlock struct {
	vibe_prose string
	meta_refs []*MetaRef

	line_start, line_end uint64
}

// can do meta_ref.(type) to get type
type MetaRef interface {
	MetaRefData | MetaRefUseImportSpace | MetaRefUseImportAgent | MetaRefTask | MetaRefPath

	// Line() int
	// Col() int
	// ToStr() string
}

type MetaRefData struct {
	ident []byte

	line, col uint64
}

type MetaRefUseImportSpace struct {
	ident []byte

	line, col uint64

	imported []byte
}


type MetaRefUseImportAgent struct {
	ident []byte

	line, col uint64

	imported []byte
}

type MetaRefTask struct {
	ident []byte

	line, col uint64

	args []Params
}

type MetaRefPath struct {
	ident []byte

	line, col uint46
}
