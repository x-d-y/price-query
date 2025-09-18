## 价格查询服务


### 支持功能

- [x] 币安秒级市场价格查询 [币安相关接口](https://developers.binance.com/docs/binance-spot-api-docs/rest-api/market-data-endpoints)
- [x] Coinpaprika 秒级市场价格查询 [Coinpaprika 文档](https://api.coinpaprika.com/)
- [ ] CoinMarketCap // todo
- [ ] cryptoToCompare // todo
- [ ] okx // todo

##### 数据库
- 使用 timescaledb [docker启动](https://docs.tigerdata.com/self-hosted/latest/install/installation-docker/)

- 启动
```aiignore

docker network create hushine-tech

docker run -d --name timescaledb --network hushine-tech -p 5432:5432 -v /home/xdy/develop/data/price-query/timescalDb_data:/pgdata -e PGDATA=/pgdata -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg17

CREATE DATABASE cryptocurrency;

CREATE TABLE price_ticks_$(interval)
(
    time   TIMESTAMPTZ     NOT NULL,
    symbol VARCHAR(10)     NOT NULL, -- 币种代码，比如 BTCUSDT, ETHUSDT 等
    price  NUMERIC(12, 4)  NOT NULL,
    source VARCHAR(32)     NOT NULL,
    volume NUMERIC(30, 18) NULL,     -- 交易量
    quote  NUMERIC(30, 18) NULL,     -- 交易额
    taker_buy_volume       NUMERIC(30, 18) NULL,   -- taker buy 交易量
    taker_buy_quote        NUMERIC(30, 18) NULL,   -- taker buy 交易额
    trade_num              INT NULL, -- 交易订单数
    PRIMARY KEY (time, source, symbol)
);


其中 $(interval) 是可以根据间隔更换的， 例如 1s, 1m ， 15m, 1w , 1M 等



```




### 代码编译
make build

### 运行
make run



