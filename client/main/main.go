package main

import "RPC/client"

func main() {
	client := client.NewClient()
	client.DialAndRequest("0.0.0.0:8080")
}
