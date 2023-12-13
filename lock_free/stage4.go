package lock_free

// import (
// 	"errors"
// 	"fmt"
// 	"main/tree_api"
// 	"reflect"
// 	"sync"
// )

func (t *LockFreeTree) Stage4(finalModList [](map[*Node]([]*Modification)), palmMaxThreadCount int) {
	orphanedKeys := make([]int, 0)
	for _, modMap := range finalModList {
		for node, modList := range modMap {
			for _, mod := range modList {
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
				orphanedKeys = append(orphanedKeys, mod.OrphanedKeys...)
			}
		}
	}
	if t.Root.NumKeys > maxOrder {
		newKeys, newNodes := bigSplit(t.Root)
		newRoot, _ := makeNode()
		newRoot.Keys[0] = t.Root.Keys[0]
		newRoot.NumKeys++
		newRoot.Pointers[0] = t.Root
		for i, updateKey := range newKeys {
			left_index := getLeftIndex(newRoot, newNodes[i].(*Node))
			insertIntoNode(newRoot, left_index, updateKey, newNodes[i].(*Node))
		}
		t.Root = newRoot
	} else if t.Root.NumKeys == 0 {
		t.Root = nil
	}
	// defer t.Palm(len(orphanedKeys), palmMaxThreadCount)
}
