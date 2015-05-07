// The memory package implements an in-memory storage engine. This is primarily
// for testing.
package memory

import (
	"sync"

	"github.com/chop-dbhi/origins/engine"
)

type Tx struct {
	e *Engine
}

func (t *Tx) Get(p, k string) ([]byte, error) {
	if b, ok := t.e.parts[p]; ok {
		return b[k], nil
	}

	return nil, nil
}

func (t *Tx) Set(p, k string, v []byte) error {
	var (
		ok bool
		b  map[string][]byte
	)

	if b, ok = t.e.parts[p]; !ok {
		b = make(map[string][]byte)
		t.e.parts[p] = b
	}

	b[k] = v

	return nil
}

func (t *Tx) Incr(p, k string) (uint64, error) {
	var (
		ok bool
		id uint64
		v  []byte
		b  map[string][]byte
	)

	if b, ok = t.e.parts[p]; !ok {
		b = make(map[string][]byte)
		t.e.parts[p] = b
	}

	if v, ok = b[k]; ok {
		id = engine.DecodeCounter(v)
	}

	id++

	b[k] = engine.EncodeCounter(id)

	return id, nil
}

// Engine is an in-memory store that keeps data in keyed parts.
type Engine struct {
	parts map[string]map[string][]byte

	sync.Mutex
}

func (e *Engine) Get(p, k string) ([]byte, error) {
	e.Lock()
	defer e.Unlock()

	t := &Tx{e}

	return t.Get(p, k)
}

func (e *Engine) Set(p, k string, v []byte) error {
	e.Lock()
	defer e.Unlock()

	t := &Tx{e}

	return t.Set(p, k, v)
}

func (e *Engine) Incr(p, k string) (uint64, error) {
	e.Lock()
	defer e.Unlock()

	t := &Tx{e}

	return t.Incr(p, k)
}

func (e *Engine) Multi(f func(tx engine.Tx) error) error {
	e.Lock()
	defer e.Unlock()

	return f(&Tx{e})
}

// Open initializes a new Engine and returns it.
func Init(opts engine.Options) (*Engine, error) {
	e := Engine{
		parts: make(map[string]map[string][]byte),
	}

	return &e, nil
}