/*****************************************************************************
 * http_proxy_DNS.go
 * Names:
 * NetIds:
 *****************************************************************************/

// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	_, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	var host, address string
	// log.Printf("b = %v", string(b[:]))
	reqs := strings.Split(string(b[:]), " ")
	if len(reqs) <= 2 {
		log.Printf("invalid request, b=%v\n", string(b[:]))
		return
	}
	host = reqs[1]

	hostURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostURL.Opaque == "443" {
		address = hostURL.Scheme + ":443"
	} else {
		if strings.Index(hostURL.Host, ":") == -1 {
			address = hostURL.Host + ":80"
		} else {
			address = hostURL.Host
		}
	}

	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}

	req := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", hostURL.Host)
	// log.Printf("req = %s", req)
	_, err = server.Write([]byte(req))
	if err != nil {
		log.Println(err)
		return
	}

	htmlBytes := make([]byte, 0)
	buf := bytes.NewBuffer(htmlBytes)
	_, err = io.Copy(buf, server)
	if err != nil {
		log.Println(err)
		return
	}
	go dnsPrefetch(buf.String())

	_, err = io.Copy(client, buf)
	if err != nil {
		log.Println(err)
		return
	}
}

func dnsPrefetch(htmlResp string) {
	htmlTokens := html.NewTokenizer(strings.NewReader(htmlResp))

loop:
	for {
		tt := htmlTokens.Next()
		switch tt {
		case html.ErrorToken:
			break loop
		case html.StartTagToken:
			t := htmlTokens.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						urlHost, err := url.Parse(attr.Val)
						if err != nil {
							log.Println(err)
							break loop
						}
						_, err = net.LookupHost(urlHost.Host)
						if err != nil {
							log.Println(err)
							break loop
						}
					}
				}
			}
		default:
		}
	}

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 2 {
		log.Fatal("Usage: ./http_proxy [port]")
	}

	port := os.Args[1]
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Panic(err)
	}

	defer l.Close()

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
}
