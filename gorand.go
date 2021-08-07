package gorand

import (
	"log"
	"math/rand"
	"sync"
)

var (
	_ rand.Source   = &Source{}
	_ rand.Source64 = &Source{}
	_ rand.Source64 = rand.NewSource(1).(rand.Source64)
)

type Source struct {
	pool sync.Pool
}

func NewSource(seed int64) *Source {
	seedGenerator := newRand(seed)

	return &Source{
		pool: sync.Pool{
			New: func() interface{} {
				return rand.NewSource(seedGenerator.Int63())
			},
		},
	}
}

func (r *Source) Int63() int64 {
	var (
		generator = r.pool.Get().(rand.Source64)
		number    = generator.Int63()
	)

	r.pool.Put(generator)

	return number
}

func (r *Source) Uint64() uint64 {
	var (
		generator = r.pool.Get().(rand.Source64)
		number    = generator.Uint64()
	)

	r.pool.Put(generator)

	return number
}

func (r *Source) Seed(seed int64) {
	log.Println("[gorand] seed not support")
}

// lockedSource allows a random number generator to be used by multiple goroutines concurrently.
// The code is very similar to math/rand.lockedSource, which is unfortunately not exposed.
type lockedSource struct {
	src rand.Source
	mut sync.Mutex
}

// NewRand returns a rand.Rand that is threadsafe.
func newRand(seed int64) *rand.Rand {
	return rand.New(&lockedSource{src: rand.NewSource(seed)}) // nolint:gosec // it's ok
}

func (r *lockedSource) Int63() (n int64) {
	r.mut.Lock()
	n = r.src.Int63()
	r.mut.Unlock()

	return
}

// Seed implements Seed() of Source.
func (r *lockedSource) Seed(seed int64) {
	r.mut.Lock()
	r.src.Seed(seed)
	r.mut.Unlock()
}
