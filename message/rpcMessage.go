package message

type RPCMessage struct {
	MethodName   string   // 函数名
	ServiceName  string   // 服务名
	Parameters   []string // 参数
	ReturnValues []string // 返回值
	Error        string   // 错误信息
}
