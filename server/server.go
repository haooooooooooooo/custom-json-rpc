package server

import (
	"RPC/message"
	"RPC/serviceInfo"
	"RPC/testFunction"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type Server struct {
	IP      string
	Port    int
	Service map[string]reflect.Value
}

func NewServer(args []string) (*Server, error) {
	var ip string
	var port int

	if len(args) > 1 && args[1] == "-h" {
		fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
		fmt.Printf("-l\t-- Command to start the server\n")
		fmt.Printf("<IP>\t-- IP address the server listens on (optional) (default IP: 0.0.0.0)\n")
		fmt.Printf("-p\t-- Command to set the port\n")
		fmt.Printf("<Port>\t-- Port number the server listens on\n")
		return nil, errors.New("server initialization error")
	} else if len(args) == 4 {
		if args[1] == "-l" && args[2] == "-p" {
			ip = "0.0.0.0"
			port, _ = strconv.Atoi(args[3])
		} else {
			fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
			fmt.Printf("default IP: 0.0.0.0 (if <IP> is null)\n")
			return nil, errors.New("server initialization error")
		}
	} else if len(args) != 5 {
		fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
		fmt.Printf("default IP: 0.0.0.0 (if <IP> is null)\n")
		return nil, errors.New("server initialization error")
	} else {
		if args[1] == "-l" && args[3] == "-p" {
			ip = args[2]
			port, _ = strconv.Atoi(args[4])
		} else {
			fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
			fmt.Printf("default IP: 0.0.0.0 (if <IP> is null)\n")
			return nil, errors.New("server initialization error")
		}
	}

	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP format\n")
		return nil, errors.New("server initialization error")
	} else {
		if strings.Contains(ip, ".") {
			fmt.Printf("IPv4: %s\n", ip)
		} else if strings.Contains(ip, ":") {
			fmt.Printf("IPv6: %s\n", ip)
		}
	}

	server := &Server{
		IP:      ip,
		Port:    port,
		Service: map[string]reflect.Value{},
	}
	server.Register("heartCheck", testFunction.HeartCheck)

	return server, nil
}

func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

func (s *Server) Register(functionName string, function any) {
	if _, exist := s.Service[functionName]; exist {
		return
	}

	functionValue := reflect.ValueOf(function)
	s.Service[functionName] = functionValue
}

func (s *Server) ListenAndServe(address string) error {
	//s.register(address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	} else {
		fmt.Println("listening on:", address)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		} else {
			fmt.Println("accept client:", conn.RemoteAddr())
		}
		go s.handle(conn)
	}
}

func (s *Server) register(address string) {
	serverInfo := serviceInfo.NewServiceInfo(1, address, s.IP, s.Port, 0)
	marshal, _ := json.Marshal(serverInfo)
	//fmt.Printf("Service is %s\n", string(marshal))

	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Println("Failed to connect to the registration center:", err)
		return
	}
	_, err = conn.Write(marshal)
	if err != nil {
		log.Println("Failed to connect to the registration center:", err)
		return
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("conn close error:", err)
		}
	}()

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("client disconnected:", conn.RemoteAddr())
			} else {
				log.Println("request read error:", err)
			}
			return
		}
		requestByte := buf[:n]

		var requestData message.RPCMessage
		err = json.Unmarshal(requestByte, &requestData)
		if err != nil {
			log.Println("request decode error:", err)
			continue
		}
		fmt.Println("requestData:", requestData)

		function, exist := s.Service[requestData.MethodName]
		if !exist {
			log.Printf("unexpected call: function %s is not exist\n", requestData.MethodName)
			continue
		}

		inArgs := make([]reflect.Value, len(requestData.Parameters))
		for i, inArg := range requestData.Parameters {
			inArgs[i] = reflect.ValueOf(inArg)
		}

		outArgs := function.Call(inArgs)
		returnValues := outArgs[0].Interface().([]any)
		errValue := ""
		if !outArgs[1].IsNil() {
			errValue = outArgs[1].Interface().(error).Error()
		}

		responseData := message.RPCMessage{
			MethodName:   requestData.MethodName,
			Parameters:   requestData.Parameters,
			ReturnValues: returnValues,
			Error:        errValue,
		}
		fmt.Println("responseData =", responseData)

		responseByte, err := json.Marshal(responseData)
		if err != nil {
			log.Println("response encode error:", err)
			continue
		}

		_, err = conn.Write(responseByte)
		if err != nil {
			log.Println("response write error:", err)
			return
		}
	}
}
