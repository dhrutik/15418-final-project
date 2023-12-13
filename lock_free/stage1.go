package lock_free

import (
	"fmt"
	"main/tree_api"
	"sync"
)

// Evenly distributes queries across all threads, return slice corresponding to ith thread
func (t *LockFreeTree) PartitionInput(Q []tree_api.Query, i int, num_threads int) []tree_api.Query {
	num_queries := len(Q)
	start := i + num_threads - 1 // because will never have 0 threads
	end := start + (num_queries / num_threads)
	// TODO: CHECK THIS MATH, Ensure all queries are in some partition
	res := Q[start:end]
	fmt.Printf("\nIn PartitionInput for %d threads, thread_id: %d\n", num_threads, i)
	fmt.Printf("Queries:\n")
	for _, q := range res {
		fmt.Printf("res: QueryMethod: %d, QueryKey: %d\n", q.Method, q.Key)
	}

	return res
}

// func printNodeSlice(s [](*Node)) {
// 	fmt.Printf("key=%d\n", s.)
// }

func (t *LockFreeTree) FindMultiple(Q []tree_api.Query) [](*Node) {
	res := [](*Node){}
	fmt.Printf("b4 len of res: %d\n", len(res))
	verbose := false // debugging purposes
	for _, q := range Q {
		// TODO: potentially need to remove q from Q once found, or mark as done (serviced query)
		if q.Method == tree_api.MethodFind {
			key := q.Key
			node := t.findLeaf(key, verbose)
			res = append(res, node) // TODO: Check for errors
			q.Done = true           // Ensure this can actually be done in place like this
			fmt.Printf("after len of res: %d\n", len(res))
		} else {
			continue
		}
	}
	return res
}

func (t *LockFreeTree) Stage1Logic(Q []tree_api.Query, i int, num_threads int) []*Node {
	// Stage 1
	Q_i := t.PartitionInput(Q, i, num_threads)
	L_i := t.FindMultiple(Q_i)
	fmt.Printf("stage 1 done\n")
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
	return L_i
	// fmt.Printf("printing tree now...\n")
	// t.PrintTree()

	// defer wg.Done()

	// t.Sync(i, num_threads)
}

func (t *LockFreeTree) modifySharedLeaves(index int, sharedLeafData [][]*Node, queries []tree_api.Query, palmMaxThreadCount int, wg *sync.WaitGroup) {
	defer wg.Done()

	res := t.Stage1Logic(queries, index, palmMaxThreadCount)
	sharedLeafData[index] = res
}

func (t *LockFreeTree) Stage1(queries []tree_api.Query, palmMaxThreadCount int) [][]*Node {
	var wg1 sync.WaitGroup
	// dbg := false
	sharedLeafData := make([][]*Node, palmMaxThreadCount+1)
	for i := 1; i <= palmMaxThreadCount; i++ {
		sharedLeafData[i] = make([]*Node, 0)
	}
	for i := 1; i <= palmMaxThreadCount; i++ {
		wg1.Add(1) // Increment the counter for each goroutine
		go t.modifySharedLeaves(i, sharedLeafData, queries, palmMaxThreadCount, &wg1)
	}
	wg1.Wait()
	fmt.Println("All workers have completed.")

	// if dbg {
	// 	fmt.Println("Printing sharedLeafData vals")
	// 	for idx, L_i := range sharedLeafData {
	// 		fmt.Printf("index: %d\n", idx)
	// 		for _, l := range L_i {
	// 			fmt.Printf("Leaf: ")
	// 			for i := 0; i < l.NumKeys; i++ {
	// 				if verbose_output {
	// 					fmt.Printf("%d \n", l.Pointers[i])
	// 				}
	// 				fmt.Printf("%d ", l.Keys[i])
	// 			}
	// 			fmt.Printf("\n")
	// 		}
	// 	}
	// }
	return sharedLeafData
}