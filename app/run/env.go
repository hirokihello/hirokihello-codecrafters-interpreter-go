package run

type Env struct {
	variables       *map[string]EvaluateNode
	parentVariables *map[string]EvaluateNode
}

func NewEnv() *Env {
	return &Env{
		variables:       &map[string]EvaluateNode{},
		parentVariables: &map[string]EvaluateNode{},
	}
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

	// コピー
	for k, v := range *e.variables {
		(*newEnv.variables)[k] = v
		(*newEnv.parentVariables)[k] = v
	}
	return newEnv
}
