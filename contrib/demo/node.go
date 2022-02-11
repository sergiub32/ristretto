package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/dustin/go-humanize"
	"github.com/sergiub32/ristretto/z"
)

type node struct {
	next *node
	val  int
}

var (
	nodeSz = int(unsafe.Sizeof(node{})) //nolint:unused,varcheck,gochecknoglobals,lll,deadcode,revive // adopt fork, do not touch it
	alloc  *z.Allocator                 //nolint:unused,varcheck,gochecknoglobals,lll,deadcode,revive // adopt fork, do not touch it
)

func printNode(n *node) {
	if n == nil {
		return
	}
	if n.val%100000 == 0 {
		fmt.Printf("node: %d\n", n.val)
	}
	printNode(n.next)
}

func main() {
	N := 2000001
	root := newNode(-1)
	n := root
	for i := 0; i < N; i++ {
		nn := newNode(i)
		n.next = nn
		n = nn
	}
	fmt.Printf("Allocated memory: %s Objects: %d\n",
		humanize.IBytes(uint64(z.NumAllocBytes())), N)

	runtime.GC()
	printNode(root)
	fmt.Println("printing done")

	if alloc != nil {
		alloc.Release()
	} else {
		n = root
		for n != nil {
			left := n
			n = n.next
			freeNode(left)
		}
	}
	fmt.Printf("After freeing. Allocated memory: %d\n", z.NumAllocBytes())

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("HeapAlloc: %s\n", humanize.IBytes(ms.HeapAlloc))
}
