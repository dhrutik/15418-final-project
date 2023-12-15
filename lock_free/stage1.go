package lock_free

import (
	"fmt"
	"main/tree_api"
	"slices"
	"sync"
)

// Evenly distributes queries across all threads, return slice corresponding to ith thread
func (t *LockFreeTree) PartitionInput(Q []tree_api.Query, i int, num_threads int) []tree_api.Query {
	num_queries := len(Q)
	start := i * (num_queries / num_threads) // because will never have 0 threads
	end := start + (num_queries / num_threads)
	res := Q[start:end]

	return res
}

func (t *LockFreeTree) FindMultiple(Q []tree_api.Query) [](*Node) {
	res := [](*Node){}
	verbose := false // debugging purposes
	for _, q := range Q {
		key := q.Key
		node := t.findLeaf(key, verbose)
		if !slices.Contains(res, node) {
			res = append(res, node)
		}
	}
	return res
}

func (t *LockFreeTree) Stage1Logic(Q []tree_api.Query, i int, num_threads int) []*Node {
	// Stage 1
	Q_i := t.PartitionInput(Q, i, num_threads)
	L_i := t.FindMultiple(Q_i)
	return L_i
}

func (t *LockFreeTree) modifySharedLeaves(index int, sharedLeafData [][]*Node, queries []tree_api.Query, palmMaxThreadCount int, wg *sync.WaitGroup) {
	defer wg.Done()

	res := t.Stage1Logic(queries, index, palmMaxThreadCount)
	sharedLeafData[index] = res
}

func (t *LockFreeTree) Stage1(queries []tree_api.Query, palmMaxThreadCount int) [][]*Node {
	var wg1 sync.WaitGroup
	dbg := false
	sharedLeafData := make([][]*Node, palmMaxThreadCount)
	for i := 0; i < palmMaxThreadCount; i++ {
		sharedLeafData[i] = make([]*Node, 0)
	}
	for i := 0; i < palmMaxThreadCount; i++ {
		wg1.Add(1) // Increment the counter for each goroutine
		go t.modifySharedLeaves(i, sharedLeafData, queries, palmMaxThreadCount, &wg1)
	}
	wg1.Wait()

	if dbg {
		fmt.Println("Printing sharedLeafData vals")
		for idx, L_i := range sharedLeafData {
			fmt.Printf("index: %d\n", idx)
			for _, l := range L_i {
				fmt.Printf("Leaf: ")
				for i := 0; i < l.NumKeys; i++ {
					if verbose_output {
						fmt.Printf("%d \n", l.Pointers[i])
					}
					fmt.Printf("%d ", l.Keys[i])
				}
				fmt.Printf("\n")
			}
		}
	}
	return sharedLeafData
}
