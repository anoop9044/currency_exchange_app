package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	r "github.com/go-redis/redis/v8"

	"currency_exchange_app/internal/model"
	"currency_exchange_app/internal/redis"
)
const (
	historicalKeyPrefix = "historical_rates:"
	redisKeyPrefix = "CEA:"
)

func UpdateExchangeRate(ctx context.Context, username string, currency string, rate float64) {
	exchangeRate := model.ExchangeRate{
		Currency: currency,
		Rate:     rate,
	}

	data, err := json.Marshal(exchangeRate)
	if err != nil {
		log.Printf("Could not marshal exchange rate: %v", err)
		return
	}

	// Store current rate in Redis
	err = redis.Client.Set(ctx, redisKeyPrefix+currency, data, 0).Err()
	if err != nil {
		log.Printf("Could not set exchange rate in Redis: %v", err)
		return
	}

	// Store historical rate with timestamp
	timestamp := time.Now().Unix()
	historicalKey := historicalKeyPrefix + currency
	err = redis.Client.ZAdd(ctx, redisKeyPrefix+historicalKey, &r.Z{
		Score:  float64(timestamp),
		Member: data,
	}).Err()
	if err != nil {
		log.Printf("Could not set historical exchange rate in Redis: %v", err)
		return
	}

	// Fetch all exchange rates and publish
	var rates []model.ExchangeRate
	for _, currency := range []string{"INR", "USD", "EURO"} {
		key := redisKeyPrefix + currency
		val, err := redis.Client.Get(ctx,key ).Result()
		println(key,val,err)
		if err != nil {
			log.Printf("Could not get exchange rate from Redis for curency : %v", currency)
			continue
		}

		var rate model.ExchangeRate
		err = json.Unmarshal([]byte(val), &rate)
		if err != nil {
			log.Printf("Could not unmarshal exchange rate: %v", err)
			continue
		}
		rates = append(rates, rate)
	}

	
	// Publish updated rates to all connected clients
	redis.PublishRateUpdate(ctx, rates)
		
}


func GetHistoricalRates(ctx context.Context, currency string) ([]model.ExchangeRate, error) {
    historicalKey := redisKeyPrefix + historicalKeyPrefix + currency
    zRange := redis.Client.ZRangeWithScores(ctx, historicalKey, 0, -1)
    if zRange.Err() != nil {
        return nil, zRange.Err()
    }

    var rates []model.ExchangeRate

    for _, z := range zRange.Val() {
        var rate model.ExchangeRate
        err := json.Unmarshal([]byte(z.Member.(string)), &rate)
        if err != nil {
            log.Printf("Could not unmarshal historical exchange rate: %v", err)
            continue
        }
        timestamps :=  int64(z.Score)
		rate.UpdateTime = &timestamps
        rates = append(rates, rate)

    }

    return rates, nil
}
