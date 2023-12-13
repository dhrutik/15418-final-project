package lock_free

import (
	// "errors"
	"fmt"
	"main/tree_api"
	// "reflect"
	"slices"
	"sync"
)

func (t *LockFreeTree) PrintLeaf(l *Node, idx int) {
	// for _, l := range leaf {
	if idx > -1 {
		fmt.Printf("Leaf %d: ", idx)
	} else {
		fmt.Printf("Leaf: ")
	}
	for i := 0; i < l.NumKeys; i++ {
		if verbose_output {
			fmt.Printf("%d \n", l.Pointers[i])
		}
		fmt.Printf("%d ", l.Keys[i])
	}
	fmt.Printf("\n")
	// }
}

func (t *LockFreeTree) RedistributeWorkLeaves(index int, sharedLeafData [][]*Node) []*Node {
	if index == 0 {
		return sharedLeafData[index]
	}
	L_i_prime := make([]*Node, 0)
	curr_L_i := sharedLeafData[index]
	for _, lam := range curr_L_i {
		for j := 0; j < index; j++ {
			L_j := sharedLeafData[j]
			if slices.Contains(L_j, lam) { // in L_i, not in any L_j prior
				break
			}
			if !slices.Contains(L_i_prime, lam) { // might be an unnecessary check
				L_i_prime = append(L_i_prime, lam)
			}
			// L_i_prime = append(L_i_prime, lam)
		}
	}
	return L_i_prime
}

func (t *LockFreeTree) Stage2Logic(i int, num_threads int, sharedLeafData [][]*Node, R [][]*tree_api.Record) []*Node /*map[*Node]([][]Modification)*/ {
	// Redistribute Work
	L_i_prime := t.RedistributeWorkLeaves(i, sharedLeafData)
	return L_i_prime
	// Modify leaves independently
}

func (t *LockFreeTree) modifySharedModLists(index int, sharedLeafData [][]*Node, sharedModLists [](map[*Node]([]Modification)), R [][]*tree_api.Record, palmMaxThreadCount int, wg *sync.WaitGroup, testing [][]*Node) {
	defer wg.Done()
	// fmt.Printf("in modsharedlists index: %d\n", index)
	res := t.Stage2Logic(index, palmMaxThreadCount, sharedLeafData, R)
	testing[index] = res
	// sharedModLists[index] = res
}

func (t *LockFreeTree) Stage2(sharedLeafData [][]*Node, palmMaxThreadCount int) ([](map[*Node]([]Modification)), [][]*tree_api.Record) {
	fmt.Printf("Starting Stage 2\n")
	var wg2 sync.WaitGroup
	dbg := true

	testing := make([][]*Node, palmMaxThreadCount)
	for i := 0; i < palmMaxThreadCount; i++ {
		testing[i] = make([]*Node, 0)
	}

	// Set up
	sharedModLists := make([]map[*Node]([]Modification), palmMaxThreadCount)
	for i := 0; i < palmMaxThreadCount; i++ {
		sharedModLists[i] = make(map[*Node]([]Modification))
	}
	R := make([][]*tree_api.Record, palmMaxThreadCount+1) // could potentially have off by 1 issues (just not using 0 i guess)
	for i := 0; i < palmMaxThreadCount; i++ {
		R[i] = make([]*tree_api.Record, 0)
	}

	// Do threads
	for i := 0; i < palmMaxThreadCount; i++ {
		wg2.Add(1) // Increment the counter for each goroutine
		go t.modifySharedModLists(i, sharedLeafData, sharedModLists, R, palmMaxThreadCount, &wg2, testing)
	}
	wg2.Wait()
	fmt.Println("All workers have completed Stage 2.")

	if dbg {
		fmt.Println("Printing testing redistributed leaves vals")
		for idx, L_i := range testing {
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

	return sharedModLists, R
}
