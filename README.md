# Goroutine Pools
Managed goroutine with pooling systems

# Example
Simply initialize Goroutine pools instance and shutdown on application exit.

## Starting Pools
```go
const poolSize = 100
const queueSize = 100
const shutdownPeriod = 100 * time.Second
const cooldownPerExecutionPeriod = 10 * time.Millisecond

routinePools := pools.New(poolSize, queueSize, shutdownPeriod, cooldownPerExecutionPeriod)
routinePools.Start()
```

## Fire and Forget Execution
```go
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
```

## Awaited Execution
```go
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
```

## Shutting Down Pools
This will simply await all goroutine execution to finish or until the shutdown period is reached
```
routinePools.Shutdown()
```
