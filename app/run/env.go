package run

type Env struct {
	variables map[string]EvaluateNode
}

func NewEnv() *Env {
	return &Env{
		variables: make(map[string]EvaluateNode),
	}
}

var globalEnv *Env

func getGlobalEnv() *Env {
	if globalEnv == nil {
		globalEnv = NewEnv()
	}
	return globalEnv
}
