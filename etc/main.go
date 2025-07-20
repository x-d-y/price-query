package main

import (
	"context"
	"eth-service-price_query/lib/safeGo"
	"eth-service-price_query/service/cronQuery/apiQuery"
	"fmt"
	"github.com/sourcegraph/jsonrpc2"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// EthService 结构体
type EthService struct{}

// 定义处理方法（以太坊常见方法名）
func (h *EthService) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "eth_blockNumber":
		// 这里返回一个假的区块高度，实际代码中你可以连接真实以太坊节点获得数据
		result := "0x5BAD55" // 示例区块高度 HEX 编码，返回字符串类型
		conn.Reply(ctx, req.ID, result)
	default:
		fmt.Println("Unknown method:", req.Method)
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: "方法未实现或未找到",
		})
	}
}

func listenOnPort(port int) {
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Println("TCP Listen err:", err)
			return
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("接受TCP连接失败:", err)
				continue
			}
			fn := func(ctx context.Context) {
				jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), &EthService{})
			}

			safeGo.SafeGo(context.Background(), fn)
		}
	}()
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	listenOnPort(8848)
	apiQuery.Init()
	//contractQuery.Init()
	<-sigs
	fmt.Println("shut down ... ...")
	// todo log flush to disk
}
