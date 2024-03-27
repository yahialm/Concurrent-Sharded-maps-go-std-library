package main

import (
	"errors"
	"sync"
)

/*
The sync.Map type is optimized for two common use cases: (1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, or (2) when multiple goroutines read, write, and overwrite entries for disjoint sets of keys. In these two cases, use of a Map may significantly reduce lock contention compared to a Go map paired with a separate Mutex or RWMutex.
*/

type shardedMap struct {
	shards map[int]*shard // Let's suppose keys are integers
	muxSm sync.RWMutex
}

func NewShardedMap() *shardedMap {
	shards := make(map[int]*shard, 128) // Maybe is better to use integers as keys or maybe is better to use a slice of sync.Map : check the performance
	for i:=0; i < 128; i++  {
		shards[i] = NewShard() // fixed
	}
	return &shardedMap{
		shards: shards,
		muxSm: sync.RWMutex{},
	}
}

func (sm *shardedMap) Len() int {
    sm.muxSm.RLock()
    defer sm.muxSm.RUnlock()
    count := 0
    for _, shard := range sm.shards {
        if shard != nil {
            count++
        }
    }
    return count
}

// Get the index of shard based on the key provided 
func (sMap *shardedMap) getShardIndex(key string) (int, error) {

	// Hash the key
	hashedKey, err := hashToInteger(key)

	if err != nil {
		return -1, err
	}

	// Perform modulo operation
	shardIndex := hashedKey % uint64(len(sMap.shards))

	return int(shardIndex), nil
}

func (sMap *shardedMap) get(key string) (any, error){

	// map length is zero
	if len(sMap.shards) == 0 {
		return -1, errors.New("key not found, empty map")
	}

	// Get the index of the responsible shard for this key	
	index, err := sMap.getShardIndex(key)

	if err != nil {
		return -1, err
	}

	keyIsPresent := sMap.shards[index].isPresent(key)

	if !keyIsPresent {
		return -1, errors.New("key not found")
	}

	// use loadValue method to get the value associated with that key from that shard
	value := sMap.shards[index].loadValue(key)

	// check if the value is nil ---> no value is present -----> returns an (nil, error)
	if value == nil {
		return nil, errors.New("no value is assigned to this key")
	}

	// return value, nil
	return value, nil
}

// Returns true if the key-value pair is stored successfully, false if not
func (sMap *shardedMap) store(key string, value any) (bool, error){

	// Hash the key and perform the modulo operation to get the index of the shard where
	// the key value pair will be stored.
	index, err := sMap.getShardIndex(key)

	if err != nil {
		return false, err
	}

	// check for existing some key == key 
	keyIsPresent := sMap.shards[index].isPresent(key)

	if keyIsPresent {
		return false, errors.New("this key already exist, try a new one")
	}

	// store the key-value pair using the returned index from modulo func
	_ , err = sMap.shards[index].storeValue(key, value)

	if err != nil {
		return false, err
	}

	return true, nil

}


func (sMap *shardedMap) delete(key string) error{

	// Hash the key and perform the modulo operation to get the index of the shard where
	// the key value pair will be stored.
	index, err := sMap.getShardIndex(key)

	if err != nil {
		return err
	}

	keyIsPresent := sMap.shards[index].isPresent(key)

	if !keyIsPresent {
		return errors.New("key not found")
	}

	// delete the key-value pair
	err = sMap.shards[index].deleteValue(key)

	if err != nil {
		return err
	}

	return nil
}
