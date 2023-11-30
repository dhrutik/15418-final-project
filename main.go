package main

import (
	"fmt"
	"main/benchmark"
	"main/global_lock_tree"
	"main/seq_tree"
)

func main() {
	keyCount := 2500000
	seqTree := seq_tree.NewTree()
	// globalLockTree := global_lock_tree.NewTree()
	fmt.Println("Sequential Tree Insert Benchmark")
	benchmark.RunInsertBenchmark(seqTree, keyCount, 1)
	fmt.Println("Global Lock Tree Insert Benchmark")
	seqGlobalTree := global_lock_tree.NewTree()
	duration, _ := benchmark.RunInsertBenchmark(seqGlobalTree, keyCount, 1)
	bigTrees := []*global_lock_tree.GlobalLockTree{}
	bigTrees = append(bigTrees, seqGlobalTree.(*global_lock_tree.GlobalLockTree))
	for threadCount := 2; threadCount <= 64; threadCount *= 2 {
		tree := global_lock_tree.NewTree()
		newDuration, _ := benchmark.RunInsertBenchmark(tree, keyCount, threadCount)
		fmt.Printf("Speedup over sequential: %f\n", duration.Seconds()/newDuration.Seconds())
		bigTrees = append(bigTrees, tree.(*global_lock_tree.GlobalLockTree))
	}
	fmt.Println("Sequential Tree Find Benchmark")
	benchmark.RunFindBenchmark(seqTree, keyCount, 1)

	fmt.Println("Global Lock Tree Find Benchmark")
	duration, _ = benchmark.RunFindBenchmark(seqGlobalTree, keyCount, 1)
	treeIdx := 0
	for threadCount := 2; threadCount <= 64; threadCount *= 2 {
		newDuration, _ := benchmark.RunFindBenchmark(bigTrees[treeIdx], keyCount, threadCount)
		fmt.Printf("Speedup over sequential: %f\n", duration.Seconds()/newDuration.Seconds())
		treeIdx++
	}
	treeIdx = 0
	fmt.Println("Sequential Tree Delete Benchmark")
	benchmark.RunDeleteBenchmark(seqTree, keyCount, 1)

	fmt.Println("Global Lock Tree Delete Benchmark")
	duration, _ = benchmark.RunDeleteBenchmark(seqGlobalTree, keyCount, 1)
	for threadCount := 2; threadCount <= 64; threadCount *= 2 {
		newDuration, _ := benchmark.RunDeleteBenchmark(bigTrees[treeIdx], keyCount, threadCount)
		fmt.Printf("Speedup over sequential: %f\n", duration.Seconds()/newDuration.Seconds())
		treeIdx++
	}

}
