package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gilang-anggara/go-routine-pools/pools"
)

const poolSize = 100
const queueSize = 100
const shutdownPeriod = 100 * time.Second
const cooldownPerExecutionPeriod = 10 * time.Millisecond

func main() {
	routinePools := pools.New(poolSize, queueSize, shutdownPeriod, cooldownPerExecutionPeriod)
	routinePools.Start()

	// fire and forget style
	for i := 0; i < 50; i++ {
		i := i
		exec := func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("executing: %d\n", i)
		}
		err := routinePools.Send(pools.Routine{
			ID:          fmt.Sprint(i),
			ExecuteFunc: exec,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// you can also optionally wait for specific goroutines to finish
	finished := make(chan bool)
	err := routinePools.Send(pools.Routine{
		ID: "awaited id",
		ExecuteFunc: func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("executing: awaited goroutine")
		},
		Finished: finished,
	})
	if err != nil {
		log.Fatal(err)
	}
	<-finished

	routinePools.Shutdown()
}
