package main

import (
    "fmt"
	"time"
    "net"
	"strings"
)

const (
	pingPort = "25468"
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

	if myAddress == "" {
		return fmt.Errorf("own IP address not found")
	}

	fmt.Println("my IP is", myAddress)

    broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:" + pingPort)
    if err != nil {
		return err
    }

    conn, err := net.DialUDP("udp", nil, broadcastAddr)
    if err != nil {
		return err
    }
    defer conn.Close()

	msg := "Hello Network!"
    _, err = conn.Write([]byte(msg))
    if err != nil {
		return err
    }

    listenAddr, err := net.ResolveUDPAddr("udp", ":" + pingPort)
    if err != nil {
		return err
    }

    listenConn, err := net.ListenUDP("udp", listenAddr)
    if err != nil {
		return err
    }
    defer listenConn.Close()

	go func() {
		buffer := make([]byte, 1024)

		for {
	        n, addr, err := listenConn.ReadFromUDP(buffer)
	        if err != nil {
	            fmt.Println("Error receiving data:", err)
	            return
	        }

    	    fmt.Printf("Received message from %s: %s\n", addr, string(buffer[:n]))
    	}
	}()

	time.Sleep(5 * time.Second)

	return nil
}
