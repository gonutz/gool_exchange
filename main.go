package main

import (
    "fmt"
    "net"
    "golang.org/x/sys/windows"
    "syscall"
    "reflect"
)

func main() {
    // Resolve the broadcast address and port
    broadcastAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:12345")
    if err != nil {
        fmt.Println("Error resolving broadcast address:", err)
        return
    }

    // Create a UDP connection
    conn, err := net.ListenUDP("udp4", broadcastAddr)
    if err != nil {
        fmt.Println("Error setting up UDP listener:", err)
        return
    }
    defer conn.Close()

    // Enable broadcasting on the UDP connection
    fd := connFd(conn)
    if fd == syscall.InvalidHandle {
        fmt.Println("Error getting file descriptor")
        return
    }

    if err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_BROADCAST, 1); err != nil {
        fmt.Println("Error setting socket option:", err)
        return
    }

    // Start a goroutine to send broadcasts
    go func() {
        for {
            message := []byte("Hello from Server!")
            _, err = conn.WriteToUDP(message, broadcastAddr)
            if err != nil {
                fmt.Println("Error sending broadcast:", err)
                return
            }
            fmt.Println("Broadcast message sent")
        }
    }()

    buffer := make([]byte, 1024)

    // Listen for incoming broadcast messages
    for {
        n, addr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error receiving data:", err)
            continue
        }

        fmt.Printf("Received message from %s: %s\n", addr, string(buffer[:n]))
    }
}

func connFd(conn *net.UDPConn) syscall.Handle {
    connVal := reflect.ValueOf(conn)
    connFd := connVal.Elem().FieldByName("fd").Elem().FieldByName("pfd").Elem().FieldByName("Sysfd")
    return syscall.Handle(connFd.Uint())
}
