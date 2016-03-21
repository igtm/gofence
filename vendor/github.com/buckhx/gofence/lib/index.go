package geofence

import (
	"fmt"
	"sync"

	"github.com/buckhx/diglet/geo"
)

//FenceIndex is a dictionary of multiple fences. Useful if you have multiple data sets that need to be searched
type FenceIndex interface {
	// Set the GeoFence
	Set(name string, fence GeoFence)
	// Get the GeoFence at the key, return nil if doesn't exist
	Get(name string) GeoFence
	// Add a feature to the GeoFence at the key
	Add(name string, feature *geo.Feature) error
	// Search for the coordinate at the key
	Search(name string, c geo.Coordinate) ([]*geo.Feature, error)
	// List the keys of the indexed fences
	Keys() []string
}

// Returns a thread-safe FenceIndex
func NewFenceIndex() FenceIndex {
	return NewMutexFenceIndex()
}

type UnsafeFenceIndex struct {
	fences map[string]GeoFence
}

func NewUnsafeFenceIndex() *UnsafeFenceIndex {
	return &UnsafeFenceIndex{fences: make(map[string]GeoFence)}
}

func (idx *UnsafeFenceIndex) Set(name string, fence GeoFence) {
	idx.fences[name] = fence
}

func (idx *UnsafeFenceIndex) Get(name string) (fence GeoFence) {
	return idx.fences[name]
}

func (idx *UnsafeFenceIndex) Add(name string, feature *geo.Feature) (err error) {
	fence, ok := idx.fences[name]
	if !ok {
		return fmt.Errorf("FenceIndex does not contain fence %q", name)
	}
	fence.Add(feature)
	return
}

func (idx *UnsafeFenceIndex) Search(name string, c geo.Coordinate) (matchs []*geo.Feature, err error) {
	fence, ok := idx.fences[name]
	if !ok {
		err = fmt.Errorf("FenceIndex does not contain fence %q", name)
		return
	}
	matchs = fence.Get(c)
	return
}

func (idx *UnsafeFenceIndex) Keys() (keys []string) {
	for k := range idx.fences {
		keys = append(keys, k)
	}
	return
}

type MutexFenceIndex struct {
	fences *UnsafeFenceIndex
	sync.RWMutex
}

func NewMutexFenceIndex() *MutexFenceIndex {
	return &MutexFenceIndex{fences: NewUnsafeFenceIndex()}
}

func (idx *MutexFenceIndex) Set(name string, fence GeoFence) {
	idx.Lock()
	defer idx.Unlock()
	idx.fences.Set(name, fence)
}

func (idx *MutexFenceIndex) Get(name string) GeoFence {
	idx.RLock()
	defer idx.RUnlock()
	return idx.fences.Get(name)
}

func (idx *MutexFenceIndex) Add(name string, feature *geo.Feature) error {
	idx.Lock()
	defer idx.Unlock()
	return idx.fences.Add(name, feature)
}

func (idx *MutexFenceIndex) Search(name string, c geo.Coordinate) ([]*geo.Feature, error) {
	idx.RLock()
	defer idx.RUnlock()
	return idx.fences.Search(name, c)
}

func (idx *MutexFenceIndex) Keys() []string {
	idx.RLock()
	defer idx.RUnlock()
	return idx.fences.Keys()
}
