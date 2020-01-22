package memory_cache

import (
	"github.com/paulusrobin/go-memory-cache/logs"
	"math"
	"time"
)

const Infinite = math.MaxInt32
const Forever = 87660 * time.Hour

type (
	Option struct {
		MaxEntrySize        int
		MaxEntriesKey       int
		MaxEntriesInWindow  int
		MaxPercentageMemory float64
		OnRemove            func(key string, value interface{})
		OnRemoveWithReason  func(key string, reason string)
		OnMemoryExceed      func(memoryUsedPercentage float64, maxMemoryPercentage float64, memoryUsed float64)
	}
	Cache interface {
		Set(key string, value interface{}, ttl time.Duration) error
		Get(key string) (interface{}, error)
		Remove(key string) error
		Truncate() error
		Len() int
		Size() uintptr
		Cleaner(time.Duration, <-chan bool)
	}
)

func initializeOption(option Option) Option {
	if option.MaxEntriesKey == 0 {
		option.MaxEntriesKey = Infinite
	}
	if option.MaxEntriesInWindow == 0 {
		option.MaxEntriesInWindow = 2 * 1024 * 1024 * 1024
	}
	if option.MaxEntrySize == 0 {
		option.MaxEntrySize = 1024 * 1024
	}
	if option.MaxPercentageMemory == 0 {
		option.MaxPercentageMemory = 95
	}
	return option
}

func NewWithOption(option Option) (Cache, error) {
	return &cache{
		option: initializeOption(option),
		log:    logs.DefaultLog(),
		data:   make(map[string]interface{}),
		size:   uintptr(0),
	}, nil
}

func New() (Cache, error) {
	return NewWithOption(Option{})
}
