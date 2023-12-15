package main

import (
	"fmt"
	"main/benchmark"
	"main/crab"
	"main/global_lock_tree"
	"main/lock_free"
	"main/seq_tree"
	"main/tree_api"
	"time"
)

func runSpeedup(trees []tree_api.BPTree, keyCount int, threadCount int, benchmarkFunc func(tree_api.BPTree, int, int) (time.Duration, float64)) {
	duration, _ := benchmarkFunc(trees[0], keyCount, 1)
	for i := 1; i < len(trees); i++ {
		threadCount := 1 << i // 2 ** i
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

func main() {

	// Set up trees
	keyCount := 1000000
	seqTree := seq_tree.NewTree()
	maxThreadCount := 128
	globalLockTrees := makeTreeList(maxThreadCount, global_lock_tree.NewTree)
	crabTrees := makeTreeList(maxThreadCount, crab.NewTree)

	FLAG_run_benchmarks := true
	FLAG_test_palm := true

	if FLAG_run_benchmarks {
		runBenchmark("Insert", benchmark.RunInsertBenchmark, seqTree, globalLockTrees, crabTrees, keyCount, maxThreadCount)
		runBenchmark("Find", benchmark.RunFindBenchmark, seqTree, globalLockTrees, crabTrees, keyCount, maxThreadCount)
		runBenchmark("Delete", benchmark.RunDeleteBenchmark, seqTree, globalLockTrees, crabTrees, keyCount, maxThreadCount)
	}

	if FLAG_test_palm {
		palmTotalKeyCount := 1000000
		palmKeyCount := 45000
		lock_free_tree := lock_free.NewTree()
		palmMaxThreadCount := 64
		benchmark.InsertQueries(lock_free_tree, palmKeyCount, palmMaxThreadCount)
		benchmark.PalmInsertBenchmark(lock_free_tree.(*lock_free.LockFreeTree), palmTotalKeyCount, palmKeyCount, palmMaxThreadCount)
		benchmark.PalmFindBenchmark(lock_free_tree.(*lock_free.LockFreeTree), palmTotalKeyCount, palmKeyCount, palmMaxThreadCount)
		benchmark.PalmDeleteBenchmark(lock_free_tree.(*lock_free.LockFreeTree), palmTotalKeyCount, palmKeyCount, palmMaxThreadCount)

	}
}
