package main

import (
	"fmt"
	"time"

	"github.com/gilang-anggara/go-routine-pools/pools"
)

const maxPoolSize = 100
const shutdownPeriod = 10 * time.Second
const cooldownPerExecutionPeriod = 10 * time.Millisecond

func main() {
	routinePools := pools.New(maxPoolSize, shutdownPeriod, cooldownPerExecutionPeriod)
	routinePools.Start()

	for i := 0; i < 100; i++ {
		i := i

		exec := func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println(i)
		}

		routinePools.Send(pools.Job{
			ID:          fmt.Sprint(i),
			ExecuteFunc: exec,
		})
	}

	routinePools.Shutdown()
}
