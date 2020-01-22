package memory_cache

import (
	"fmt"
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/paulusrobin/go-memory-cache/logs"
	"github.com/pkg/errors"
	"reflect"
	"sync"
	"time"
)

var tooManyKeysInvoked = "Too many keys invoked"
var windowsTooLarge = "Windows too large"
var valueTooLarge = "Entry value too large"
var memoryExceed = "Memory exceed"

type cache struct {
	option Option
	mu     sync.Mutex
	log    logs.Logger
	data   map[string]interface{}
	queue  []string
	size   uintptr
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newItemSize := c.getSize(value)
	if newItemSize > uintptr(c.option.MaxEntrySize) {
		return errors.New(valueTooLarge)
	}

	for {
		if c.size+newItemSize > uintptr(c.option.MaxEntriesInWindow) {
			c.forceRemove(c.queue[0], windowsTooLarge)
		} else {
			break
		}
	}
	if len(c.data) >= c.option.MaxEntriesKey {
		c.forceRemove(c.queue[0], tooManyKeysInvoked)
	}

	c.data[key] = value
	c.queue = append(c.queue, key)
	c.size += newItemSize

	time.AfterFunc(ttl, func() { _ = c.Remove(key) })
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value, ok := c.data[key]; ok {
		return value, nil
	}
	return nil, errors.New(fmt.Sprintf("key %s is not found", key))
}

func (c *cache) Remove(key string) error {
	data, err := c.Get(key)
	if err != nil {
		return errors.WithStack(err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	c.removeQueue()

	if c.option.OnRemove != nil {
		c.option.OnRemove(key, data)
	}
	return nil
}

func (c *cache) Truncate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.data {
		delete(c.data, key)
		c.removeQueue()
	}
	return nil
}

func (c *cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *cache) Size() uintptr {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}

func (c *cache) forceRemove(key string, reason string) {
	delete(c.data, key)
	c.removeQueue()

	if c.option.OnRemoveWithReason != nil {
		c.option.OnRemoveWithReason(key, reason)
	}
}

func (c *cache) removeQueue() {
	if len(c.queue) > 1 {
		sizeItemToRemove := c.getSize(c.queue[0])
		c.queue = c.queue[1:len(c.queue)]
		c.size -= sizeItemToRemove
	} else {
		c.queue = make([]string, 0)
		c.size = 0
	}
}

func (c *cache) getSize(T interface{}) uintptr {
	return reflect.TypeOf(T).Size()
}

func (c *cache) Cleaner(duration time.Duration, done <-chan bool) {
	tick := time.Tick(duration)
	for {
		select {
		case <-tick:
			mem := sigar.Mem{}

			if err := mem.Get(); err != nil {
				c.log.Errorf("error get memory = ", err.Error())
				break
			}

			percentageUsed := float64(mem.ActualUsed) / float64(mem.Total) * 100
			if percentageUsed > c.option.MaxPercentageMemory {
				if c.option.OnMemoryExceed != nil {
					c.option.OnMemoryExceed(percentageUsed, c.option.MaxPercentageMemory, float64(mem.ActualUsed))
				}
				c.forceRemove(c.queue[0], memoryExceed)
			}

		case <-done:
			return
		}
	}
}
