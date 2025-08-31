package timescalDb

import (
	"context"
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
	connStr := "postgres://postgres:password@timescaledb:5432/cryptocurrency?sslmode=disable"
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
