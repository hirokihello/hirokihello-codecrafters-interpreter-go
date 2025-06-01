package run

type Env struct {
	variables map[string]EvaluateNode
	parentEnv *Env
}

type Function struct {
	name       string
	parameters []string
	statements []Statement
	closure    *Env
}

func NewEnv() *Env {
	return &Env{
		variables: map[string]EvaluateNode{},
		parentEnv: nil,
	}
}


// NewChildEnv creates a new child environment that inherits from the current environment.
func (e *Env) NewChildEnv() *Env {
	return &Env{
		variables: map[string]EvaluateNode{},
		parentEnv: e,
	}
}

// Get looks up a variable in this environment or parent environments
func (e *Env) Get(name string) (EvaluateNode, bool) {
	if val, ok := e.variables[name]; ok {
		return val, true
	}
	if e.parentEnv != nil {
		return e.parentEnv.Get(name)
	}
	return EvaluateNode{}, false
}

// Set sets a variable in the environment where it's defined
func (e *Env) Set(name string, value EvaluateNode) bool {
	if _, ok := e.variables[name]; ok {
		e.variables[name] = value
		return true
	}
	if e.parentEnv != nil {
		return e.parentEnv.Set(name, value)
	}
	return false
}

// Define creates a new variable in the current environment
func (e *Env) Define(name string, value EvaluateNode) {
	e.variables[name] = value
}

