package main

import (
	"context"
	"fmt"
	"github.com/sourcegraph/jsonrpc2"
	"net"
	"time"
)

func main() {
	ctx := context.Background()
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8848", 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	client := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), nil)
	res := ""
	if err := client.Call(ctx, "eth_blockNumber", "pong", &res, nil); err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
