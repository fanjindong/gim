package resource

import (
	"sync"
)

var BufPool *sync.Pool

func init() {
	BufPool = &sync.Pool{New: func() interface{} { return []byte{} }}
}
