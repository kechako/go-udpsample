package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8080", "UDP server address")
	flag.Parse()

	var wg sync.WaitGroup
	defer wg.Wait()

	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go read(ctx, conn, &wg)

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		log.Print(err)
		return
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := s.Text()
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Print(err)
			return
		}
		if text == "exit" {
			break
		}
	}
	if err := s.Err(); err != nil {
		log.Print(err)
		return
	}
}

func read(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	var buf [1500]byte
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Print(err)
				break loop
			}

			msg := string(buf[:n])

			fmt.Printf("server respond : %s\n", msg)
		}
	}
}
