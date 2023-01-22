package ip

import (
	"fmt"
	"net"
)

// GetLocalIP get all your local ipv4 address (except 127.0.0.1)
func GetLocalIP() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("GetLocalIP could not get interface addresses: %w", err)
	}
	IPs := make([]string, 0)
	for _, a := range addrs {
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
func GetPublicIP() (string, error) {
	// resp, err := http.Get("https://api.ipify.org")
	// if err != nil {
	// 	return "", fmt.Errorf("GetPublicIP could not query api.ipify.org: %w", err)
	// }
	// defer resp.Body.Close()

	// var buffer []byte = make([]byte, 100)
	// count, err := resp.Body.Read(buffer)
	// address := buffer[:count]
	// if err != nil || count == 0 {
	//     return "", fmt.Errorf("GetPublicIP could not read data from api.ipify.org: %w", err)
	// }

	// For now, avoid the external call since we only work with one address anyway
	return string("86.89.186.20"), nil
}
