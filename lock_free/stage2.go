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

func assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
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

func (t *LockFreeTree) ResolveHazards(L_i_prime []*Node, queries []tree_api.Query) ([]*tree_api.Record, map[*Node]([]tree_api.Query)) {
	res := make([]*tree_api.Record, 0)
	findQueries := make([]tree_api.Query, 0)
	O_L_i := make(map[*Node]([]tree_api.Query))

	// Extract queries relevant to this (index-th) thread
	for _, q := range queries {
		// Iterate over *my* leaves
		for _, node := range L_i_prime {
			if slices.Contains(node.Keys, q.Key) {
				// Found a query that affects one of this thread's leaves
				if q.Method == tree_api.MethodFind {
					findQueries = append(findQueries, q) // to be serviced here
					break
				} else if q.Method == tree_api.MethodInsert || q.Method == tree_api.MethodDelete {
					// Add other queries into map, to be serviced later
					val, ok := O_L_i[node]
					if !ok {
						val := make([]tree_api.Query, 0)
						O_L_i[node] = append(val, q)
					} else {
						O_L_i[node] = append(val, q)
					}
					break
				}
			}
		}
	}
	// Service appropriate find queries
	for _, q := range findQueries {
		val, err := t.Find(q.Key, false)
		if err == nil {
			res = append(res, val)
		}
	}
	return res, O_L_i
}

// From Stage 3 code
func bigSplit(node *Node) ([]int, []interface{}) {
	newKeys := make([]int, 0)
	newNodes := make([]interface{}, 0)
	newNodeCount := ((node.NumKeys + minOrder - 1) / minOrder) - 1
	assert(newNodeCount > 0)
	currIndex := minOrder
	for nodeNum := 1; nodeNum < newNodeCount; nodeNum++ {
		newNode, _ := makeNode()
		newNode.Parent = node.Parent
		newNode.NumKeys = minOrder - 1
		for j := 0; j < minOrder; j++ {
			newNode.Keys[j] = node.Keys[j+currIndex]
		}
		for j := 0; j < minOrder; j++ {
			newNode.Pointers[j] = node.Pointers[j+minOrder]
		}
		newKeys = append(newKeys, newNode.Keys[nodeNum*minOrder])
		newNodes = append(newNodes, newNode)
		currIndex += minOrder
	}
	node.NumKeys = minOrder - 1
	node.Keys = node.Keys[:node.NumKeys]
	node.Pointers = node.Pointers[:node.NumKeys+1]
	return newKeys, newNodes
}

func addModificationIntoList(node *Node, mod *Modification, M_i map[*Node]([]*Modification)) map[*Node]([]*Modification) {
	val, ok := M_i[node]
	if !ok {
		val = make([]*Modification, 0)
		M_i[node] = append(val, mod)
	} else {
		M_i[node] = append(val, mod)
	}
	return M_i
}

func (t *LockFreeTree) ModifyLeafNode(queriesToBeServiced map[*Node]([]tree_api.Query), L_i_prime []*Node) map[*Node]([]*Modification) {
	M_i := make(map[*Node]([]*Modification))
	for node, queries := range queriesToBeServiced {
		for _, q := range queries {
			if q.Method == tree_api.MethodInsert {
				insertIntoLeaf(node, q.Key, q.Pointer)
			} else if q.Method == tree_api.MethodDelete {
				removeEntryFromNode(node, q.Key, q.Pointer)
			}
		}
	}
	for node, _ := range queriesToBeServiced {
		if node.NumKeys > maxOrder {
			// Split Case
			newKeys, newNodes := bigSplit(node)
			mod := &Modification{Split, node.Parent, &SplitData{newKeys, newNodes}, nil, nil}
			M_i = addModificationIntoList(node, mod, M_i)
		} else if node.NumKeys < minOrder {
			// Underflow Case
			leafKeys := node.Keys
			keyToRemove := node.Keys[0]
			childKeys := make([]int, 0)
			childKeys = append(childKeys, keyToRemove)
			childPtrs := make([]interface{}, 0)
			childPtrs = append(childPtrs, node)
			mod := &Modification{Underflow, node.Parent, nil, &UnderflowData{childKeys, childPtrs}, leafKeys}
			M_i = addModificationIntoList(node, mod, M_i)
		}
	}
	return M_i
}

func (t *LockFreeTree) Stage2Logic(i int, num_threads int, sharedLeafData [][]*Node, queries []tree_api.Query, R [][]*tree_api.Record) /*[]*Node*/ map[*Node]([]*Modification) {
	// Redistribute Work
	L_i_prime := t.RedistributeWorkLeaves(i, sharedLeafData)
	res, O_L_i := t.ResolveHazards(L_i_prime, queries)
	// Update shared results slice
	R[i] = res
	M_i := t.ModifyLeafNode(O_L_i, L_i_prime)
	return M_i
	// Modify leaves independently
}

func (t *LockFreeTree) modifySharedModLists(index int, sharedLeafData [][]*Node, sharedModLists [](map[*Node]([]*Modification)), queries []tree_api.Query, R [][]*tree_api.Record, palmMaxThreadCount int, wg *sync.WaitGroup /*, testing [][]*Node*/) {
	defer wg.Done()
	// fmt.Printf("in modsharedlists index: %d\n", index)
	res := t.Stage2Logic(index, palmMaxThreadCount, sharedLeafData, queries, R)
	// testing[index] = res
	sharedModLists[index] = res
}

func (t *LockFreeTree) Stage2(sharedLeafData [][]*Node, palmMaxThreadCount int, queries []tree_api.Query) ([](map[*Node]([]*Modification)), [][]*tree_api.Record) {
	fmt.Printf("Starting Stage 2\n")
	var wg2 sync.WaitGroup
	dbg := false

	// testing := make([][]*Node, palmMaxThreadCount)
	// for i := 0; i < palmMaxThreadCount; i++ {
	// 	testing[i] = make([]*Node, 0)
	// }

	// Set up
	sharedModLists := make([]map[*Node]([]*Modification), palmMaxThreadCount)
	for i := 0; i < palmMaxThreadCount; i++ {
		sharedModLists[i] = make(map[*Node]([]*Modification))
	}
	R := make([][]*tree_api.Record, palmMaxThreadCount+1) // could potentially have off by 1 issues (just not using 0 i guess)
	for i := 0; i < palmMaxThreadCount; i++ {
		R[i] = make([]*tree_api.Record, 0)
	}

	// Do threads
	for i := 0; i < palmMaxThreadCount; i++ {
		wg2.Add(1) // Increment the counter for each goroutine
		go t.modifySharedModLists(i, sharedLeafData, sharedModLists, queries, R, palmMaxThreadCount, &wg2 /*, testing*/)
	}
	wg2.Wait()
	fmt.Println("All workers have completed Stage 2.")

	if dbg {
		// fmt.Println("Printing testing redistributed leaves vals")
		// for idx, L_i := range testing {
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

	return sharedModLists, R
}
