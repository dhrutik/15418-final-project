package global_lock_tree

import (
	"main/seq_tree"
	"main/tree_api"
	"sync"
)

type GlobalLockTree struct {
	tree *seq_tree.Tree
	lock sync.Mutex
}

func NewTree() tree_api.BPTree {
	return &GlobalLockTree{tree: seq_tree.NewTree(), lock: sync.Mutex{}}
}

func (t *GlobalLockTree) Insert(key int, value []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.tree.Insert(key, value)
}

func (t *GlobalLockTree) Delete(key int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.tree.Delete(key)
}

func (t *GlobalLockTree) Find(key int, verbose bool) (*tree_api.Record, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.tree.Find(key, verbose)
}

func (t *GlobalLockTree) PrintTree() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.tree.PrintTree()
}

func (t *GlobalLockTree) FindAndPrint(key int, verbose bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.tree.FindAndPrint(key, verbose)
}

func (t *GlobalLockTree) FindAndPrintRange(key_start, key_end int, verbose bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.tree.FindAndPrintRange(key_start, key_end, verbose)
}

func (t *GlobalLockTree) PrintLeaves() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.tree.PrintLeaves()
}

func (t *GlobalLockTree) PalmBasic(key_count int, num_threads int) {}

func (t *GlobalLockTree) Palm(query []tree_api.Query, num_threads int) [][]*tree_api.Record {
	return nil
}

// func (t *GlobalLockTree) Stage1(Q []tree_api.Query, i int, num_threads int) {}
// func (t *GlobalLockTree) Stage2(Q []tree_api.Query, i int, num_threads int) {}
// func (t *GlobalLockTree) Stage3(Q []tree_api.Query, i int, num_threads int) {}
// func (t *GlobalLockTree) Stage4(Q []tree_api.Query, i int, num_threads int) {}
