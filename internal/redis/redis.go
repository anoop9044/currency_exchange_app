package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"

	"currency_exchange_app/internal/model"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

const (
	historicalKeyPrefix = "historical_rates:"
	redis_key_prefix = "CEA:"
)
// InitRedis initializes the Redis client
func InitRedis(addr string) error {
	Client = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func PublishRateUpdate(ctx context.Context, rates []model.ExchangeRate) {
	data, err := json.Marshal(rates)
	if err != nil {
		log.Printf("Could not marshal exchange rates: %v", err)
		return
	}
	println("Publishing rate update",data)
	err = Client.Publish(ctx, "rate_updates", data).Err()
	if err != nil {
		log.Printf("Could not publish rate update: %v", err)
	}
}
