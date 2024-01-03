package pools

import (
	"errors"
	"fmt"
	"time"
)

type Pools interface {
	Start()
	Shutdown()
	Send(j Job) error
}

func New(maxPoolSize int, shutdownPeriod time.Duration, cooldownPerExecutionPeriod time.Duration) Pools {
	return &routine{
		maxPoolSize:                maxPoolSize,
		shutdownPeriod:             shutdownPeriod,
		cooldownPerExecutionPeriod: cooldownPerExecutionPeriod,
	}
}

type routine struct {
	maxPoolSize                int
	shutdownPeriod             time.Duration
	cooldownPerExecutionPeriod time.Duration
	jobs                       chan Job
}

type Job struct {
	ID          string
	ExecuteFunc func()
}

var (
	ErrWorkerNotStarted = errors.New("worker has not started yet")
	ErrPoolsUnavailable = errors.New("routine pools are unavailable")
)

func (r *routine) Start() {
	fmt.Println("Starting goroutine pools")

	r.jobs = make(chan Job, r.maxPoolSize)
	go r.worker()

	fmt.Println("Goroutine pools started")
}

func (r *routine) worker() {
	for job := range r.jobs {
		fmt.Println("execution started: ", job.ID)

		job.ExecuteFunc()

		fmt.Println("execution finished: ", job.ID)
	}
}

func (r *routine) Shutdown() {
	r.close()

	fmt.Println("Gracefully shutting down goroutine pools")

	r.await()

	fmt.Println("Goroutine pools shut down")
}

func (r *routine) await() {
	done := make(chan bool)

	go func() {
		start := time.Now()

		isTimedOut := func() bool {
			return time.Since(start) > r.shutdownPeriod
		}

		for len(r.jobs) > 0 && !isTimedOut() {
			fmt.Printf("%d goroutine remaining\n", len(r.jobs))
			time.Sleep(1 * time.Second)
		}

		done <- true
	}()

	<-done
}

func (r *routine) close() {
	if r.jobs == nil {
		return
	}

	close(r.jobs)
	fmt.Printf("closing goroutine pools with %d remaining\n", len(r.jobs))
}

func (r *routine) Send(job Job) error {
	if r.jobs == nil {
		return ErrWorkerNotStarted
	}

	select {
	case r.jobs <- job:
		fmt.Printf("job queued: %s\n", job.ID)
	default:
		return ErrPoolsUnavailable
	}

	return nil
}
