package main

import (
	"RPC/server"
	"RPC/testFunction"
	"log"
	"os"
)

func main() {
	args := os.Args

	// 服务注册
	server, err := server.NewServer(args)
	if err != nil {
		log.Fatalln(err)
	}
	server.Register("add", testFunction.Add)
	server.Register("subtract", testFunction.Subtract)
	server.Register("multiply", testFunction.Multiply)
	server.Register("divide", testFunction.Divide)
	server.Register("pow", testFunction.Pow)
	server.Register("sqrt", testFunction.Sqrt)
	server.Register("random", testFunction.Random)
	server.Register("swap", testFunction.Swap)
	server.Register("sort", testFunction.Sort)
	server.Register("sleep", testFunction.Sleep)

	// 服务，启动！
	err = server.ListenAndServe(server.GetAddress())
	if err != nil {
		log.Fatalln("服务器运行时发生错误")
	}
}
