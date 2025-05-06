package client

import (
	"RPC/message"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) DialAndRequest(serverAddress string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	} else {
		c.conn = conn
		fmt.Println("dialing on:", serverAddress)
	}
	defer conn.Close()

	var methodName string
	var params []any

	for {
		fmt.Print("Enter method name (enter exit to exit): ")
		reader := bufio.NewReader(os.Stdin)

		methodNameInput, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read method name:", err)
			return fmt.Errorf("failed to read method name: %w", err)
		}
		methodName = strings.TrimSpace(methodNameInput)
		if methodName == "exit" {
			break
		}

		fmt.Print("Enter parameters (space-separated): ")
		paramsInput, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read parameters:", err)
			return fmt.Errorf("failed to read parameters: %w", err)
		}
		paramsInput = strings.TrimSpace(paramsInput)

		inputParts := strings.Fields(paramsInput)
		params = make([]any, len(inputParts))

		// Parse each element as JSON
		for i, part := range inputParts {
			err = json.Unmarshal([]byte(part), &params[i])
			if err != nil {
				log.Println("failed to parse parameter:", part, "error:", err)
				return fmt.Errorf("failed to parse parameter '%s': %w", part, err)
			}
		}

		fmt.Printf("method: %v\nparams: %v\n", methodName, params)
		res, err := c.call(methodName, params...)
		if err != nil {
			log.Println("call err:", err)
			return fmt.Errorf("failed to call method '%s': %w", methodName, err)
		}
		if len(res) == 0 {
			log.Println("call returned no results")
			return errors.New("no results returned from call")
		}
		fmt.Printf("result: %v\n", res)
	}

	return nil
}

func (c *Client) call(methodName string, params ...any) ([]any, error) {
	requestData := message.RPCMessage{
		MethodName: methodName,
		Parameters: params,
	}

	requestByte, err := json.Marshal(requestData)
	if err != nil {
		log.Println("request encode error:", err)
	}

	_, err = c.conn.Write(requestByte)
	if err != nil {
		log.Println("request write err:", err)
	}

	buf := make([]byte, 4096)
	n, err := c.conn.Read(buf)
	if err != nil {
		log.Println("response read error:", err)
		return nil, err
	}
	responseByte := buf[:n]

	var responseData message.RPCMessage
	err = json.Unmarshal(responseByte, &responseData)
	if err != nil {
		log.Println("response decode error:", err)
		return nil, err
	}

	if responseData.Error == "" {
		return responseData.ReturnValues, nil
	}
	return nil, errors.New(responseData.Error)
}
