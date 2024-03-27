package main

import (
	"errors"
	"sync"
)

/*
The sync.Map type is optimized for two common use cases: (1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, or (2) when multiple goroutines read, write, and overwrite entries for disjoint sets of keys. In these two cases, use of a Map may significantly reduce lock contention compared to a Go map paired with a separate Mutex or RWMutex.
*/

type shard struct {
	s sync.Map
	muxS sync.Mutex
}

func NewShard() *shard {
	return &shard{
		s: sync.Map{},
		muxS: sync.Mutex{},
	}
}

// Check if a key exists in the shard
func (s *shard) isPresent(key string) bool {
	var found bool
	s.s.Range(func(k, v any) bool {
		if k == key {
			found = true
			return false // Stop iterating
		}
		return true // Continue iterating
	})
	return found
}

func (s *shard) loadValue(key string) (any) {
	value, _ := s.s.Load(key)
	// returns value if present, nil if not
	return value
}

func (s *shard) storeValue(key string, value any) (any, error) {
	actual, ok := s.s.LoadOrStore(key, value)
	// if key-value pair already exists. Load the value associated with key.
	if ok {
		return actual, errors.New("value already exists for this key")
	}
	// Otherwise, we store the new value-pair and returns the value.
	return actual, nil
}

func (s *shard) deleteValue(key string) error {
	_ , ok := s.s.LoadAndDelete(key)
	// If the key doesn't exist
	if !ok {
		return errors.New("key not found")
	}
	// Otherwise
	return nil
}
