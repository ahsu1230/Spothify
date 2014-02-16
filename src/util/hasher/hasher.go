package hasher

import (
		"hash/fnv"
)

// For Testing Purposes:
const MAX_HASH = 100000

// Partitioning:  Defined here so that all implementations
// use the same mechanism.
func Storehash(key string) uint32 {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	//return hasher.Sum32()
	return hasher.Sum32() % MAX_HASH
}
