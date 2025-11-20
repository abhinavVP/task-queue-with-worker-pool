package task

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const MQ = "tasks"
const DQ = "tasks:delayed"
const DLQ = "tasks:dead"
const MAX_RETRIES = 3

type Task struct {
	Id		string                 `json:"id"`
    Type    string                 `json:"type"`
    Payload map[string]any		   `json:"payload"`
	Attempts int				   `json:"attempts"`
	Error string                   `json:"error"`
}

type Config struct{
	Rdb *redis.Client
	Ctx context.Context
	MainQueue string
	DelayedQueue string
	DeadLetterQueue string
}

func Setup(ctx context.Context) (*Config, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
		Protocol: 2,
    })

    return &Config{
        Rdb: rdb,
        Ctx: ctx,
		MainQueue: MQ,
		DelayedQueue: DQ,
		DeadLetterQueue: DLQ,
    }, nil
}