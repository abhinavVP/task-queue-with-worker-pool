package main

import (
	"context"
	"goq/internal/task"
	"goq/internal/worker"
	"log"

	
)


func main(){

	cfg, err := task.Setup(context.Background())
	if err != nil{
		log.Printf("error during config setup: %s", err.Error())
	}

	workerCount := 2

	log.Printf("starting %d workers\n", workerCount)

	for i:=range workerCount{
		go worker.Worker(i+1, cfg)
	}
	go worker.Scheduler(cfg)

	select{}
}