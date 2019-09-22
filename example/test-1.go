package example

import (
	"fmt"
	"github.com/tiketrobin/Go-Memory-Cache/logs"
	"github.com/tiketrobin/Go-Memory-Cache/memory-cache"
	"sync"
	"time"
)

var log logs.Logger
var cache memory_cache.Cache
var err error

type counter struct {
	count int
	mu    sync.Mutex
}

func (c *counter) Add(val int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = c.count + val
}

func init() {
	log = logs.DefaultLog()
}

func getData(key string) bool {
	if _, err := cache.Get(key); err != nil {
		return false
	} else {
		return true
	}
}

func Test1() {
	cache, err = memory_cache.NewWithOption(memory_cache.Option{})
	if err != nil {
		log.Error(err)
		panic(err)
	}

	count := counter{}
	key := "test"

	var wg sync.WaitGroup
	iterate := 1
	wg.Add(iterate)
	for i := 1; i <= iterate; i++ {
		go func(i int) {
			_ = cache.Set(fmt.Sprintf("%s#%d", key, i), "example", 3*time.Second)
			getData(key)
			time.Sleep(3 * time.Second)
			if getData(key) {
				count.Add(-1)
			} else {
				count.Add(1)
			}
			log.Infof("Len: %d", cache.Len())
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)
	log.Infof("Len: %d", cache.Len())
	log.Info(count.count)
}
