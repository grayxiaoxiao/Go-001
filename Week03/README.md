学习笔记

#### 作业
> 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
> [作业提交地址](https://github.com/Go-000/Go-000/issues/69)

#### 解答思路
> 级联注销或关闭功能，本质上就是进程或线程间的通讯问题，在各自的程序体内设置对应信号的监听接受就可以解决级联注销的问题。
> golang中chanel刚好可以用来实现信号的发送与接受功能，而errgroup则是用来实现存储各自程序体可能产生的错误信息，以便针对处理。

```golang
package main

import (
    "fmt"
    "net/http"
    "context"
    "os"
    "golang.org/x/sync/errgroup"
)

func server(stop <-chan struct{}) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request){
        fmt.Fprintln(resp, "Hello, Serve is running...")
    })

    server := http.Server{
        Addr: "0.0.0.0:9090",
        Handler: mux,
    }
    go func() {
        <-stop
        server.Shutdown(context.Background())
    }()
    return server.ListenAndServe()

}

func listener(stop <-chan struct{}) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/listener", func(resp http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(resp, "Hello, I am is listener!!")
        // 通过传递参数模拟出错以达到关闭监听目的
        stop, _ := req.URL.Query()["stop"]
        if len(stop) > 0 {
            os.Exit(0)
        }
    })
    listener := http.Server{
        Addr: "127.0.0.1:9091",
        Handler: mux,
    }
    go func() {
        <-stop
        listener.Shutdown(context.Background())
    }()
    return listener.ListenAndServe()
}

func main() {
    stop := make(chan struct{})
    eg   := new(errgroup.Group)
    eg.Go(func() error {
        return server(stop)
    })
    eg.Go(func() error {
        return listener(stop)
    })

    if err := eg.Wait(); err != nil {
        close(stop)
    }
}
```

#### 程序测试结果
###### 使用命令查看端口监听情况 `lsof -iTCP:9090`
```
# 结果如下
COMMAND   PID     USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
main    94738 renzhong    3u  IPv6 0xe235d34f5ea16ba5      0t0  TCP *:websm (LISTEN)
```
###### 使用命令查看端口监听情况 `lsof -iTCP:9091`
```
#结果如下
COMMAND   PID     USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
main    94738 renzhong    4u  IPv4 0xe235d34f7aaeeaa5      0t0  TCP localhost:xmltec-xmlmail (LISTEN)
```
###### 使用`curl -X GET http://localhost:9090`可以获得输出 `"Hello, Serve is running..."`
###### 使用`curl -X GET http://localhost:9091/listener`可以获得输出 `"Hello, I am is listener!!"`
###### 使用`curl -X GET http://localhost:9091/listener?stop=ss`模拟出错关闭监听服务，得到 `"curl: (52) Empty reply from server"`
再次使用`lsof`命令查看端口监听情况，此时两个对应端口已没有监听存在.
