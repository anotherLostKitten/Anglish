package prompts

type TaskTemplateArgs struct {
	TaskIdentifier string
	TaskVibeBlock  string
	TaskInputs     string
	TaskOutputs    string
	PathEndpoints  string
}

type DataTemplateArgs struct {
	DataIdentifier string
	DataVibeBlock  string
}

const (
	GenericTaskTemplateString = `Function:
Name: {{.TaskIdentifier}}
Type: Function
Description: {{.TaskVibeBlock}}
Inputs: {{.TaskInputs}}
Outputs: {{.TaskOutputs}}
Endpoints: {{.PathEndpoints}}`

	UiSpaceTaskTemplateString = `Component:
Name: {{.TaskIdentifier}}
Type: Component
Description: {{.TaskVibeBlock}}
Inputs: {{.TaskInputs}}
Outputs: {{.TaskOutputs}}
Endpoints: {{.PathEndpoints}}`

	GenericDataTemplateString = `Data:
Name: {{.DataIdentifier}}
Type: Data
Description: {{.DataVibeBlock}}`

	UiSpaceAgentPerSpaceTemplate = `Create an interface according to the following specifications..
Description of interface goal/purpose: {{.SpaceVibeBlock}}
Interface Content:
{{}}{{end}}`
)
