package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"goq/internal/task"

)

const MAX_RETRIES = 3
var msg_retries int

var handlers = map[string]func(int, task.Task) error{
	"email" : handleEmail,
	"message" : handleMessage,
}

func handleEmail(id int, t task.Task) error{
	log.Printf("[worker %d] - [task %s] - executing mailing task...", id, t.Id[:4])
	time.Sleep(time.Second * 3)
	log.Printf("[worker %d] - [task %s] - task completed !", id, t.Id[:4])
	return nil
}

func handleMessage(id int, t task.Task) error{
	msg_retries++
	log.Printf("[worker %d] - [task %s] - executing messaging task...", id, t.Id[:4])
	time.Sleep(time.Second * 2)
	if msg_retries < 2{
		err_msg := fmt.Sprintf("[worker %d] - [task %s] - falied to send message",id, t.Id[:4])
		return errors.New(err_msg)
	}
	log.Printf("[worker %d] - [task %s] - task completed !", id, t.Id[:4])
	return nil
}

func Worker(id int, cfg *task.Config){
	
	for{
		data, err := cfg.Rdb.BRPop(cfg.Ctx, 0,cfg.MainQueue).Result()
		if err != nil{
			log.Printf("[worker %d] - error during redis brpop: %s\n", id, err.Error())
			continue
		}

		var t task.Task

		if err := json.Unmarshal([]byte(data[1]), &t); err != nil{
			log.Printf("[worker %d] - error trying to decode the request body: %s\n", id, err.Error())
			continue
		}

		handler, ok := handlers[t.Type]
		if !ok {
			log.Printf("[worker %d] - no handler for task type of: %s\n", id, t.Type)
			continue
		}

		if err := handler(id, t); err != nil{
			retry(id, t, cfg, err)
		}

	}
}