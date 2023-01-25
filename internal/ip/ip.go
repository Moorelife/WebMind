package ip

import (
	"fmt"
	"net"
	"strings"
)

// GetLocalIP get all your local ipv4 address (except 127.0.0.1)
func GetLocalIP() ([]string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("GetLocalIP could not get interface addresses: %w", err)
	}
	IPs := make([]string, 0)
	for _, a := range address {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				IPs = append(IPs, ipNet.IP.To4().String())
			}
		}
	}
	return IPs, nil
}

// GetOutboundIP get the outbound ip, especially useful when you have multi
// local ipv4 ip, and you want figure out which one is being used
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("GetOutboundIP could not reach Google DNS: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GetPublicIP gets your public ip
func GetPublicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}
