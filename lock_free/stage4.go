package lock_free

import (
	"main/tree_api"
)

func (t *LockFreeTree) MakeOrphanedKeyInsertQueries(orphanedKeys []int) []tree_api.Query {
	queries := make([]tree_api.Query, 0)
	for _, key := range orphanedKeys {
		value := tree_api.Record{Value: []byte("value")}
		queries = append(queries, tree_api.Query{Method: tree_api.MethodInsert, Key: key, Done: false, Pointer: &value})
	}
	return queries
}

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
	queries := t.MakeOrphanedKeyInsertQueries(orphanedKeys)
	if len(queries) == 0 {
		return
	}
	defer t.Palm(queries, palmMaxThreadCount)
}
