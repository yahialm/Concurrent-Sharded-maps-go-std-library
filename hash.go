package main

import "hash/crc32"

func hashToInteger(key string) (uint64, error) {
	hasher := crc32.NewIEEE()
	_, err := hasher.Write([]byte(key))

	if err != nil {
		return 0, err
	}

	hashedKeySlice := hasher.Sum(nil)

	hashedKey := uint64(hashedKeySlice[0]) | uint64(hashedKeySlice[1])<<8 | uint64(hashedKeySlice[2])<<16 | uint64(hashedKeySlice[3])<<24

	return hashedKey, nil
}
