package run

// FunctionValue represents a function as a value
type FunctionValue struct {
	function Function
}

func (f *FunctionValue) getValue(env *Env) EvaluateNode {
	// Function自体を返す
	return EvaluateNode{
		value:     "",
		valueType: "function_value",
	}
}

func (f *FunctionValue) getType() string {
	return "function"
}