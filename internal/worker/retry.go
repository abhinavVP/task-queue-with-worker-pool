package worker

import (
	"math/rand"
	"encoding/json"
	"goq/internal/task"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func pushToDLQ(id int, t task.Task, cfg *task.Config){
	
	log.Printf("[worker %d] - [task %s] - max retires exceeded; killing task!", id, t.Id[:4])

	data,err := json.Marshal(t);
	if err!=nil{
		fallback := map[string]any{
			"id":       t.Id,
			"type":     t.Type,
			"attempts": t.Attempts,
			"error":    t.Error,
		}
		
		data, _ = json.Marshal(fallback)
	}
	
	rn := time.Now().Unix()
	cfg.Rdb.ZAdd(cfg.Ctx, cfg.DeadLetterQueue, redis.Z{
		Score: float64(rn),
		Member: data,
	})
}

func pushToDelayed(id int, t task.Task, cfg *task.Config){

	data, err := json.Marshal(t)
	if err != nil{
		log.Printf("[worker %d] - [task %s] - error trying to marshal during re-enqueue !!", id, t.Id[:4])
		return
	}

	backoff := 1 << t.Attempts
	jitter := rand.Intn(backoff)
	delay := (backoff + jitter)
	run_at := time.Now().Add(time.Second * time.Duration(delay))

	if zerr := cfg.Rdb.ZAdd(cfg.Ctx, cfg.DelayedQueue, redis.Z{
		Score: float64(run_at.Unix()),
		Member: data,
	}).Err(); zerr!= nil{
		log.Printf("[worker %d] - [task %s] - error trying to add to delayed queue : %s", id, t.Id[:4], zerr.Error())
	}
	
	log.Printf("[worker %d] - [task %s] - re-enqueued task; retries remaining: %d", id, t.Id[:4], MAX_RETRIES - t.Attempts)
}

func retry(id int, t task.Task, cfg *task.Config, e error){
	t.Attempts++
	t.Error = e.Error()

	if t.Attempts >= MAX_RETRIES{
		pushToDLQ(id, t, cfg)
		return
	}
	
	pushToDelayed(id, t, cfg)
		
}