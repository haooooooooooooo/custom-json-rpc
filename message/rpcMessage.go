package message

type RPCMessage struct {
	MethodName   string
	Parameters   []any
	ReturnValues []any
	Error        string
}
