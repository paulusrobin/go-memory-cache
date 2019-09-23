package example

import (
	"fmt"
	"github.com/paulusrobin/go-memory-cache/memory-cache"
	"sync"
)

func Test4() {
	cache, err = memory_cache.NewWithOption(memory_cache.Option{
		MaxEntrySize: 8,
		MaxEntriesInWindow: 1024 * 1024,
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
	iterate := 65537
	wg.Add(iterate)
	for i := 1; i <= iterate; i++ {
		go func(i int) {
			if err := cache.Set(fmt.Sprintf("%s#%d", key, i), i, memory_cache.Forever); err != nil {
				log.Error(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	_ = cache.Truncate()
	log.Info(cache.Size())
}
