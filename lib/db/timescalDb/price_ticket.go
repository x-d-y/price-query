package timescalDb

import (
	"context"
	"fmt"
	"time"
)

// 场景不复杂，先不用接口定义，直接写

type PriceTicks struct {
	Time   time.Time
	Symbol string
	Price  float64
	Source string
	Volume int
}

func (p *timescaleDb) InsertPriceTick(ctx context.Context, tick *PriceTicks) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3) // 暂定 3 秒超时
	defer cancel()
	sql := `INSERT INTO %s (time, symbol, price, source, volume) VALUES ($1, $2, $3, $4, $5)`
	sql = fmt.Sprintf(sql, tablePriceTicksName)
	_, err := p.conn.Exec(ctx,
		sql,
		tick.Time, tick.Symbol, tick.Price, tick.Source, tick.Volume)
	return err
}
