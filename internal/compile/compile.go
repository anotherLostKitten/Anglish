package compile

import (
	"fmt"

	"github.com/anotherLostKitten/Anglish/internal/llm"
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
	for i := range len {
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
	switch me := c.(type) {
	case *parse.SpaceDecl:
		out = compileSpace(me, deps_outputs)
	case *parse.AgentDecl:
		out = compileAgent(me, deps_outputs)
	case *parse.PathDecl:
		out = compilePath(me, deps_outputs)
	case *parse.TaskDecl:
		out = compileTask(me, deps_outputs)
	default:
		panic("Unknown compilation unit type")
	}

	air.generated[name] = &out
}

func compileSpace(me *parse.SpaceDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling space\n")
	var systemPrompt string
	var systemPromptParseErr error
	switch me.Space_type {
	case parse.UI:
		systemPrompt, systemPromptParseErr = llm.ReadPrompt(llm.UISpaceAgentSystem)
	case parse.IO:
		// TODO: define prompt
		return ""
	case parse.DATA:
		// TODO: define prompt
		return ""
	case parse.CALL, parse.CHAT:
		systemPrompt, systemPromptParseErr = llm.ReadPrompt(llm.CallSpaceAgentSystem)
	case parse.UnknownSpace:
		panic("Unknown space type")
	}
	if systemPromptParseErr != nil {
		panic("Cannot read system prompt!")
	}
	if me.Space_type == parse.CHAT {
		promptExtension, err := llm.ReadPrompt(llm.AFAgentAgentSystem)
		if err != nil {
			panic("Cannot read system prompt extension!")
		}

		systemPrompt = systemPrompt + "\n" + promptExtension
	}
	return "test space"
}

func compileAgent(me *parse.AgentDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling agent\n")
	var systemPrompt string
	switch me.Agent_type {
	case parse.AF:
		break
	case parse.DF:
		break
	case parse.UnknownAgent:
		panic("Unknown agent type")
	}
	return "test agent"
}

func compilePath(me *parse.PathDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling path\n")
	var systemPrompt string
	switch me.Path_type {
	case parse.INVOKE:
		break
	case parse.ATTEND:
		break
	case parse.UnknownPath:
		panic("Unknown path type")
	}
	return "test path"
}

func compileTask(me *parse.TaskDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling task\n")
	return "test task"
}
