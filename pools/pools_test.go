package pools_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gilang-anggara/go-routine-pools/pools"
)

func Test_RoutinePools_Send_NotStarted_ReturnError(t *testing.T) {
	routinePools := pools.New(1, 1, 0, 0)

	count := 0
	err := routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count += 1
		},
	})

	routinePools.Shutdown()

	assert.Equal(t, pools.ErrWorkerNotStarted, err)
	assert.Equal(t, 0, count)
}

func Test_RoutinePools_Send_EmptyChannels_WillExecute(t *testing.T) {
	routinePools := pools.New(1, 1, 100*time.Second, 0)
	routinePools.Start()

	count := 0
	err := routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count += 1
		},
	})

	routinePools.Shutdown()

	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func Test_RoutinePools_Send_FullChannel_ReturnError(t *testing.T) {
	routinePools := pools.New(1, 1, 0, 0)
	routinePools.Start()

	var count atomic.Int32
	err := routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			time.Sleep(10 * time.Second)
			count.Add(1)
		},
	})

	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	// this will stay in the buffer
	err = routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count.Add(1)
		},
	})

	assert.Nil(t, err)

	err = routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count.Add(1)
		},
	})

	assert.Equal(t, pools.ErrFullQueue, err)

	routinePools.Shutdown()

	assert.Equal(t, int32(0), count.Load())
}

func Test_RoutinePools_Send_WithFinishFlag_WillExecute(t *testing.T) {
	routinePools := pools.New(1, 1, 100*time.Second, 0)
	routinePools.Start()

	count := 0
	finished := make(chan bool)
	err := routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count += 1
		},
		Finished: finished,
	})

	<-finished

	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	routinePools.Shutdown()
}

func Test_RoutinePools_Send_WithMultiplePools_WillExecute(t *testing.T) {
	routinePools := pools.New(10, 100, 100*time.Second, 0)
	routinePools.Start()

	var count atomic.Int32
	for i := 0; i < 100; i++ {
		err := routinePools.Send(pools.Routine{
			ExecuteFunc: func() {
				count.Add(1)
			},
		})

		assert.Nil(t, err)
	}

	routinePools.Shutdown()

	assert.Equal(t, int32(100), count.Load())
}

func Test_RoutinePools_Send_WithDelayPerExecution_WillExecute(t *testing.T) {
	routinePools := pools.New(1, 3, 100*time.Second, 200*time.Millisecond)
	routinePools.Start()

	start := time.Now()
	var count atomic.Int32
	routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count.Add(1)
		},
	})
	routinePools.Send(pools.Routine{
		ExecuteFunc: func() {
			count.Add(1)
		},
	})

	routinePools.Shutdown()

	duration := time.Since(start)

	assert.Equal(t, int32(2), count.Load())
	assert.True(t, duration > 400*time.Millisecond)
	assert.True(t, duration < 500*time.Millisecond)
}
