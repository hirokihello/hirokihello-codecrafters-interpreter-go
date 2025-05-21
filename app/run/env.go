package run

type Env struct {
	variables       *map[string]EvaluateNode
	parentVariables *map[string]EvaluateNode
	functions       *map[string]Function
	parentFunctions *map[string]Function
}

type Function struct {
	name        string
	parameters  []string
	statements  []Statement
	environment *Env
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

	return &Env{
		variables:       &newVariables,
		functions:       &newFunctions,
		parentVariables: e.variables,
		parentFunctions: e.functions,
	}
}
