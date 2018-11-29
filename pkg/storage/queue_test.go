package storage

import (
	"testing"
)

func BenchmarkQueue_Push(b *testing.B) {
	queue := NewQueue()
	for i := 0; i < b.N; i++ {
		queue.Push([]byte("hello,world"))
	}
}

func BenchmarkQueue_Pop(b *testing.B) {
	queue := NewQueue()
	for i := 0; i < b.N; i++ {
		queue.Push([]byte("hello,world"))
	}
}
