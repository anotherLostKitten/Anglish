package compile

import (
	"context"
	"fmt"

	"github.com/anotherLostKitten/Anglish/internal/llm"
	"github.com/anotherLostKitten/Anglish/internal/parse"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools"
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
		promptExtension, promptExtensionParseErr := llm.ReadPrompt(llm.ChatSpaceAgentSystem)
		if promptExtensionParseErr != nil {
			panic("Cannot read system prompt extension!")
		}

		systemPrompt += "\n" + promptExtension
	}
	tools := []tools.Tool{llm.EchoTool{}}
	spaceGenAgentExec, newSpaceGenAgentExecErr := llm.NewAgentExecutor(systemPrompt, tools, nil)
	if newSpaceGenAgentExecErr != nil {
		panic("Cannot create space generation agent executor!")
	}

	var generationPrompt CompileOutput
	generationPrompt = "Project Name: " + CompileOutput(me.GetName().N)

	for _, task := range me.Tasks {
		if me.Space_type == parse.CHAT && task.GetName().N == "agentic_main" {
			// TODO: handle CHAT space's agent
			continue
		}
		generationPrompt += "\n" + compileTask(&task, deps)
	}

	// TODO: compile agents

	fmt.Printf("Generation Prompt:\n%s\n", generationPrompt)
	ctx := context.Background()

	response, runErr := chains.Run(ctx, spaceGenAgentExec, string(generationPrompt))

	if runErr != nil {
		panic("Space generation agent execution failed!")
	}

	fmt.Printf("Response:\n%s\n", response)

	return CompileOutput(response)
}

func compileAgent(me *parse.AgentDecl, deps []*CompileOutput) CompileOutput {
	fmt.Printf("compiling agent\n")
	// TODO: implement agent compilation
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
	// TODO: implement path compilation
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

	taskGenPrompt := "Function: " + me.GetName().N + "\nIO:"
	for _, param := range me.Params {
		taskGenPrompt += "\n" + param.ToStr()
	}

	taskGenPrompt += "\nDescription:"
	for _, line := range me.Vibe_desc.Vibe_prose {
		taskGenPrompt += "\n" + line
	}

	return CompileOutput(taskGenPrompt)
}
