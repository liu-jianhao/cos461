/*****************************************************************************
 * server-go.go
 * Name:
 * NetId:
 *****************************************************************************/

package main

import (
	"io"
	"log"
	"net"
	"os"
)

const RECV_BUFFER_SIZE = 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	l, err := net.Listen("tcp", "127.0.0.1:"+server_port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// 注意这里不能起goroutine去执行，不然test10和test15过不了
		handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var n int64 = 1
	var err error
	for n > 0 {
		// 用io.Copy最省事，用其他输入输出函数都会有些问题QwQ
		n, err = io.Copy(os.Stdout, conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
