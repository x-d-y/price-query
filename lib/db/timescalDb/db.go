package timescalDb

import (
	"context"
	"eth-service-price_query/etc/config"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

var TimescaleDb *timescaleDb

const (
	tablePriceTicksName = "price_ticks"
)

type timescaleDb struct {
	conn *pgx.Conn
	//todo log 等
}

// todo 数据库连接内容放到配置文件, log 要更换
func init() {
	connStr := "postgres://%s:%s@%s:5432/cryptocurrency?sslmode=disable"
	dbConfig := config.GetConfig().DBConfig
	connStr = fmt.Sprintf(connStr, dbConfig.User, dbConfig.Pass, dbConfig.Addr)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if conn == nil {
		panic("Unable to connect to database")
	}
	TimescaleDb = &timescaleDb{
		conn: conn,
	}

}
