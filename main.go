package main

import (
    "fmt"
    "net"
	"strings"
)

func main() {
	myAddress := ""

    addrs, err := net.InterfaceAddrs()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, addr := range addrs {
        // Check if the address is an IP address
        ipNet, ok := addr.(*net.IPNet)
        if !ok {
            continue
        }

        // Skip loopback addresses
        if ipNet.IP.IsLoopback() {
            continue
        }

        // Check if the address is an IPv4 address
        if ipNet.IP.To4() != nil {
            fmt.Println("Local IP address:", ipNet.IP.String())
        }

		if strings.HasPrefix(ipNet.IP.String(), "10.") {
			myAddress = ipNet.IP.String()
		}
    }

	fmt.Println("my IP is", myAddress)
}
