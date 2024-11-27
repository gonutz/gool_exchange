package main

import (
    "fmt"
    "net"
	"strings"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	myAddress := ""

    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return err
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

		if strings.HasPrefix(ipNet.IP.String(), "10.") {
			myAddress = ipNet.IP.String()
		}
    }

	fmt.Println("my IP is", myAddress)
	return nil
}
