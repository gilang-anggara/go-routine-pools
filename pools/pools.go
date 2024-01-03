package pools

import (
	"errors"
	"fmt"
	"time"
)

type Pools interface {
	Start()
	Shutdown()
	Send(Routine) error
}

func New(maxPoolSize int, shutdownPeriod time.Duration, cooldownPerExecutionPeriod time.Duration) Pools {
	return &pools{
		maxPoolSize:                maxPoolSize,
		shutdownPeriod:             shutdownPeriod,
		cooldownPerExecutionPeriod: cooldownPerExecutionPeriod,
	}
}

type pools struct {
	maxPoolSize                int
	shutdownPeriod             time.Duration
	cooldownPerExecutionPeriod time.Duration
	routines                   chan Routine
}

type Routine struct {
	ID          string
	ExecuteFunc func()
}

var (
	ErrWorkerNotStarted = errors.New("worker has not started yet")
	ErrPoolsUnavailable = errors.New("routine pools are unavailable")
)

func (r *pools) Start() {
	fmt.Println("Starting goroutine pools")

	r.routines = make(chan Routine, r.maxPoolSize)
	go r.worker()

	fmt.Println("Goroutine pools started")
}

func (r *pools) worker() {
	for routine := range r.routines {
		fmt.Println("execution started: ", routine.ID)

		routine.ExecuteFunc()

		fmt.Println("execution finished: ", routine.ID)
	}
}

func (r *pools) Shutdown() {
	r.close()

	fmt.Println("Gracefully shutting down goroutine pools")

	r.await()

	fmt.Println("Goroutine pools shut down")
}

func (r *pools) await() {
	done := make(chan bool)

	go func() {
		start := time.Now()

		isTimedOut := func() bool {
			return time.Since(start) > r.shutdownPeriod
		}

		for len(r.routines) > 0 && !isTimedOut() {
			fmt.Printf("%d goroutine remaining\n", len(r.routines))
			time.Sleep(1 * time.Second)
		}

		done <- true
	}()

	<-done
}

func (r *pools) close() {
	if r.routines == nil {
		return
	}

	close(r.routines)
	fmt.Printf("closing goroutine pools with %d remaining\n", len(r.routines))
}

func (r *pools) Send(routine Routine) error {
	if r.routines == nil {
		return ErrWorkerNotStarted
	}

	select {
	case r.routines <- routine:
		fmt.Printf("routine queued: %s\n", routine.ID)
	default:
		return ErrPoolsUnavailable
	}

	return nil
}
