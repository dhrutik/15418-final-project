package benchmark

import (
	"fmt"
	"main/tree_api"
	"math/rand"
	"time"
)

func insertKeys(tree tree_api.BPTree, startIndex, endIndex int, threadId int, keys []int) {
	for i := startIndex; i < endIndex; i++ {
		tree.Insert(keys[i], []byte("value"))
	}
}

func deleteKeys(tree tree_api.BPTree, startIndex, endIndex int, threadId int, keys []int) {
	for i := startIndex; i < endIndex; i++ {
		tree.Delete(keys[i])
	}
}

func findKeys(tree tree_api.BPTree, startIndex, endIndex int, threadId int, keys []int) {
	for i := startIndex; i < endIndex; i++ {
		tree.Find(keys[i], false)
	}
}

func makeShuffledKeysList(numKeys int) []int {
	keys := make([]int, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = i
	}
	// Shuffle keys
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	return keys
}

func RunInsertBenchmark(tree tree_api.BPTree, numKeys int, threads int) (time.Duration, float64) {
	done := make(chan bool)
	numKeysPerThread := numKeys / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			insertKeys(tree, startIndex, endIndex, index, keys)
			done <- true
		}(i)
	}
	for doneCount := 0; doneCount < threads; doneCount++ {
		<-done
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(numKeys) / elapsedTime.Seconds()
	fmt.Printf("Insert %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
	return elapsedTime, throughput
}

func RunFindBenchmark(tree tree_api.BPTree, numKeys int, threads int) (time.Duration, float64) {
	done := make(chan bool)
	numKeysPerThread := numKeys / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			findKeys(tree, startIndex, endIndex, index, keys)
			done <- true
		}(i)
	}
	for doneCount := 0; doneCount < threads; doneCount++ {
		<-done
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(numKeys) / elapsedTime.Seconds()
	fmt.Printf("Find %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
	return elapsedTime, throughput
}

func RunDeleteBenchmark(tree tree_api.BPTree, numKeys int, threads int) (time.Duration, float64) {
	done := make(chan bool)
	numKeysPerThread := numKeys / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			deleteKeys(tree, startIndex, endIndex, index, keys)
			done <- true
		}(i)
	}
	for doneCount := 0; doneCount < threads; doneCount++ {
		<-done
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(numKeys) / elapsedTime.Seconds()
	fmt.Printf("Delete %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
	return elapsedTime, throughput
}
