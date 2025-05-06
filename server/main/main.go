package main

import (
	"RPC/server"
	"RPC/testFunction"
	"log"
	"os"
)

func main() {
	args := os.Args

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

	err = server.ListenAndServe(server.GetAddress())
	if err != nil {
		log.Fatalln("Server encountered an error while running:", err)
	}
}
