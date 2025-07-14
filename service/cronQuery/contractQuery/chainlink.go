package contractQuery

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang-jwt/jwt/v5"
)

func query() {
	// 读取jwtsecret文件
	secretFile := "/home/xdy/workplace/eth/blockData/geth/jwtsecret"
	secret, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}

	// 生成JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token.Claims = claims

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JWT token:", tokenString)

	// 自定义 HTTP 客户端，在 Header 中加入 Authorization
	header := http.Header{}
	header.Set("Authorization", "Bearer "+tokenString)

	rpcClient, err := rpc.DialHTTPWithClient("http://localhost:8551", &http.Client{
		Transport: &transportWithHeader{header: header},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)
	defer client.Close()

	// 请求 Chainlink 合约示例：查询 latestRoundData
	contractAddress := common.HexToAddress("0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419")
	aggregatorV3ABI := `[{"inputs":[],"name":"latestRoundData","outputs":[{"internalType":"uint80","name":"roundId","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"}]`

	parsedABI, err := abi.JSON(strings.NewReader(aggregatorV3ABI))
	if err != nil {
		log.Fatal(err)
	}

	data, err := parsedABI.Pack("latestRoundData")
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	ctx := context.Background()
	output, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		log.Fatal(err)
	}

	var out struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	}
	err = parsedABI.UnpackIntoInterface(&out, "latestRoundData", output)
	if err != nil {
		log.Fatal(err)
	}

	priceFloat := new(big.Float).SetInt(out.Answer)
	ethPriceUSD, _ := priceFloat.Quo(priceFloat, big.NewFloat(1e8)).Float64()

	fmt.Printf("ETH price in USD: %.8f\n", ethPriceUSD)
}

// transportWithHeader 在每个请求附加自定义 HTTP Header
type transportWithHeader struct {
	header http.Header
}

func (t *transportWithHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.header {
		req.Header[k] = v
	}
	return http.DefaultTransport.RoundTrip(req)
}
