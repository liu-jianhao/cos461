# cos461
普利斯顿大学的计算机网络课程

## assignment1
[assignment1](https://github.com/PrincetonUniversity/COS461-Public/tree/master/assignments/assignment1)

这次作业的任务大概分为两大部分：搭建环境和socket编程。

[assignment1实现](https://github.com/liu-jianhao/cos461/blob/main/assignment1/client_server/README.md)


## assignment5
由于虚拟机里的环境实在搭建不好qwq，所以跳过assignment2-4

这次作业要实现一个HTTP代理服务和DNS Prefetching。

### HTTP代理服务器
简单来说就是客户端给HTTP代理服务器的请求：
```
GET http://www.princeton.edu/ HTTP/1.1
```
然后代理服务器将这个请求转发给远程服务器：
```
GET / HTTP/1.1
Host: www.princeton.edu
Connection: close
(Additional client specified headers, if any...)
```
最后把远程服务器的返回再透传给客户端。

知道这些后就可以开始写代码了。

- 首先，熟悉的TCP服务器代码，就不多说了。
```go
func main() {
	log.SetFlags(log.LstdFlags|log.Lshortfile)

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
```
- 然后handleClientRequest可以分为几部分：
1. 读取请求
```go
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
```
2. 解析请求中的URL
```go
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
```
3. 请求远程服务器
```go
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
```
4. 返回再透传给客户端
```go
	_, err = io.Copy(client, server)
	if err != nil {
		log.Println(err)
		return
	}
```

这样就完成了一个简单的HTTP代理服务器。

最后跑一下测试脚本：
```
$ python test_scripts/test_proxy.py http_proxy
$ python test_scripts/test_proxy_conc.py http_proxy
```

### DNS Prefetching
需求就是解析远程服务器返回的HTML，然后对里面的链接做DMS Prefetching，这种优化就是客户在页面上点其他链接的时候会比较快。

要使用到html库，google用法就行

```go
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
```