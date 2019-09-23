# Go-Memory-Cache

### Installation
```
~ $ go get -u github.com/paulusrobin/go-memory-cache
``` 

#### Example
```go
cache, err := memory_cache.New()
if err != nil {
    log.Error(err)
    panic(err)
}

_ = cache.Set("example", "example", time.Minute)
if val, err := cache.Get("example"); err == nil {
    fmt.Println(val)
}
```

#### Options
```
# Max entry size in bytes
MaxEntrySize       int

# Max number of key invoked to memory, will evict old keys
MaxEntriesKey      int

# Max memory used in bytes, will evict old keys
MaxEntriesInWindow int

# Callback Function when removing data from memory
OnRemove           func(key string, value interface{})

# Callback Function when removing data from memory, called when there is special reason (MaxEntriesKey or MaxEntriesInWindow)
OnRemoveWithReason func(key string, reason string)
```

#### Example using Option
```go
cache, err := memory_cache.NewWithOption(memory_cache.Option{
    MaxEntrySize:       1024,
    MaxEntriesKey:      100,
    MaxEntriesInWindow: 1024 * 1024,
    OnRemove:           func(key string, value interface{}) {
        fmt.Printf("Key %s is removed\n", key)
    },
    OnRemoveWithReason: func(key string, reason string) {
        fmt.Printf("Key %s is removed because of %s\n", key, reason)
    },
})
if err != nil {
    log.Error(err)
    panic(err)
}

_ = cache.Set("example", "example", memory_cache.Forever)
if val, err := cache.Get("example"); err == nil {
    fmt.Println(val)
}
```

### Performance Benchmarks
```
goos: darwin
goarch: amd64

BenchmarkRead-4                  1000000              1132 ns/op
BenchmarkWrite-4                 1000000              1060 ns/op
BenchmarkWriteRead-4             1000000              1239 ns/op
BenchmarkWriteRemove-4           2000000               776 ns/op
BenchmarkWriteLen-4              2000000              1119 ns/op
BenchmarkWriteSize-4             2000000              1097 ns/op
BenchmarkWriteTruncate-4         2000000               617 ns/op
```
