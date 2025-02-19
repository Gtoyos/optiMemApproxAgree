package main

//Memory.go: includes the implementation of the shared memory used to communicate between
//the two processes.
//We implement a classical snapshot.

import (
	"sync"
)

// A generic snapshot type. With concurrency controls for non-atomic concurrent
// read opertaions.
type Snapshot[T any] struct {
	raw_array []T
	mus       []sync.RWMutex
}

// Creates an instance of shared memory.
func NewSnapshot[T any](size int) *Snapshot[T] {
	arr := make([]T, size)
	return &Snapshot[T]{raw_array: arr}
}

// Return a consistent copy of the shared array
func (s *Snapshot[T]) Snap() []T {
	snap := make([]T, len(s.raw_array))
	for i := 0; i < len(s.raw_array); i++ {
		s.mus[i].RLock()
		defer s.mus[i].RUnlock()
	}
	copy(snap, s.raw_array)
	return snap
}

// Writes a value in the snapshot. Note that writes are not done concurrently
// for the same entry. Each process has its entry.
func (s *Snapshot[T]) Write(value T, index int) {
	if index >= 0 && index < len(s.raw_array) {
		s.mus[index].Lock()
		s.raw_array[index] = value
		s.mus[index].Unlock()
	}
}

// Allow only numbers that provide consistency on 64 bit systems
type AtomicTypes[T any] interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | bool | *T
}

// A snapshot for primitive atomic operations which do not require a lock on each entry.
type SnapshotAtomic[T AtomicTypes[T]] struct {
	raw_array []T
	mu        sync.RWMutex
}

// Creates an instance of shared memory.
func NewSnapshotAtomic[T AtomicTypes[T]](size int) *SnapshotAtomic[T] {
	arr := make([]T, size)
	return &SnapshotAtomic[T]{raw_array: arr}
}

// Return a consistent copy of the shared array
func (s *SnapshotAtomic[T]) Snap() []T {
	snap := make([]T, len(s.raw_array))
	s.mu.RLock()
	copy(snap, s.raw_array)
	s.mu.RUnlock()
	return snap
}

// Writes a value in the snapshot. Note that writes are not done concurrently
// for the same entry. Each process has its entry.
func (s *SnapshotAtomic[T]) Write(value T, index int) {
	s.mu.Lock()
	if index >= 0 && index < len(s.raw_array) {
		s.raw_array[index] = value
	}
	s.mu.Unlock()
}
