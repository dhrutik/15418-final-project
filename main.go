package main

import (
	"fmt"
	"main/benchmark"
	"main/crab"
	"main/global_lock_tree"
	"main/lock_free"
	"main/seq_tree"
	"main/tree_api"
	"sync"
	"time"
)

func runSpeedup(trees []tree_api.BPTree, keyCount int, threadCount int, benchmarkFunc func(tree_api.BPTree, int, int) (time.Duration, float64)) {
	duration, _ := benchmarkFunc(trees[0], keyCount, 1)
	for i := 1; i < len(trees); i++ {
		threadCount := 1 << (i + 1) // 2 ** (i + 1)
		newDuration, _ := benchmarkFunc(trees[i], keyCount, threadCount)
		fmt.Printf("Speedup over sequential: %f\n", duration.Seconds()/newDuration.Seconds())
	}
}

func runBenchmark(benchmarkName string, benchmarkFunc func(tree_api.BPTree, int, int) (time.Duration, float64), seqTree tree_api.BPTree, globalLockTrees []tree_api.BPTree, crabTrees []tree_api.BPTree, keyCount int, maxThreadCount int) {
	fmt.Printf("Benchmark %s\n", benchmarkName)
	fmt.Printf("Sequential Tree %s Benchmark\n", benchmarkName)
	benchmarkFunc(seqTree, keyCount, 1)
	fmt.Printf("Global Lock Tree %s Benchmark\n", benchmarkName)
	runSpeedup(globalLockTrees, keyCount, maxThreadCount, benchmarkFunc)
	fmt.Printf("Crab Tree %s Benchmark\n", benchmarkName)
	runSpeedup(crabTrees, keyCount, maxThreadCount, benchmarkFunc)
}

func makeTreeList(threadCount int, treeConstructor func() tree_api.BPTree) []tree_api.BPTree {
	trees := []tree_api.BPTree{}
	for i := threadCount; i > 0; i /= 2 {
		trees = append(trees, treeConstructor())
	}
	return trees
}

func modifySharedLeaves(index int, sharedLeafData [][]Node*, queries 
						[]tree_api.Query, palmMaxThreadCount int, wg *sync.WaitGroup) {
	defer wg.Done()

	res := lock_free_tree.Stage1(queries, i, palmMaxThreadCount)
	sharedArray[index] = res
}

func main() {

	// Set up trees
	keyCount := 1000000
	seqTree := seq_tree.NewTree()
	maxThreadCount := 64
	globalLockTrees := makeTreeList(maxThreadCount, global_lock_tree.NewTree)
	crabTrees := makeTreeList(maxThreadCount, crab.NewTree)

	FLAG_run_benchmarks := false
	FLAG_test_palm := true

	if FLAG_run_benchmarks {
		runBenchmark("Insert", benchmark.RunInsertBenchmark, seqTree, globalLockTrees, crabTrees, keyCount, maxThreadCount)
		runBenchmark("Find", benchmark.RunFindBenchmark, seqTree, globalLockTrees, crabTrees, keyCount, maxThreadCount)
		runBenchmark("Delete", benchmark.RunDeleteBenchmark, seqTree, globalLockTrees, crabTrees, keyCount, maxThreadCount)
	}

	if FLAG_test_palm {
		palmKeyCount := 10
		lock_free_tree := lock_free.NewTree()
		palmMaxThreadCount := 3
		var wg1 sync.WaitGroup

		// Construct Tree
		benchmark.InsertQueries(lock_free_tree, palmKeyCount, palmMaxThreadCount)
		queries := make([]tree_api.Query, 0)
		for i := 0; i < palmKeyCount; i++ {
			queries = append(queries, tree_api.Query{tree_api.MethodFind, i, false})
		}

		// Run Stage 1
		sharedLeafData := make([][]Node*, 1)
		for i := 0; i < palmMaxThreadCount; i++ {
			sharedLeafData[i] = make([]Node*, 0)
		}
		for i := 1; i <= palmMaxThreadCount; i++ {
			wg1.Add(1) // Increment the counter for each goroutine
			go modifySharedLeaves(i, sharedLeafData, queries, palmMaxThreadCount, &wg1)
			// go lock_free_tree.Stage1(queries, i, palmMaxThreadCount, &wg1)
		}
		// Sync
		wg1.Wait()
		fmt.Println("All workers have completed.")
		fmt.Println("Printing sharedLeafData vals")
		for _, L_i := range sharedLeafData {
			for _, l := range L_i {
				fmt.Printf("Leaf: ")
				for i = 0; i < l.NumKeys; i++ {
					if verbose_output {
						fmt.Printf("%d \n", l.Pointers[i])
					}
					fmt.Printf("%d ", l.Keys[i])
				}
				fmt.Printf("\n")
			}
		}
		

		// We are assured that the results reflect the state of the tree
		// when each query was dispatched, because no modifications to the tree have occurred yet.

		// Run Stage 2
		// var wg2 sync.WaitGroup

		// Stage 1 Test
		// benchmark.RunStage1(lock_free_tree, queries, palmKeyCount, palmMaxThreadCount, &wg1)
	}
}
