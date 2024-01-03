package main

import (
	"fmt"
	"time"

	"github.com/gilang-anggara/go-routine-pools/pools"
)

const maxPoolSize = 100
const shutdownPeriod = 100 * time.Second
const cooldownPerExecutionPeriod = 10 * time.Millisecond

func main() {
	routinePools := pools.New(maxPoolSize, shutdownPeriod, cooldownPerExecutionPeriod)
	routinePools.Start()

	// fire and forget style
	for i := 0; i < 10; i++ {
		i := i
		exec := func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("executing: %d\n", i)
		}
		routinePools.Send(pools.Routine{
			ID:          fmt.Sprint(i),
			ExecuteFunc: exec,
		})
	}

	// you can also optionally wait for specific goroutines to finish
	finished := make(chan bool)
	routinePools.Send(pools.Routine{
		ID: "awaited id",
		ExecuteFunc: func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("executing: awaited goroutine")
		},
		Finished: finished,
	})
	<-finished

	routinePools.Shutdown()
}
