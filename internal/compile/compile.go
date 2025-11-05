package compile

import (
	"fmt"
	"github.com/anotherLostKitten/Anglish/internal/parse"
)

type CompileOutput string

type CompilationUnit interface {
	GetName() parse.Ident
}

type Air struct {
	generated map[parse.Ident]*CompileOutput
}

func Compile(po *parse.ParseOrder) Air {
	air := Air{
		generated: make(map[parse.Ident]*CompileOutput),
	}

	len := po.Length()
	for i := 0; i < len; i++ {
		n := po.GetNode(i)
		deps := po.GetDepIdents(n)

		air.compileUnit(n.Ast_node, deps)
	}

	return air
}

func (air *Air) compileUnit(c CompilationUnit, deps []parse.Ident) {
	name := c.GetName()
	deps_outputs := make([]*CompileOutput, len(deps))
	for i := range deps {
		val, ok := air.generated[deps[i]]
		if !ok {
			panic("Dependency missing")
		}
		deps_outputs[i] = val
	}
	var out CompileOutput
	switch c.(type) {
	case *parse.SpaceDecl:
		s, ok := c.(*parse.SpaceDecl)
		if !ok {panic("Type cast failed")}
		out = compileSpace(s, deps_outputs)
	case *parse.AgentDecl:
		a, ok := c.(*parse.AgentDecl)
		if !ok {panic("Type cast failed")}
		out = compileAgent(a, deps_outputs)
	case *parse.PathDecl:
		p, ok := c.(*parse.PathDecl)
		if !ok {panic("Type cast failed")}
		out = compilePath(p, deps_outputs)
	case *parse.TaskDecl:
		t, ok := c.(*parse.TaskDecl)
		if !ok {panic("Type cast failed")}
		out = compileTask(t, deps_outputs)
	default:
		panic("Unknown compilation unit type")
	}

	air.generated[name] = &out
}

func compileSpace(me *parse.SpaceDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling space\n")
	return "test space"
}

func compileAgent(me *parse.AgentDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling agent\n")
	return "test agent"
}

func compilePath(me *parse.PathDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling path\n")
	return "test path"
}

func compileTask(me *parse.TaskDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling task\n")
	return "test task"
}
