package main

import (
	"context"
	"encoding/json"
	"fmt"
	"goq/internal/task"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Protocol: 2,
})

var ctx = context.Background()
var q = task.MQ


func enqueueHandler(w http.ResponseWriter, r *http.Request){
	var task task.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	task.Id = uuid.New().String()[:4]
	task.Error = ""
	data, err := json.Marshal(task)
	if err != nil{
		http.Error(w, "error trying to marshal the request", http.StatusInternalServerError)
		return
	}
	
	err = rdb.LPush(ctx, q, data).Err()
	if err != nil{
		http.Error(w, "failed to enqueue task", http.StatusInternalServerError)
		return
	}
	fmt.Printf("task enqueued with id %s", task.Id)
	w.Write([]byte("task enqueued\n"))

}

func main(){
	http.HandleFunc("/enqueue", enqueueHandler)	
	log.Println("server running on 8080")
	http.ListenAndServe(":8080", nil)
}