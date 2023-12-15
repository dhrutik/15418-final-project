package benchmark

import (
	"fmt"
	"main/lock_free"
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

func makeInsertQueryLists(totalKeyCount, perRoundKeyCount int) [][]tree_api.Query {
	res := make([][]tree_api.Query, 0)
	for i := 0; i < totalKeyCount; i += perRoundKeyCount {
		queries := make([]tree_api.Query, 0)
		for j := i; j < i+perRoundKeyCount; j++ {
			queries = append(queries,
				tree_api.Query{
					Method:  tree_api.MethodInsert,
					Key:     j,
					Pointer: &tree_api.Record{Value: []byte("value")},
				},
			)
		}
		res = append(res, queries)
	}
	return res
}

func makeFindQueryLists(totalKeyCount, perRoundKeyCount int) [][]tree_api.Query {
	res := make([][]tree_api.Query, 0)
	for i := 0; i < totalKeyCount; i += perRoundKeyCount {
		queries := make([]tree_api.Query, 0)
		for j := i; j < i+perRoundKeyCount; j++ {
			queries = append(queries,
				tree_api.Query{
					Method: tree_api.MethodFind,
					Key:    j,
				},
			)
		}
		res = append(res, queries)
	}
	return res
}

func makeDeleteQueryLists(totalKeyCount, perRoundKeyCount int) [][]tree_api.Query {
	res := make([][]tree_api.Query, 0)
	for i := 0; i < totalKeyCount; i += perRoundKeyCount {
		queries := make([]tree_api.Query, 0)
		for j := i; j < i+perRoundKeyCount; j++ {
			queries = append(queries,
				tree_api.Query{
					Method: tree_api.MethodDelete,
					Key:    j,
				},
			)
		}
		res = append(res, queries)
	}
	return res
}

func PalmInsertBenchmark(tree *lock_free.LockFreeTree, totalKeyCount, perRoundKeyCount, threadCount int) (time.Duration, float64) {
	queryLists := makeInsertQueryLists(totalKeyCount, perRoundKeyCount)
	startTime := time.Now()
	for i := 0; i < len(queryLists); i++ {
		queries := queryLists[i]
		tree.Palm(queries, threadCount)
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(totalKeyCount) / elapsedTime.Seconds()
	fmt.Printf("Insert %d keys in %f seconds with %d threads, throughput: %f keys/s\n", totalKeyCount, elapsedTime.Seconds(), threadCount, throughput)
	return elapsedTime, throughput
}

func PalmFindBenchmark(tree *lock_free.LockFreeTree, totalKeyCount, perRoundKeyCount, threadCount int) (time.Duration, float64) {
	queryLists := makeFindQueryLists(totalKeyCount, perRoundKeyCount)
	startTime := time.Now()
	for i := 0; i < len(queryLists); i++ {
		queries := queryLists[i]
		tree.Palm(queries, threadCount)
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(totalKeyCount) / elapsedTime.Seconds()
	fmt.Printf("Find %d keys in %f seconds with %d threads, throughput: %f keys/s\n", totalKeyCount, elapsedTime.Seconds(), threadCount, throughput)
	return elapsedTime, throughput
}

func PalmDeleteBenchmark(tree *lock_free.LockFreeTree, totalKeyCount, perRoundKeyCount, threadCount int) (time.Duration, float64) {
	queryLists := makeDeleteQueryLists(totalKeyCount, perRoundKeyCount)
	startTime := time.Now()
	for i := 0; i < len(queryLists); i++ {
		queries := queryLists[i]
		tree.Palm(queries, threadCount)
	}
	elapsedTime := time.Since(startTime)
	throughput := float64(totalKeyCount) / elapsedTime.Seconds()
	fmt.Printf("Find %d keys in %f seconds with %d threads, throughput: %f keys/s\n", totalKeyCount, elapsedTime.Seconds(), threadCount, throughput)
	return elapsedTime, throughput
}

// func RunStage1(tree tree_api.BPTree, queries []tree_api.Query, numKeys int, threads int, wg *sync.WaitGroup) (time.Duration, float64) {
// 	// ASSERT tree is of type lock_free
// 	done := make(chan bool)
// 	numKeysPerThread := (numKeys + threads - 1) / threads
// 	// keys := makeShuffledKeysList(numKeys)
// 	startTime := time.Now()
// 	for i := 0; i < threads; i++ {
// 		go func(index int) {
// 			startIndex := index * numKeysPerThread
// 			endIndex := startIndex + numKeysPerThread
// 			if endIndex > numKeys {
// 				endIndex = numKeys
// 			}
// 			if startIndex > numKeys {
// 				startIndex = numKeys
// 			}
// 			tree.Stage1(queries, index, threads)
// 			// findKeys(tree, startIndex, endIndex, index, keys)
// 			done <- true
// 		}(i)
// 	}
// 	for doneCount := 0; doneCount < threads; doneCount++ {
// 		<-done
// 	}
// 	elapsedTime := time.Since(startTime)
// 	throughput := float64(numKeys) / elapsedTime.Seconds()
// 	fmt.Printf("STAGE 1: Find %d keys in %f seconds with %d threads, throughput: %f keys/s\n", numKeys, elapsedTime.Seconds(), threads, throughput)
// 	return elapsedTime, throughput

// }
