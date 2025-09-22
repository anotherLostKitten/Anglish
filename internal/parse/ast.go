package parse

type Contract struct {
	spaces []SpaceDecl
	// agents []AgentDecl // todo
	// paths []PathDecl // todo
}

type SpaceDecl struct {
	ident string
	space_type SpaceType
	replicable bool
	params []Param
	vibe_desc VibeBlock

	// todo inner decls
	// todo imports

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
	space_source *SpaceDecl
	space_dest *SpaceDecl
	vibe_desc VibeBlock

	line uint64
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

// type DataDecl struct {
// 	ident string
// 	line uint64
// }

type Param struct {
	in_param bool
	data_name string

	line, col uint64
}

type VibeBlock struct {
	vibe_prose string
	// meta_refs []*MetaRef // fixme

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

	args []Param
}

type MetaRefPath struct {
	ident []byte

	line, col uint64
}
