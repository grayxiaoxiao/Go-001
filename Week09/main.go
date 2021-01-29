package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:9898")
	if err != nil {
		log.Fatalf("Listen Error: %v\n", err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Accept Error: %v\n", err)
			continue
		}
		// reader := bufio.NewReader(conn)
		// request, _ := reader.ReadString('\n')
		// log.Printf(request)
		// go handleConn(conn)
		msg := make(chan string)
		go connReader(context.TODO(), conn, msg)
		go connWriter(context.TODO(), conn, msg)
	}
}

func connReader(ctx context.Context, conn net.Conn, message chan string) {
	reader := bufio.NewReader(conn)
	for {
		request, err := reader.ReadString('\n')
		switch err {
		case nil:
			request = strings.TrimSpace(request)
			if request == ":QUIT" {
				log.Printf("Client request server to close the connection")
				close(message)
			} else {
				message <- request
				log.Println(request)
			}
		case io.EOF:
			log.Printf("Client closed the connection by teminated the process")
			close(message)
		default:
			log.Printf("Read Error: %v\n", err)
			close(message)
		}
	}
}

func connWriter(ctx context.Context, conn net.Conn, msg chan string) {
	for {
		writer := bufio.NewWriter(conn)
		writer.WriteString("Response : ")
		writer.WriteString(<-msg)
		writer.Flush()
		go func() {
			<-msg
			log.Printf("Current msg value is %s", <-msg)
			conn.Close()
		}()
	}
}

// func handleConn(conn net.Conn) {
// 	defer conn.Close()
// 	reader := bufio.NewReader(conn)
// 	writer := bufio.NewWriter(conn)
// 	for {
// 		line, err := reader.ReadString('\n')
// 		switch err {
// 		case nil:
// 			if line == "$QUIT" {
// 				log.Printf("Client request server to close the connection")
// 				return
// 			}
// 		case io.EOF:
// 			log.Printf("Client closed the connection by teminated the process")
// 			return
// 		default:
// 			log.Printf("Read Error: %v\n", err)
// 			return
// 		}
// 		writer.WriteString("ResponseWriter: ")
// 		writer.WriteString(line)
// 		writer.Flush()
// 	}
// }
