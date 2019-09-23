package main

import (
	"fmt"
	"github.com/paulusrobin/Go-Memory-Cache/memory-cache"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	cache, _ := memory_cache.New()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(fmt.Sprintf("%d", i))
	}
}

func BenchmarkWrite(b *testing.B) {
	cache, _ := memory_cache.New()
	for i := 0; i < b.N; i++ {
		_ = cache.Set(fmt.Sprintf("%d", i), i, memory_cache.Forever)
	}
}

func BenchmarkWriteRead(b *testing.B) {
	cache, _ := memory_cache.New()
	for i := 0; i < b.N; i++ {
		_ = cache.Set(fmt.Sprintf("%d", i), i, memory_cache.Forever)
		_, _ = cache.Get(fmt.Sprintf("%d", i))
	}
}

func BenchmarkWriteRemove(b *testing.B) {
	cache, _ := memory_cache.New()
	for i := 0; i < b.N; i++ {
		_ = cache.Set(fmt.Sprintf("%d", i), i, memory_cache.Forever)
		_ = cache.Remove(fmt.Sprintf("%d", i))
	}
}

func BenchmarkWriteLen(b *testing.B) {
	cache, _ := memory_cache.New()

	for i := 0; i < b.N; i++ {
		_ = cache.Set(fmt.Sprintf("%d", i), i, memory_cache.Forever)
		_ = cache.Len()
	}
}

func BenchmarkWriteSize(b *testing.B) {
	cache, _ := memory_cache.New()

	for i := 0; i < b.N; i++ {
		_ = cache.Set(fmt.Sprintf("%d", i), i, memory_cache.Forever)
		_ = cache.Size()
	}
}

func BenchmarkWriteTruncate(b *testing.B) {
	cache, _ := memory_cache.New()

	for i := 0; i < b.N; i++ {
		_ = cache.Set(fmt.Sprintf("%d", i), i, memory_cache.Forever)
		_ = cache.Truncate()
	}
}

