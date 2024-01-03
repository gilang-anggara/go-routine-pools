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

func New(poolSize int, queueSize int, shutdownPeriod time.Duration, cooldownPerExecutionPeriod time.Duration) Pools {
	return &pools{
		poolSize:                   poolSize,
		queueSize:                  queueSize,
		shutdownPeriod:             shutdownPeriod,
		cooldownPerExecutionPeriod: cooldownPerExecutionPeriod,
	}
}

type pools struct {
	poolSize                    int
	queueSize                   int
	shutdownPeriod              time.Duration
	cooldownPerExecutionPeriod  time.Duration
	routines                    chan Routine
	workerCompletedNotification chan bool
	allWokerCompleted           chan bool
}

type Routine struct {
	ID          string
	ExecuteFunc func()
	Finished    chan bool
}

var (
	ErrWorkerNotStarted = errors.New("worker has not started yet")
	ErrFullQueue        = errors.New("routine queue is full")
)

func (r *pools) Start() {
	fmt.Println("Starting goroutine pools")

	r.routines = make(chan Routine, r.queueSize)
	r.workerCompletedNotification = make(chan bool, r.poolSize)
	r.allWokerCompleted = make(chan bool)

	go r.checkCompletedWorker()

	for i := 0; i < r.poolSize; i++ {
		go r.worker()
	}

	fmt.Println("Goroutine pools started")
}

func (r *pools) worker() {
	defer func() {
		r.workerCompletedNotification <- true
	}()

	for routine := range r.routines {
		fmt.Println("Goroutine execution started: ", routine.ID)

		routine.ExecuteFunc()

		if routine.Finished != nil {
			routine.Finished <- true
		}

		fmt.Println("Goroutine execution finished: ", routine.ID)
	}
}

func (r *pools) checkCompletedWorker() {
	defer func() {
		close(r.workerCompletedNotification)
		r.allWokerCompleted <- true
	}()

	completedWorkerCount := 0
	for range r.workerCompletedNotification {
		completedWorkerCount += 1
		if completedWorkerCount >= r.poolSize {
			break
		}
	}
}

func (r *pools) Shutdown() {
	r.close()

	fmt.Println("Gracefully shutting down goroutine pools")

	r.await()

	fmt.Println("Goroutine pools shut down")
}

func (r *pools) close() {
	if r.routines == nil {
		return
	}

	fmt.Printf("closing goroutine pools with %d remaining\n", len(r.routines))
	close(r.routines)
}

func (r *pools) await() {
	start := time.Now()

	for time.Since(start) < r.shutdownPeriod {
		select {
		case <-r.allWokerCompleted:
			fmt.Println("all execution finished")
			return
		default:
			fmt.Printf("%d goroutine remaining\n", len(r.routines))
			time.Sleep(1 * time.Second)
		}
	}
}

func (r *pools) Send(routine Routine) error {
	if r.routines == nil {
		return ErrWorkerNotStarted
	}

	select {
	case r.routines <- routine:
		fmt.Printf("goroutine queued: %s\n", routine.ID)
	default:
		return ErrFullQueue
	}

	return nil
}
