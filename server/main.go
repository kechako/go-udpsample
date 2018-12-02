package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/kechako/go-udpsample/netutil"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "UDP server address")
	flag.Parse()

	listenConfig := &net.ListenConfig{
		Control: netutil.ListenControl,
	}
	conn, err := listenConfig.ListenPacket(context.Background(), "udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var buf [1500]byte
	for {
		n, addr, err := conn.ReadFrom(buf[:])
		if err != nil {
			log.Print(err)
			break
		}

		msg := string(buf[:n])
		if msg != "hello" {
			log.Print("invalid message received")
			continue
		}

		log.Printf("client joined : %v", addr)

		d := net.Dialer{
			LocalAddr: conn.LocalAddr(),
			Control:   netutil.ListenControl,
		}
		clConn, err := d.Dial(addr.Network(), addr.String())
		if err != nil {
			log.Print(err)
			continue
		}

		go clientRead(clConn)

	}
}

func clientRead(conn net.Conn) {
	defer conn.Close()

	var buf [1500]byte
	for {
		// timeout in 30 seconds
		conn.SetDeadline(time.Now().Add(30 * time.Second))

		n, err := conn.Read(buf[:])
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				log.Printf("connection timeout : %v", conn.RemoteAddr())
				break
			}
			log.Print(err)
			break
		}

		msg := string(buf[:n])
		if msg == "exit" {
			log.Printf("client leaved : %v", conn.RemoteAddr())
			break
		}

		log.Printf("received from %v, Message : %s", conn.RemoteAddr(), msg)

		// echo client
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Print(err)
			break
		}
	}
}
