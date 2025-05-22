package evil

import (
	"runtime"
	"sync"
	"unsafe"
)

var reg *registry

type registry struct {
	// values is a map of unsafe pointers keyed by name
	pinner runtime.Pinner
	mu     sync.RWMutex
	refs   map[string]unsafe.Pointer
}

func addToRegistry[T any](name string, val *T) {
	if reg == nil {
		reg = &registry{
			refs: make(map[string]unsafe.Pointer),
		}
	}

	reg.mu.Lock()
	defer reg.mu.Unlock()

	if reg.refs == nil {
		reg.refs = make(map[string]unsafe.Pointer)
	}

	ptr := unsafe.Pointer(val)

	// pin the pointer, we never want registered refs to be garbage collected
	reg.pinner.Pin(ptr)
	reg.refs[name] = ptr
}

func getFromRegistry[T any](name string) *T {
	if reg == nil {
		panic("registry is nil")
	}
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	if reg.refs == nil {
		panic("refs map is nil")
	}

	if ptr, ok := reg.refs[name]; ok {
		return (*T)(ptr)
	}

	panic("could not find value in registry")
}
