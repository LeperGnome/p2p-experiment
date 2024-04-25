package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"
)

const (
	rxPort = "8888" // puch hole + listen
	txPort = "9999" // send data
)

func main() {
	peerIP := os.Args[1]
	laddr := net.UDPAddrFromAddrPort(
		netip.MustParseAddrPort("0.0.0.0:" + rxPort),
	)
	raddr := net.UDPAddrFromAddrPort(
		netip.MustParseAddrPort(peerIP + ":" + txPort),
	)

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening %s udp...\n", rxPort)
	defer conn.Close()

	// Punching a hole
	go func() {
		_, err := conn.WriteToUDP([]byte("punch"), raddr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("punched a hole from %s to %s\n", laddr, raddr)
	}()

	// Spamming data to connected peer
	go spam(conn, raddr)

	// Listening for incomig data
	for {
		buf := make([]byte, 1024)
		nRead, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("failed reading packet: %s\n", err.Error())
			continue
		}
		fmt.Printf("got data:\n%s\nfrom %v\n", string(buf[:nRead]), addr)
	}
}

func spam(conn *net.UDPConn, raddr *net.UDPAddr) {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			_, err := conn.WriteToUDP([]byte("sample message\n"), raddr)
			if err != nil {
				fmt.Printf(
					"error sending data to remote: %s",
					err.Error(),
				)
			}
		}
	}
}
