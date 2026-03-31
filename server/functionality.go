package server

import (
	"fmt"
	"net"
	"strings"
	"time"
)

/* CONNECTION MODULE */

func (si *ServerInstance) establishnodetonode(connection *conn) {
	return
}

/*BROADCAST FUNCTIONALITY*/

// Send an alive message
func (si *ServerInstance) sendbroadcast() {
	go func() {
		addr := net.UDPAddr{
			Port: 12345,
			IP:   net.ParseIP("255.255.255.255"),
		}

		con, err := net.DialUDP("udp", nil, &addr)
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer con.Close()

		// Learning: If i'm just looping over one channel I can do this
		for si.broadcasting {
			message := []byte("[QFSERVER]ALIVEPING")
			_, err = con.Write(message)

			if err != nil {
				fmt.Printf("Error: %v", err)
			}

			time.Sleep(time.Second * 2)
			fmt.Println("BROADCAST: Sending a broadcast")
		}
	}()
}

// Listen for alive
func (si *ServerInstance) listenbroadcast() {
	go func() {
		addr := net.UDPAddr{
			Port: 12345,
			IP:   net.ParseIP("0.0.0.0"),
		}

		con, err := net.ListenUDP("udp", &addr)
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer con.Close()

		con.SetDeadline(time.Now().Add(5 * time.Second))

		si.buffer = make([]byte, 1024)

		for si.broadcasting {

			n, addr, err := con.ReadFromUDP(si.buffer)
			if err != nil {
				break
			}

			// Check duplicates
			_, exists := si.pingpool[addr.String()]
			senderhostname, errhostname := net.LookupHost(addr.IP.String())
			if !exists {

				if errhostname != nil {
					// Store [address] = hostname
					si.pingpool[addr.String()] = strings.Join(senderhostname, " ")
				} else {
					fmt.Printf("Could not resolve hostname!\n")
					si.pingpool[addr.String()] = ""
				}
			}

			fmt.Printf("Received response from %s: %s\n", addr.String(), string(si.buffer[:n]))
		}
	}()
}
