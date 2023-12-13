package tree_api

// import (
// 	"sync"
// )

type Record struct {
	Value []byte
}

// type Node struct {
// 	Pointers []interface{}
// 	Keys     []int
// 	Parent   *Node
// 	IsLeaf   bool
// 	NumKeys  int
// 	Next     *Node
// 	lock     sync.Mutex
// }

type Method int

const (
	MethodFind Method = iota
	MethodInsert
	MethodDelete
)

type Query struct {
	Method Method
	Key    int
	Done   bool
}

type BPTree interface {
	// type Node struct {
	// 	Pointers []interface{}
	// 	Keys     []int
	// 	Parent   *Node
	// 	IsLeaf   bool
	// 	NumKeys  int
	// 	Next     *Node
	// 	lock     sync.Mutex
	// }

	Insert(key int, value []byte) error
	Delete(key int) error
	Find(key int, verbose bool) (*Record, error)
	// PrintTree()
	// FindAndPrint(key int, verbose bool)
	// FindAndPrintRange(key_start, key_end int, verbose bool)
	// Stage1(Q []Query, i int, num_threads int)
	Palm(key_count int, num_threads int)
	// Stage1(Q []Query, i int, num_threads int)
	// Stage2(Q []Query, i int, num_threads int)
	// Stage3(Q []Query, i int, num_threads int)
	// Stage4(Q []Query, i int, num_threads int)
}
