# Scheduler

```go
s := scheduler.New()

s.Every(time.Second, func(time.Duration) {
  fmt.Println("Hello every second")
})

s.Every(100*time.Millisecond, func(time.Duration) {
  fmt.Print(".")
})

s.Start()
```


