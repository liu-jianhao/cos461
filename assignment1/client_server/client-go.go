/*****************************************************************************
 * client-go.go
 * Name:
 * NetId:
 *****************************************************************************/

package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const SEND_BUFFER_SIZE = 2048

/* TODO: client()
 * Open socket and send message from stdin.
 */
func client(server_ip string, server_port string) {
	conn, err := net.Dial("tcp", server_ip+":"+server_port)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	input := bufio.NewScanner(os.Stdin)
	// 默认Scanner是ScanLine，这里需要改成ScanBytes
	input.Split(bufio.ScanBytes)
	for input.Scan() {
		_, err = conn.Write(input.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Main parses command-line arguments and calls client function
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: ./client-go [server IP] [server port] < [message file]")
	}
	server_ip := os.Args[1]
	server_port := os.Args[2]
	client(server_ip, server_port)
}
