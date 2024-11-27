package main

import (
    "fmt"
	"time"
    "net"
	"strings"
)

const (
	broadcastIP = "10.0.0.255"
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

	localAddr, err := net.ResolveUDPAddr("udp", myAddress+":0")
    if err != nil {
		return err
    }

    broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:" + pingPort)
    if err != nil {
		return err
    }

    conn, err := net.DialUDP("udp", localAddr, broadcastAddr)
    if err != nil {
		return err
    }
    defer conn.Close()

	// TODO

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





/*package main

import (
    "net"
    "fmt"
    "os"
    "syscall"
)

func main() {
    // Local IP address of the specific network interface
    localIP := "192.168.1.100" // Change this to your specific local IP address

    // Create a local UDP address for the specific interface
    localAddr, err := net.ResolveUDPAddr("udp", localIP+":0")
    if err != nil {
        fmt.Println("Error resolving local address:", err)
        os.Exit(1)
    }

    // Broadcast address and port
    broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:12345")
    if err != nil {
        fmt.Println("Error resolving broadcast address:", err)
        os.Exit(1)
    }

    // Create a UDP connection bound to the specific interface
    conn, err := net.DialUDP("udp", localAddr, broadcastAddr)
    if err != nil {
        fmt.Println("Error creating connection:", err)
        os.Exit(1)
    }
    defer conn.Close()

    // Set the socket option to enable broadcasting
    f, err := conn.File()
    if err != nil {
        fmt.Println("Error getting file descriptor:", err)
        os.Exit(1)
    }
    defer f.Close()

    if err := syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil {
        fmt.Println("Error setting socket option:", err)
        os.Exit(1)
    }

    // Data to broadcast
    message := []byte("Hello, UDP Broadcast!")

    // Send the data
    _, err = conn.Write(message)
    if err != nil {
        fmt.Println("Error sending data:", err)
        os.Exit(1)
    }

    fmt.Println("Broadcast message sent successfully on specific interface!")
}

*/