package main

import (	"fmt"
		"flag"
		"hash/fnv"
)

// Run Example:
// go run hash_string.go ahsu
// go run hash_string.go jfan89


// Partitioning:  Defined here so that all implementations
// use the same mechanism.
func Storehash(key string) uint32 {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return hasher.Sum32()
}

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	targetStr := flag.Arg(0)
	if targetStr == "" {
		fmt.Println("Empty String!")
		return
	}
	fmt.Println("Target String:", targetStr)
	
	hashedKey := Storehash(targetStr)
	fmt.Println("Hash Key:", hashedKey)
}
