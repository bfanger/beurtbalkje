// Improves reload DX when using `dotnet watch`
//
// This server proxies TCP requests to another server.
// If the target server is down, it waits until the target server is up
// If the target server is still down  after the specified timeout, the connection is closed.

package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

var connectionCount int32

func main() {

	pFlag := flag.Int("port", 8888, "proxy server port")
	tFlag := flag.String("target", "", "target server")
	waitFlag := flag.Int("timeout", 30, "timeout in seconds")
	flag.Parse()

	port := *pFlag
	target := *tFlag
	timeout := time.Duration(*waitFlag) * time.Second

	if target == "" {
		args := flag.Args()
		if len(args) != 1 {
			fmt.Println("usage: beurtbalkje [-port=8888] [-timeout=30] target")
			os.Exit(1)
		}
		target = args[0]
	}
	if regexp.MustCompile("^[0-9]+$").MatchString(target) {
		target = "localhost:" + target
	}

	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	fmt.Println(" ______________")
	fmt.Print("/ Beurtbalkje /\\\n\n")
	fmt.Printf("   port : %d\n", port)
	fmt.Printf(" target : %s\n", target)
	fmt.Printf("timeout : %s\n\n", timeout)
	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn, target, timeout)
	}
}

func handleConnection(conn net.Conn, remote string, timeout time.Duration) {
	defer conn.Close()
	atomic.AddInt32(&connectionCount, 1)
	defer func() {
		atomic.AddInt32(&connectionCount, -1)
		printConnections()
	}()

	proxy, err := connectAndRetry(remote, timeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer proxy.Close()
	printConnections()
	read := make(chan error)
	write := make(chan error)

	go copyData(conn, proxy, read)
	go copyData(proxy, conn, write)

	select {
	case <-read:
	case <-write:
	}
	conn.Close()
}

func connectAndRetry(remote string, timeout time.Duration) (net.Conn, error) {
	started := time.Now()
	spinner := []string{"-", "\\", "|", "/"}
	i := 0

	for {
		proxy, err := net.Dial("tcp", remote)
		if err != nil {
			if time.Now().Sub(started) > timeout {
				return nil, fmt.Errorf("timed out after %s", timeout)
			}
			i++
			fmt.Printf("  %s retrying...\r", spinner[i%len(spinner)])
			time.Sleep(50 * time.Millisecond)
		} else {
			if i != 0 {
				fmt.Printf("\rconnected after %s\n", time.Now().Sub(started))
			}
			return proxy, nil
		}
	}
}

func printConnections() {
	fmt.Printf("  %d connections\r", connectionCount)
}

func copyData(dst io.Writer, src io.Reader, result chan error) {
	_, err := io.Copy(dst, src)
	result <- err
}
