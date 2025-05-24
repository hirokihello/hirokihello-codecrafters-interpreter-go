package run

type Env struct {
	variables       *map[string]EvaluateNode
	parentVariables *map[string]EvaluateNode
	functions       *map[string]Function
	parentFunctions *map[string]Function
}

type Function struct {
	name       string
	parameters []string
	statements []Statement
	closure    *Env
}

func NewEnv() *Env {
	env := &Env{
		variables:       &map[string]EvaluateNode{},
		parentVariables: &map[string]EvaluateNode{},
		functions:       &map[string]Function{},
	}
	return env
}

var globalEnv *Env

var funcGlobalEnv *Env = &Env{
	variables:       &map[string]EvaluateNode{},
	parentVariables: &map[string]EvaluateNode{},
	functions:       &map[string]Function{},
}

func getGlobalEnv() *Env {
	if globalEnv == nil {
		globalEnv = NewEnv()
	}
	return globalEnv
}

// NewChildEnv creates a new child environment that inherits from the current environment.
func (e *Env) NewChildEnv() *Env {
	newVariables := make(map[string]EvaluateNode)
	newFunctions := make(map[string]Function)

	// Copy the variables and functions from the parent environment
	for k, v := range *e.variables {
		newVariables[k] = v
	}
	for k, v := range *e.functions {
		newFunctions[k] = v
	}

	return &Env{
		variables:       &newVariables,
		functions:       &newFunctions,
		parentVariables: e.variables,
		parentFunctions: e.functions,
	}
}

func (e *Env) newClosureEnv() *Env {
	newVariables := make(map[string]EvaluateNode)
	newFunctions := make(map[string]Function)
	newParentVariables := make(map[string]EvaluateNode)
	newParentFunctions := make(map[string]Function)

	// Copy the variables and functions from the parent environment
	for k, v := range *e.variables {
		newVariables[k] = v
		newParentVariables[k] = v
	}
	for k, v := range *e.functions {
		newFunctions[k] = v
		newParentFunctions[k] = v
	}

	return &Env{
		variables:       &newVariables,
		functions:       &newFunctions,
		parentVariables: &newParentVariables,
		parentFunctions: &newParentFunctions,
	}
}
