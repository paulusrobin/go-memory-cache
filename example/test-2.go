package example

import (
	"fmt"
	"github.com/paulusrobin/go-memory-cache/memory-cache"
	"sync"
	"time"
)

func Test2() {
	cache, err = memory_cache.NewWithOption(memory_cache.Option{
		MaxEntriesKey: 100,
		OnRemoveWithReason: func(key string, reason string) {
			log.Infof("Key %s is removed because of %s", key, reason)
		},
	})
	if err != nil {
		log.Error(err)
		panic(err)
	}

	key := "test"
	var wg sync.WaitGroup
	iterate := 110
	wg.Add(iterate)
	for i := 1; i <= iterate; i++ {
		go func(i int) {
			_ = cache.Set(fmt.Sprintf("%s#%d", key, i), "example", 3*time.Second)
			time.Sleep(3 * time.Second)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)
}
