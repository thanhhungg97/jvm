package runtime

import (
	"sync"
	"sync/atomic"
)

// Heap represents the JVM heap for object allocation
type Heap struct {
	// Object storage
	objects     map[uint64]any
	objectMutex sync.RWMutex

	// Object ID counter
	nextID atomic.Uint64

	// GC statistics
	allocCount atomic.Uint64
	freeCount  atomic.Uint64
	totalBytes atomic.Int64
	gcRuns     atomic.Uint64

	// GC configuration
	gcThreshold int64 // Trigger GC when heap exceeds this size
	gcEnabled   bool
}

// NewHeap creates a new heap
func NewHeap() *Heap {
	return &Heap{
		objects:     make(map[uint64]any),
		gcThreshold: 10 * 1024 * 1024, // 10MB default threshold
		gcEnabled:   true,
	}
}

// Alloc allocates a new object on the heap and returns its ID
func (h *Heap) Alloc(obj any) uint64 {
	id := h.nextID.Add(1)
	h.allocCount.Add(1)

	h.objectMutex.Lock()
	h.objects[id] = obj
	h.objectMutex.Unlock()

	// Estimate size and track
	size := int64(estimateSize(obj))
	h.totalBytes.Add(size)

	// Check if GC should run
	if h.gcEnabled && h.totalBytes.Load() > h.gcThreshold {
		h.GC(nil) // nil means we don't have root info yet
	}

	return id
}

// Get retrieves an object by ID
func (h *Heap) Get(id uint64) any {
	h.objectMutex.RLock()
	defer h.objectMutex.RUnlock()
	return h.objects[id]
}

// Free removes an object from the heap
func (h *Heap) Free(id uint64) {
	h.objectMutex.Lock()
	defer h.objectMutex.Unlock()

	if obj, exists := h.objects[id]; exists {
		size := int64(estimateSize(obj))
		h.totalBytes.Add(-size)
		delete(h.objects, id)
		h.freeCount.Add(1)
	}
}

// GC performs garbage collection using mark-sweep
func (h *Heap) GC(roots []any) {
	if !h.gcEnabled {
		return
	}

	h.gcRuns.Add(1)

	// Mark phase - mark all reachable objects
	marked := make(map[uint64]bool)

	// Mark from roots
	for _, root := range roots {
		h.mark(root, marked)
	}

	// Sweep phase - remove unmarked objects
	h.objectMutex.Lock()
	defer h.objectMutex.Unlock()

	toDelete := make([]uint64, 0)
	for id := range h.objects {
		if !marked[id] {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		if obj, exists := h.objects[id]; exists {
			size := int64(estimateSize(obj))
			h.totalBytes.Add(-size)
			delete(h.objects, id)
			h.freeCount.Add(1)
		}
	}
}

// mark recursively marks an object and its references
func (h *Heap) mark(obj any, marked map[uint64]bool) {
	if obj == nil {
		return
	}

	// Check if it's a heap-allocated object
	switch v := obj.(type) {
	case *Object:
		// Mark reference fields (these can contain other objects/arrays/strings)
		for _, ref := range v.Fields {
			h.mark(ref, marked)
		}

	case *Array:
		// Only reference arrays can hold other objects
		if v.IsRefArray() {
			for i := int32(0); i < v.Length; i++ {
				h.mark(v.GetRef(i), marked)
			}
		}

	case string:
		// Strings are immutable and managed by Go's GC
		// No action needed - they don't reference other Java objects

	case *Frame:
		// Frames contain local variables and operand stack with references
		// Mark all references in local vars and operand stack
		if v != nil && v.LocalVars != nil {
			for _, ref := range v.LocalVars.refs {
				h.mark(ref, marked)
			}
		}
		// Note: OperandStack refs would also need marking in a full implementation

	default:
		// Other types (int, long, placeholders like "System.out") don't need marking
		// They either don't hold references or are Go-managed
	}
}

// Stats returns GC statistics
func (h *Heap) Stats() HeapStats {
	return HeapStats{
		AllocCount:  h.allocCount.Load(),
		FreeCount:   h.freeCount.Load(),
		LiveObjects: h.allocCount.Load() - h.freeCount.Load(),
		TotalBytes:  h.totalBytes.Load(),
		GCRuns:      h.gcRuns.Load(),
		GCThreshold: h.gcThreshold,
	}
}

// HeapStats contains heap statistics
type HeapStats struct {
	AllocCount  uint64
	FreeCount   uint64
	LiveObjects uint64
	TotalBytes  int64
	GCRuns      uint64
	GCThreshold int64
}

// SetGCEnabled enables or disables garbage collection
func (h *Heap) SetGCEnabled(enabled bool) {
	h.gcEnabled = enabled
}

// SetGCThreshold sets the heap size threshold for triggering GC
func (h *Heap) SetGCThreshold(bytes int64) {
	h.gcThreshold = bytes
}

// estimateSize estimates the size of an object in bytes
func estimateSize(obj any) int {
	switch v := obj.(type) {
	case *Object:
		// Base size + fields
		size := 64 // Base object overhead
		size += len(v.Fields) * 16
		size += len(v.FieldSlots) * 8
		return size
	case *Array:
		size := 32 // Base array overhead
		if v.Ints != nil {
			size += len(v.Ints) * 4
		}
		if v.Longs != nil {
			size += len(v.Longs) * 8
		}
		if v.Floats != nil {
			size += len(v.Floats) * 4
		}
		if v.Doubles != nil {
			size += len(v.Doubles) * 8
		}
		if v.References != nil {
			size += len(v.References) * 8
		}
		return size
	case string:
		return 24 + len(v)
	default:
		return 16 // Default small object size
	}
}

// TriggerGC manually triggers garbage collection
func (h *Heap) TriggerGC(roots []any) {
	h.GC(roots)
}

// ObjectCount returns the number of live objects
func (h *Heap) ObjectCount() int {
	h.objectMutex.RLock()
	defer h.objectMutex.RUnlock()
	return len(h.objects)
}
