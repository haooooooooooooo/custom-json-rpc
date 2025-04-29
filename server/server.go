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

// NewServer 服务对象创建
func NewServer(args []string) (*Server, error) {
	var ip string
	var port int
	// 处理运行参数
	if len(args) > 1 && args[1] == "-h" {
		fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
		fmt.Printf("-l\t-- 服务端启动命令\n")
		fmt.Printf("<IP>\t-- 服务端监听的 IP 地址 (可为空) (默认ip: 0.0.0.0)\n")
		fmt.Printf("-p\t-- 端口设置命令\n")
		fmt.Printf("<Port>\t-- 服务端监听的端口号\n")
		return nil, errors.New("服务器初始化错误")
	} else if len(args) == 4 {
		if args[1] == "-l" && args[2] == "-p" {
			ip = "0.0.0.0"
			port, _ = strconv.Atoi(args[3])
		} else {
			fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
			fmt.Printf("default IP: 0.0.0.0 (if <IP> is null)\n")
			return nil, errors.New("服务器初始化错误")
		}
	} else if len(args) != 5 {
		fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
		fmt.Printf("default IP: 0.0.0.0 (if <IP> is null)\n")
		return nil, errors.New("服务器初始化错误")
	} else {
		if args[1] == "-l" && args[3] == "-p" {
			ip = args[2]
			port, _ = strconv.Atoi(args[4])
		} else {
			fmt.Printf("usage: %s -l <IP> -p <Port>\n", args[0])
			fmt.Printf("default IP: 0.0.0.0 (if <IP> is null)\n")
			return nil, errors.New("服务器初始化错误")
		}
	}

	// 检测IP格式
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP 格式错误\n")
		return nil, errors.New("服务器初始化错误")
	} else {
		if strings.Contains(ip, ".") {
			fmt.Printf("ipv4: %s\n", ip)
		} else if strings.Contains(ip, ":") {
			fmt.Printf("ipv6: %s\n", ip)
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

// GetAddress 获取服务地址
func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// Register 服务注册
func (s *Server) Register(functionName string, function interface{}) {
	// 若服务名存在则跳过
	if _, exist := s.Service[functionName]; exist {
		return
	}
	// 反射获取函数原型
	functionValue := reflect.ValueOf(function)
	s.Service[functionName] = functionValue
}

// ListenAndServe 在对应地址监听并启用服务
func (s *Server) ListenAndServe(address string) error {
	//s.register(address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	} else {
		fmt.Println("listening on:", address)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		} else {
			fmt.Println("accept client:", conn.RemoteAddr())
		}
		go s.handle(conn)
	}
}

// register 向服务注册中心注册
func (s *Server) register(address string) {
	// 初始化服务信息
	serverInfo := serviceInfo.NewServiceInfo(1, address, s.IP, s.Port, 0)
	marshal, _ := json.Marshal(serverInfo)
	//fmt.Printf("Service is %s\n", string(marshal))

	// 与注册中心连接
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalln("服务器与注册中心连接时发生错误")
	}
	_, err = conn.Write(marshal)
	if err != nil {
		log.Fatalln("服务器与注册中心连接时发生错误")
	}
}

// handle 处理对应套接字的远程调用
func (s *Server) handle(conn net.Conn) {
	// 关闭套接字
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("conn close error:", err)
		}
	}()

	// 初始化缓冲区并读取
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("request read error:", err)
		return
	}
	requestByte := buf[:n]

	// 使用 unmarshal 进行反序列化
	var requestData message.RPCMessage
	err = json.Unmarshal(requestByte, &requestData)
	if err != nil {
		log.Println("request decode error:", err)
		return
	}
	//fmt.Println("requestData:", requestData)

	// 从服务列表获取函数
	function, exist := s.Service[requestData.MethodName]
	if !exist {
		log.Printf("unexpected call: function %s is not exist\n", requestData.MethodName)
		return
	}

	// 获取参数
	inArgs := make([]reflect.Value, 0, len(requestData.Parameters))
	for _, arg := range requestData.Parameters {
		inArgs = append(inArgs, reflect.ValueOf(arg))
	}

	// 反射调用方法
	returnValues := function.Call(inArgs)
	outArgs := returnValues[0].Interface().([]string)
	err = errors.New(returnValues[1].Interface().(string))

	// 构造返回数据
	replyData := message.RPCMessage{
		MethodName:   requestData.MethodName,
		ServiceName:  requestData.ServiceName,
		Parameters:   requestData.Parameters,
		ReturnValues: outArgs,
		Error:        err.Error(),
	}
	//fmt.Println("replyData =", replyData)

	// 使用 marshal 进行序列化
	replyByte, err := json.Marshal(replyData)
	if err != nil {
		log.Println("reply encode error:", err)
		return
	}

	// 将远程调用结果写入套接字
	_, err = conn.Write(replyByte)
	if err != nil {
		log.Println("reply write error:", err)
	}
}
