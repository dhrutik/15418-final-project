package lock_free

import "sync"

func assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
}

func (t *LockFreeTree) getUpdatedModList(sharedModLists [](map[*Node]([]Modification)), threadId int) map[*Node]([]Modification) {
	previousThreadNodes := make(map[*Node]bool)
	updatedModList := make(map[*Node]([]Modification))
	for i := 0; i < threadId; i++ {
		for node, _ := range sharedModLists[i] {
			previousThreadNodes[node] = true
		}
	}
	for node, modList := range sharedModLists[threadId] {
		if _, ok := previousThreadNodes[node]; !ok {
			updatedModList[node] = modList
		}
	}
	return updatedModList
}

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
		for j := 0; j < minOrder ; j++ {
			newNode.Keys[j] = node.Keys[j+currIndex]
		}
		for j := 0; j < minOrder; j++ {
			newNode.Pointers[j] = node.Pointers[j+minOrder]
		}
		newKeys = append(newKeys, newNode.Keys[nodeNum * minOrder])
		newNodes = append(newNodes, newNode)
		currIndex += minOrder
	}
	node.NumKeys = minOrder - 1
	node.Keys = node.Keys[:node.NumKeys]
	node.Pointers = node.Pointers[:node.NumKeys+1]
	return newKeys, newNodes
}

func getLeafKeys(node *Node) []int {
	if node.IsLeaf {
		return node.Keys
	}
	leafKeys := make([]int, 0)
	for i := 0; i < node.NumKeys+1; i++ {
		leafKeys = append(leafKeys, getLeafKeys(node.Pointers[i].(*Node))...)
	}
	return leafKeys
}

func (t *LockFreeTree) modifyInternalNode(node *Node, mod Modification) *Modification {
	if mod.ModType == Split {
		for i, updateKey := range mod.SplitData.NewKeys {
			left_index := getLeftIndex(node, mod.SplitData.NewNodes[i].(*Node))
			insertIntoNode(node, left_index, updateKey, mod.SplitData.NewNodes[i].(*Node))
		}
	} else if mod.ModType == Underflow {
		for i, updateKey := range mod.UnderflowData.ChildKeys {
			removeEntryFromNode(node, updateKey, mod.UnderflowData.ChildPtrs[i])
		}
	}
	if node.NumKeys > maxOrder {
		newKeys, newNodes := bigSplit(node)
		return &Modification{Split, node.Parent, &SplitData{newKeys, newNodes}, nil, mod.OrphanedKeys}
	} else if node.NumKeys < minOrder {
		leafKeys := getLeafKeys(node)
		keyToRemove := node.Keys[0]
		childKeys := make([]int, 0)
		childKeys = append(childKeys, keyToRemove)
		childPtrs := make([]interface{}, 0)
		childPtrs = append(childPtrs, node)
		return &Modification{Underflow, node.Parent, nil, &UnderflowData{childKeys, childPtrs}, append(leafKeys, mod.OrphanedKeys...)}
	}

	return &Modification{NoMod, node.Parent, nil, nil, mod.OrphanedKeys}
}

func (t *LockFreeTree) stage3Thread(sharedModLists [](map[*Node]([]Modification)), newSharedModLists [](map[*Node]([]Modification)), threadId int, depth int, wg *sync.WaitGroup) {
	for d := 1; d < depth; d++ {
		updatedModList := t.getUpdatedModList(sharedModLists, threadId)
		for node, modList := range updatedModList {
			for _, mod := range modList {
				newMod := t.modifyInternalNode(node, mod)
				if newMod != nil {
					newSharedModLists[threadId][newMod.Parent] = append(newSharedModLists[threadId][newMod.Parent], *newMod)
				}
			}
		}
		wg.Done()
	}
}

func (t *LockFreeTree) Stage3(sharedModLists [](map[*Node]([]Modification)), palmMaxThreadCount int) [](map[*Node]([]Modification)) {
	depth := t.height()
	// spin off threads first, passing them sharedModLists and thread id
	wg := sync.WaitGroup{}
	newSharedModLists := make([](map[*Node]([]Modification)), palmMaxThreadCount)
	for i := 0; i < palmMaxThreadCount; i++ {
		go t.stage3Thread(sharedModLists, newSharedModLists, i, depth, &wg)
	}
	// this loop just syncs after each loop iteration
	// allocate a new sharedModList
	for d := 1; d < depth; d++ {
		wg.Add(palmMaxThreadCount)
		wg.Wait()
		// copy over newSharedModLists to sharedModLists
		for j := 0; j < len(sharedModLists); j++ {
			sharedModLists[j] = newSharedModLists[j]
		}
		wg.Add(palmMaxThreadCount)
		wg.Wait()
	}
	return nil
}
