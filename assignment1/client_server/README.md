## assignment1
[assignment1](https://github.com/PrincetonUniversity/COS461-Public/tree/master/assignments/assignment1)

这次作业的任务大概分为两大部分：搭建环境和socket编程。

### 搭建环境
按照[教程](https://github.com/PrincetonUniversity/COS461-Public/tree/master/assignments/assignment1)一步一步搭建即可，期间可能会有各种问题，这时候Google或者百度解决即可。

### socket编程
编程部分分为两部分，c语言版本和python或者go版本的server和client，我这里使用go语言。

建议先仔细看编程要求：
1. server端需要**死循环**等待client端连接
2. server端需要把从客户端传来的数据**输出到标准输出**，之后不需要做任何其他事情
3. server端注意**不要使用fork！！！**
4. client端需要读取标准输入，直到读到EOF
5. client端只需要发送一次消息即可，**发送完即可退出**，不需要死循环！！！

#### C语言部分
- 课程里给了一个教程地址：https://beej.us/guide/bgnet/html/

- 可以看这个教程快速熟悉Linux系统编程和socket编程的知识，如果想系统学习，建议去看《Unix网络编程》，入门必备。

- 其实大概框架可以参考上面的教程的例子写出来，比较麻烦的部分是输入输出部分

贴一下核心代码：
- `server-c.c`:
```c
  while (1)
  { // main accept() loop
    sin_size = sizeof their_addr;
    new_fd = accept(sockfd, (struct sockaddr *)&their_addr, &sin_size);
    if (new_fd == -1)
    {
      perror("accept");
      continue;
    }

    char buf[RECV_BUFFER_SIZE];
    n = recv(new_fd, buf, RECV_BUFFER_SIZE, 0);
    while (n > 0)
    {
      fwrite(buf, n, 1, stdout);
      fflush(stdout);
      n = recv(new_fd, buf, RECV_BUFFER_SIZE, 0);
    }

    close(new_fd);
  }
```
- `client-c.c`:
```c
  numbytes = fread(buf, 1, SEND_BUFFER_SIZE, stdin);
  while (numbytes > 0)
  {
    if (send(sockfd, buf, numbytes, 0) == -1)
    {
      fprintf(stderr, "client: failed to send\n");
      close(sockfd);
      return 1;
    }
    numbytes = fread(buf, 1, SEND_BUFFER_SIZE, stdin);
  }
```

- 其余部分详见代码
- 写完这部分代码，就可以跑测试脚本了，如果顺利的话能通过前5个测试用例。


#### Go语言部分
同理，麻烦的是处理输入输出。

贴一下核心代码：
- `server-go.go`:
```go
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
```

- `client-go.go`:
```go
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
```

- 最后跑测试用例：
```shell
$ ./test_client_server.sh go 11111
================================================================
Testing C client against C server (1/4)
================================================================

1. TEST SHORT MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

2. TEST RANDOM ALPHANUMERIC MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

3. TEST RANDOM BINARY MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

4. TEST SERVER INFINITE LOOP (multiple sequential clients to same server)

SUCCESS: Message received matches message sent!
________________________________________

5. TEST SERVER QUEUE (overlapping clients to same server)

SUCCESS: Message received matches message sent!

================================================================
Testing Go client against Go server (2/4)
================================================================

6. TEST SHORT MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

7. TEST RANDOM ALPHANUMERIC MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

8. TEST RANDOM BINARY MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

9. TEST SERVER INFINITE LOOP (multiple sequential clients to same server)

SUCCESS: Message received matches message sent!
________________________________________

10. TEST SERVER QUEUE (overlapping clients to same server)

SUCCESS: Message received matches message sent!

================================================================
Testing C client against Go server (3/4)
================================================================

11. TEST SHORT MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

12. TEST RANDOM ALPHANUMERIC MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

13. TEST RANDOM BINARY MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

14. TEST SERVER INFINITE LOOP (multiple sequential clients to same server)

SUCCESS: Message received matches message sent!
________________________________________

15. TEST SERVER QUEUE (overlapping clients to same server)

SUCCESS: Message received matches message sent!

================================================================
Testing Go client against C server (4/4)
================================================================

16. TEST SHORT MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

17. TEST RANDOM ALPHANUMERIC MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

18. TEST RANDOM BINARY MESSAGE

SUCCESS: Message received matches message sent!
________________________________________

19. TEST SERVER INFINITE LOOP (multiple sequential clients to same server)

SUCCESS: Message received matches message sent!
________________________________________

20. TEST SERVER QUEUE (overlapping clients to same server)

SUCCESS: Message received matches message sent!

================================================================

TESTS PASSED: 20/20
```

大功告成！