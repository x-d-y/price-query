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

docker run -d --name timescaledb -p 5432:5432  -v $(path_to_your_project):/pgdata -e PGDATA=/pgdata -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg17


CREATE DATABASE cryptocurrency;

CREATE TABLE price_ticks (
    time        TIMESTAMPTZ       NOT NULL,
    symbol      VARCHAR(10)       NOT NULL,  -- 币种代码，比如 BTC, ETH, LTC
    price       NUMERIC(12, 4)    NOT NULL,
    source      VARCHAR(32)       NOT NULL,
    volume      INT               NULL,
    PRIMARY KEY (time, source, symbol)
);



```




### 代码编译
make build

### 运行
make run



