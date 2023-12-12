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
		// insertedKeys[i] = true
	}
}

func deleteKeys(tree tree_api.BPTree, startIndex, endIndex int, threadId int, keys []int) {
	for i := startIndex; i < endIndex; i++ {
		tree.Delete(keys[i])
	}
}

func findKeys(tree tree_api.BPTree, startIndex, endIndex int, threadId int, keys []int) {
	for i := startIndex; i < endIndex; i++ {
		_, err := tree.Find(keys[i], false)
		if err != nil {
			panic(err)
		}
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

// func checkAllInserted(insertedKeys []bool) {
// 	for _, inserted := range insertedKeys {
// 		if !inserted {
// 			panic("Not all keys were inserted")
// 		}
// 	}
// }

func RunInsertBenchmark(tree tree_api.BPTree, numKeys int, threads int) (time.Duration, float64) {
	done := make(chan bool)
	numKeysPerThread := (numKeys + threads - 1) / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	// insertedKeys := make([]bool, numKeys)
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			if startIndex > numKeys {
				startIndex = numKeys
			}
			insertKeys(tree, startIndex, endIndex, index, keys)
			done <- true
		}(i)
	}
	for doneCount := 0; doneCount < threads; doneCount++ {
		<-done
	}
	// checkAllInserted(insertedKeys)
	elapsedTime := time.Since(startTime)
	throughput := float64(numKeys) / elapsedTime.Seconds()
	fmt.Printf("Insert %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
	return elapsedTime, throughput
}

func RunFindBenchmark(tree tree_api.BPTree, numKeys int, threads int) (time.Duration, float64) {
	done := make(chan bool)
	numKeysPerThread := (numKeys + threads - 1) / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			if startIndex > numKeys {
				startIndex = numKeys
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
	numKeysPerThread := (numKeys + threads - 1) / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			if startIndex > numKeys {
				startIndex = numKeys
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

/*************************Basic Lock-Free Tests******************************************************/
func InsertQueries(tree tree_api.BPTree, numKeys int, threads int) (time.Duration, float64) {
	done := make(chan bool)
	numKeysPerThread := (numKeys + threads - 1) / threads
	keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	// insertedKeys := make([]bool, numKeys)
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			if startIndex > numKeys {
				startIndex = numKeys
			}
			insertKeys(tree, startIndex, endIndex, index, keys)
			done <- true
		}(i)
	}
	for doneCount := 0; doneCount < threads; doneCount++ {
		<-done
	}
	// checkAllInserted(insertedKeys)
	elapsedTime := time.Since(startTime)
	throughput := float64(numKeys) / elapsedTime.Seconds()
	fmt.Printf("Insert %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
	return elapsedTime, throughput
}

func RunStage1(tree tree_api.BPTree, queries []tree_api.Query, numKeys int, threads int) (time.Duration, float64) {
	// ASSERT tree is of type lock_free
	done := make(chan bool)
	numKeysPerThread := (numKeys + threads - 1) / threads
	// keys := makeShuffledKeysList(numKeys)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go func(index int) {
			startIndex := index * numKeysPerThread
			endIndex := startIndex + numKeysPerThread
			if endIndex > numKeys {
				endIndex = numKeys
			}
			if startIndex > numKeys {
				startIndex = numKeys
			}
			tree.Stage1(queries, index, threads)
			// findKeys(tree, startIndex, endIndex, index, keys)
			done <- true
		}(i)
	}
	for doneCount := 0; doneCount < threads; doneCount++ {
		<-done
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(numKeys) / elapsedTime.Seconds()
	fmt.Printf("STAGE 1: Find %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
	return elapsedTime, throughput

}
