package main

import (
	"fmt"
	"time"
)

const (
	maxJobs    = 2000
	maxWorkers = 3000
)

var chanJob chan string
var chanWorkerPool chan chan string
var chanQuit chan struct{}

func init() {
	chanJob = make(chan string, maxJobs)
	chanWorkerPool = make(chan chan string, maxWorkers)
	chanQuit = make(chan struct{})
}

func main() {
	initWorkers()

	for i := 0; i < 100000; i++ {
		chanJob <- fmt.Sprintf("job_%d", i)
		chanWorkerPool <- chanJob
	}

	time.Sleep(2 * time.Second)
	chanQuit <- struct{}{}
}

func work(job string) {
	fmt.Println("handling job:", job)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("handled job:", job)
}

func initWorkers() {
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for {
				select {
				case chanJob := <-chanWorkerPool:
					work(<-chanJob)
				case <-chanQuit:
					return
				}
			}
		}()
	}
}
