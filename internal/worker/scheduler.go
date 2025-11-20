package worker

import (
	"goq/internal/task"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func Scheduler(cfg *task.Config){
	for{
		tasks, err := cfg.Rdb.ZRangeByScore(cfg.Ctx, cfg.DelayedQueue, &redis.ZRangeBy{
			Min: "-inf",
			Max: strconv.FormatInt(time.Now().Unix(), 10),
		}).Result()
		
		if err != nil{
			log.Printf("[scheduler] - error trying to pull from the delayed queue: %s", err.Error())
		}

		if len(tasks) > 0{
			for _, t := range tasks{
				if err := cfg.Rdb.LPush(cfg.Ctx, cfg.MainQueue, t).Err(); err != nil{
					log.Printf("[scheduler] - error trying to push back to main queue: %s", err.Error())
				}
				if err := cfg.Rdb.ZRem(cfg.Ctx, cfg.DelayedQueue, t).Err(); err != nil{
					log.Printf("[scheduler] - error trying to remove from the delayed queue: %s", err.Error())
				}
			}
		}

		time.Sleep(time.Millisecond * 500)

	}
}