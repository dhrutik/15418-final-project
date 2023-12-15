package tree_api

// All B+ trees in this repo implement this interface.

type Record struct {
	Value []byte
}

type Method int

const (
	MethodFind Method = iota
	MethodInsert
	MethodDelete
)

type Query struct {
	Method  Method
	Key     int
	Done    bool
	Pointer *Record
}

type BPTree interface {
	Insert(key int, value []byte) error
	Delete(key int) error
	Find(key int, verbose bool) (*Record, error)
	PalmBasic(key_count int, num_threads int)
	Palm(queries []Query, num_threads int) [][]*Record
}
