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
