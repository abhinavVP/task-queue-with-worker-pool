# Task Queue 

This is a lightweight background job processing system written in Go and uses Redis to manage the queues.
It supports workers, delayed retries, exponential backoff with jitter, and a dead-letter queue — all implemented from scratch for learning and experimentation.

---

## **Features**

- **Background job processing** using Redis lists and blocking pops (`BRPOP`).  
- **Task model** with unique IDs, payloads, attempt counters, and last error details.  
- **Retry engine** with **exponential backoff + jitter**.  
- **Delayed retry queue** implemented with **Redis Sorted Sets (ZSET)** using timestamps as scores.  
- **Scheduler** that promotes due delayed tasks back into the main queue.  
- **Dead Letter Queue (DLQ)** for tasks that exceed the retry limit, with safe fallback JSON.  
- **Modular architecture** (worker, scheduler, retry logic, task definitions, config).  

---

## **How It Works**

1. Tasks are enqueued into the **main Redis list** by sending an /enqueue HTTP request.  
2. Workers consume tasks using `BRPOP` and execute handlers.  
3. Failed tasks are retried with **exponential backoff and jitter**.  
4. Retries are scheduled using a **delayed ZSET(sorted sets)**, keyed by future timestamps.  
5. The scheduler periodically scans the delayed ZSET and pushes due tasks back to the main queue.  
6. Tasks that exceed the maximum retry count are moved into the **DLQ(also sorted sets)** for inspection.

---

## **Requirements**

- Go 1.22+  
- Redis (local or via Docker)

---

## **How to use**
- Start the Redis server (locally or through docker).
- Start the http server by running cmd/api/main.go.
- Start the worker-loop by running cmd/worker/main.go. (this runs the scheduler as well)
- Enqueue tasks by sending an http request to the server (Note: You can only enqueue one task per request at the moment).

---

## **Further Work**
- Implement the functionality for inspecting the DLQ.
- convert the entire thing into a CLI tool for ease of use.
