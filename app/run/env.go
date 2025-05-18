package run

type Env struct {
	variables       *map[string]EvaluateNode
	parentVariables *map[string]EvaluateNode
	functions       *map[string]Function
}

type Function struct {
	name       string
	parameters []string
	statements []Statement
}

func NewEnv() *Env {
	env := &Env{
		variables:       &map[string]EvaluateNode{},
		parentVariables: &map[string]EvaluateNode{},
		functions:       &map[string]Function{},
	}
	(*env.variables)["clock"] = EvaluateNode{
		value: "clock",
		valueType: "string",
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
	newEnv := NewEnv()

	// 環境変数の copy
	for k, v := range *e.variables {
		(*newEnv.variables)[k] = v
		(*newEnv.parentVariables)[k] = v
	}

	// functions の copy
	for k, v := range *e.functions {
		(*newEnv.functions)[k] = v
	}

	return newEnv
}
