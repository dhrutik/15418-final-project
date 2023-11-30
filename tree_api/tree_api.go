package tree_api

type Record struct {
	Value []byte
}

type BPTree interface {
	Insert(key int, value []byte) error
	Delete(key int) error
	Find(key int, verbose bool) (*Record, error)
	// PrintTree()
	// FindAndPrint(key int, verbose bool)
	// FindAndPrintRange(key_start, key_end int, verbose bool)
}
