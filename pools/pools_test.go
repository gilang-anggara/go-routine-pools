package pools_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gilang-anggara/go-routine-pools/pools"
)

func Test_RoutinePools_Send_NotStarted_ReturnError(t *testing.T) {
	routinePools := pools.New(1, 0, 0)

	count := 0
	err := routinePools.Send(pools.Routine{
		ID: "1",
		ExecuteFunc: func() {
			time.Sleep(10 * time.Second)
			count += 1
		},
	})

	assert.Equal(t, pools.ErrWorkerNotStarted, err)

	routinePools.Shutdown()
}

func Test_RoutinePools_Send_EmptyChannels_WillExecute(t *testing.T) {
	routinePools := pools.New(1, 100*time.Second, 0)
	routinePools.Start()

	count := 0
	err := routinePools.Send(pools.Routine{
		ID: "2",
		ExecuteFunc: func() {
			time.Sleep(1 * time.Second)
			count += 1
		},
	})

	routinePools.Shutdown()

	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func Test_RoutinePools_Send_FullChannel_ReturnError(t *testing.T) {
	routinePools := pools.New(1, 0, 0)
	routinePools.Start()

	count := 0
	err := routinePools.Send(pools.Routine{
		ID: "3",
		ExecuteFunc: func() {
			time.Sleep(10 * time.Second)
			count += 1
		},
	})

	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	// this will stay in the buffer
	err = routinePools.Send(pools.Routine{
		ID: "4",
		ExecuteFunc: func() {
			count += 1
		},
	})

	assert.Nil(t, err)

	err = routinePools.Send(pools.Routine{
		ID: "5",
		ExecuteFunc: func() {
			count += 1
		},
	})

	assert.Equal(t, pools.ErrPoolsUnavailable, err)

	routinePools.Shutdown()
}

func Test_RoutinePools_Send_WithFinishFlag_WillExecute(t *testing.T) {
	routinePools := pools.New(1, 100*time.Second, 0)
	routinePools.Start()

	count := 0
	finished := make(chan bool)
	err := routinePools.Send(pools.Routine{
		ID: "2",
		ExecuteFunc: func() {
			time.Sleep(1 * time.Second)
			count += 1
		},
		Finished: finished,
	})

	<-finished

	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	routinePools.Shutdown()
}
