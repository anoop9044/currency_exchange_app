package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"nhooyr.io/websocket"

	"currency_exchange_app/internal/model"
	"currency_exchange_app/internal/redis"
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    c, err := websocket.Accept(w, r, nil)
    if err != nil {
        log.Println("Failed to accept WebSocket connection:", err)
        return
    }
    defer c.Close(websocket.StatusInternalError, "the sky is falling")

    ctx := context.Background()

    // Subscribe to Redis pub/sub channel
    pubsub := redis.Client.Subscribe(ctx, "rate_updates")
    defer pubsub.Close()

    ch := pubsub.Channel()

    // Read messages from Redis pub/sub and broadcast to WebSocket clients
    for msg := range ch {
        var data []model.ExchangeRate
        err := json.Unmarshal([]byte(msg.Payload), &data)
        if err != nil {
            log.Println("Failed to Unmarshal message from Redis:", err)
            continue
        }
        log.Printf("Received message from Redis: %+v", data)
        // Convert data to JSON bytes
        jsonData, err := json.Marshal(data)
        if err != nil {
            log.Println("Failed to marshal message to JSON:", err)
            continue
        }

        // Write message to WebSocket client
        err = c.Write(ctx, websocket.MessageText, jsonData)
        if err != nil {
            log.Println("Failed to write message to WebSocket:", err)
            break
        }
    }

    c.Close(websocket.StatusNormalClosure, "")
}
	