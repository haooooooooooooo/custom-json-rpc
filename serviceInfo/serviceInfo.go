package serviceInfo

type ServiceInfo struct {
	Type        int
	ServiceName string
	IP          string
	Port        int
	Load        int
}

func NewServiceInfo(Type int, serviceName, ip string, port, load int) *ServiceInfo {
	return &ServiceInfo{
		Type:        Type,
		ServiceName: serviceName,
		IP:          ip,
		Port:        port,
		Load:        load,
	}
}
