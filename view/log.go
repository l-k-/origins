package view

import (
	"fmt"
	"time"

	"github.com/chop-dbhi/origins"
	"github.com/chop-dbhi/origins/storage"
)

type segment struct {
	// ID of the segment. This is also the ID of the transaction that
	// resulted in this segment.
	ID uint64

	// The domain this segment corresponds to.
	Domain string

	// Time the segment was committed.
	Time time.Time

	// Number of blocks in this segment.
	Blocks int

	// Total number of facts in the segment.
	Count int

	// Total number of bytes of the segment take up.
	Bytes int

	// ID of the segment that acted as the basis for this one.
	Base uint64

	// ID of the segment that follows this one. Typically this is
	// the same value as Base, except when a conflict is resolved and
	// the segment position is changed.
	Next uint64
}

// loadSegment loads a segment header from storage.
func loadSegment(e storage.Engine, d string, s uint64) (*segment, error) {
	var (
		b   []byte
		err error
		key string
	)

	seg := segment{
		ID:     s,
		Domain: d,
	}

	// Get segment header.
	key = fmt.Sprintf("segment.%d", s)

	if b, err = e.Get(d, key); err != nil {
		return nil, err
	}

	// Does not exist.
	if b == nil {
		return nil, nil
	}

	if err = unmarshalSegment(b, &seg); err != nil {
		return nil, err
	}

	return &seg, nil
}

// loadBlock loads a block of facts from storage.
func loadBlock(e storage.Engine, d string, s uint64, i int) (origins.Facts, error) {
	var (
		b   []byte
		err error
		key string
	)

	// Get block.
	key = fmt.Sprintf("segment.%d.%d", s, i)

	if b, err = e.Get(d, key); err != nil {
		return nil, err
	}

	// Does not exist.
	if b == nil {
		return nil, nil
	}

	return unmarshalFacts(b, d, s)
}

// logIter maintains state of a log that is being read.
type logIter struct {
	name    string
	domain  string
	head    uint64
	engine  storage.Engine
	segment *segment
	block   origins.Facts
	bindex  int
	bpos    int
	done    bool
}

// nextSegment
func (li *logIter) nextSegment() error {
	var id uint64

	if li.segment == nil {
		id = li.head
	} else {
		id = li.segment.Next
	}

	segment, err := loadSegment(li.engine, li.domain, id)

	if err != nil {
		return err
	}

	// Update state.
	li.segment = segment
	li.bindex = 0
	li.block = nil
	li.bpos = 0

	if segment == nil {
		li.done = true
		return nil
	}

	return nil
}

// nextBlock returns the block that has the next fact or nil or no
// more blocks exist.
func (li *logIter) nextBlock() error {
	if li.done {
		return nil
	}

	// Existing block and still has facts.
	if li.block != nil && li.bpos < len(li.block) {
		return nil
	}

	// First segment or no blocks left in segment.
	if li.segment == nil || li.bindex == li.segment.Blocks {
		if err := li.nextSegment(); err != nil {
			return err
		}

		if li.done {
			return nil
		}
	}

	// Error loading block
	block, err := loadBlock(li.engine, li.segment.Domain, li.segment.ID, li.bindex)

	if err != nil {
		return err
	}

	// Block does not exist.
	if block == nil {
		li.done = true
		return nil
	}

	li.block = block
	li.bpos = 0
	li.bindex++

	return nil
}

func (li *logIter) Next() (*origins.Fact, error) {
	if err := li.nextBlock(); err != nil {
		return nil, err
	}

	if li.done {
		return nil, nil
	}

	fact := li.block[li.bpos]
	li.bpos++

	return fact, nil
}

// A Log is an ordered sequence of facts within a domain.
type Log struct {
	Name   string
	Domain string

	head   uint64
	engine storage.Engine
}

// Iter returns a value that implements the Iterator interface. This can be
// called multiple times for independent consumers.
func (l *Log) Iter() origins.Iterator {
	return &logIter{
		domain: l.Domain,
		engine: l.engine,
		head:   l.head,
	}
}

// OpenLog opens a log for reading.
func OpenLog(e storage.Engine, d, n string) (*Log, error) {
	var (
		b   []byte
		err error
	)

	if b, err = e.Get(d, n); err != nil {
		return nil, err
	}

	// Does not exist.
	if b == nil {
		return nil, nil
	}

	l := Log{
		Name:   n,
		Domain: d,
		engine: e,
	}

	if err = unmarshalLog(b, &l); err != nil {
		return nil, err
	}

	return &l, nil
}