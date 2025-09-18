package timescalDb

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// 场景不复杂，先不用接口定义，直接写

type PriceTicks struct {
	Time           time.Time
	Symbol         string
	Price          float64
	Source         string
	Volume         float64
	Quote          float64
	TradeNum       int64
	TakerBuyVolume float64
	TakerBuyQuote  float64
}

func (p *timescaleDb) InsertPriceTickBatch(ctx context.Context, ticks []*PriceTicks, interval string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // 大批量所以 5 秒超时
	defer cancel()
	sql := `INSERT INTO %s (time, symbol, price, source, volume, quote, trade_num, taker_buy_volume, taker_buy_quote) VALUES %s`
	values := ""
	for _, tick := range ticks {
		values = values + fmt.Sprintf("('%s', '%s', %f, '%s', %f, '%f', '%d', '%f', '%f'),",
			tick.Time.Format(time.RFC3339), tick.Symbol, tick.Price, tick.Source, tick.Volume, tick.Quote, tick.TradeNum, tick.TakerBuyVolume, tick.TakerBuyQuote)
	}
	values = strings.TrimRight(values, ",")
	sql = fmt.Sprintf(sql, tablePriceTicksName+"_"+interval, values)
	fmt.Println(sql)
	_, err := p.conn.Exec(ctx, sql)
	return err
}

func (p *timescaleDb) InsertPriceTick(ctx context.Context, tick *PriceTicks) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3) // 暂定 3 秒超时
	defer cancel()
	sql := `INSERT INTO %s (time, symbol, price, source, volume, type) VALUES ($1, $2, $3, $4, $5, $6)`
	sql = fmt.Sprintf(sql, tablePriceTicksName)
	_, err := p.conn.Exec(ctx,
		sql,
		tick.Time, tick.Symbol, tick.Price, tick.Source, tick.Volume)
	return err
}
