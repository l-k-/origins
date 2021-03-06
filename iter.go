package origins

// Iterator is an interface that reads from an underlying stream of facts
// and makes facts available through the Next() method. The usage of the
// this interface is as follows:
//
//		// The fact will be nil if the underlying stream is exhaused
//		// or an error occurred.
//		for f := it.Next(); f != nil {
//			// Do something with the fact
//		}
//
//		if err := it.Err(); err != nil {
//			// Handle the error.
//		}
//
type Iterator interface {
	// Next returns the next available fact in the stream.
	Next() *Fact

	// Err returns an error if one occurred while iterating facts.
	Err() error
}

// Writer is an interface that defines the Write method. It takes a fact
// and writes it to the underlying value.
type Writer interface {
	// Write writes a fact to the underlying stream.
	Write(*Fact) error
}

// Flusher is an interface that defines the Flush method. Types that
// are buffered and require being *flushed* at the end of processing
// implement this method.
type Flusher interface {
	Flush() error
}

// Read consumes up to len(buf) facts of the iterator and puts them into a buffer.
func Read(it Iterator, buf Facts) (int, error) {
	var (
		i int
		f *Fact
	)

	// Range over the buffer and fill with facts.
	for _ = range buf {
		if f = it.Next(); f == nil {
			break
		}

		buf[i] = f
		i++
	}

	return i, it.Err()
}

type sliceIter struct {
	iter          Iterator
	offset, limit int
	count, index  int
}

func (s *sliceIter) Next() *Fact {
	// If the limit has been reached, exit.
	if s.limit > 0 && s.count >= s.limit {
		return nil
	}

	var fact *Fact

	for {
		if fact = s.iter.Next(); fact == nil {
			return nil
		}

		// Skip fact if not within range.
		if s.index >= s.offset {
			s.index++
			s.count++
			return fact
		}

		s.index++
	}

	return nil
}

func (s *sliceIter) Err() error {
	return s.iter.Err()
}

// Slice applies an offset and limit to the iterator.
func Slice(iter Iterator, offset, limit int) Iterator {
	return &sliceIter{
		iter:   iter,
		offset: offset,
		limit:  limit,
	}
}

// ReadAll reads all facts from the reader.
func ReadAll(it Iterator) (Facts, error) {
	buf := NewBuffer(nil)

	if _, err := Copy(it, buf); err != nil {
		return nil, err
	}

	return buf.Facts(), nil
}

// Copy reads all facts from the reader and copies them to the writer.
// The number of facts written is returned and an error if present.
func Copy(it Iterator, w Writer) (int, error) {
	var (
		n   int
		f   *Fact
		err error
	)

	for {
		if f = it.Next(); f == nil {
			break
		}

		if err = w.Write(f); err != nil {
			break
		}

		n++
	}

	// If a write error did not occur, set the read error.
	if err != nil {
		err = it.Err()
	}

	// If the writer implements Flusher, flush it.
	switch x := w.(type) {
	case Flusher:
		err = x.Flush()
	}

	return n, err
}

// Map takes an iterator and passes each fact into the map function.
func Map(iter Iterator, proc func(*Fact) error) error {
	var (
		err  error
		fact *Fact
	)

	for {
		if fact = iter.Next(); fact == nil {
			break
		}

		if err = proc(fact); err != nil {
			return err
		}
	}

	return iter.Err()
}

func uniqueIdents(iter Iterator, identer func(*Fact) *Ident) (Idents, error) {
	var (
		ok     bool
		err    error
		key    [2]string
		fact   *Fact
		ident  *Ident
		idents Idents
	)

	seen := make(map[[2]string]struct{})

	for {
		if fact = iter.Next(); fact == nil {
			break
		}

		// Get the identity.
		ident = identer(fact)

		key[0] = ident.Domain
		key[1] = ident.Name

		if _, ok = seen[key]; !ok {
			seen[key] = struct{}{}
			idents = append(idents, ident)
		}
	}

	if err = iter.Err(); err != nil {
		return nil, err
	}

	return idents, nil
}

// Entities extract a unique set of entity identities from the iterator.
func Entities(iter Iterator) (Idents, error) {
	return uniqueIdents(iter, func(fact *Fact) *Ident {
		return fact.Entity
	})
}

// Attributes extract a unique set of attribute identities from the iterator.
func Attributes(iter Iterator) (Idents, error) {
	return uniqueIdents(iter, func(fact *Fact) *Ident {
		return fact.Attribute
	})
}

// Values extract a unique set of values identities from the iterator.
func Values(iter Iterator) (Idents, error) {
	return uniqueIdents(iter, func(fact *Fact) *Ident {
		return fact.Value
	})
}

// Transactions extract a unique set of transaction IDs from the iterator.
func Transactions(iter Iterator) ([]uint64, error) {
	var (
		ok   bool
		err  error
		fact *Fact
		tx   uint64
		txes []uint64
	)

	seen := make(map[uint64]struct{})

	for {
		if fact = iter.Next(); fact == nil {
			break
		}

		tx = fact.Transaction

		if _, ok = seen[tx]; !ok {
			seen[tx] = struct{}{}
			txes = append(txes, tx)
		}
	}

	if err = iter.Err(); err != nil {
		return nil, err
	}

	return txes, nil
}

type filterer struct {
	iter      Iterator
	predicate func(*Fact) bool
}

func (f *filterer) Next() *Fact {
	var fact *Fact

	for {
		if fact = f.iter.Next(); fact == nil {
			break
		}

		// Fact matches, return.
		if f.predicate(fact) {
			return fact
		}
	}

	return nil
}

func (f *filterer) Err() error {
	return f.iter.Err()
}

// Filter filters facts consumed from the iterator and returns an iterator.
func Filter(iter Iterator, predicate func(*Fact) bool) Iterator {
	return &filterer{
		iter:      iter,
		predicate: predicate,
	}
}

// Entity takes an iterator and returns the all facts about an entity.
func Entity(iter Iterator, id *Ident) Iterator {
	return &filterer{
		iter: iter,
		predicate: func(f *Fact) bool {
			return f.Entity.Is(id)
		},
	}
}

// First takes an iterator and predicate function and returns the first
// fact that matches the predicate.
func First(iter Iterator, predicate func(*Fact) bool) *Fact {
	f := Filter(iter, predicate)

	return f.Next()
}

// Exists takes an iterator and predicate function and returns true if
// the predicate matches.
func Exists(iter Iterator, predicate func(*Fact) bool) bool {
	return First(iter, predicate) != nil
}

type grouper struct {
	iter   Iterator
	last   *Fact
	buffer *Buffer
	cmp    func(f1, f2 *Fact) bool
}

type FactsIterator interface {
	Next() Facts
	Err() error
}

func (g *grouper) Next() Facts {
	if g.iter.Err() != nil {
		return nil
	}

	var (
		f  *Fact
		ok = true
	)

	// Boundary fact for the last iteration. Add it to the current.
	// buffer. The last fact is also used for comparison with
	// subsequent iterations.
	if g.last != nil {
		g.buffer.Write(g.last)
	}

	// Iterate over the facts and write them to the buffer until
	// the comparator returns false.
	for {
		if f = g.iter.Next(); f == nil {
			// Unset the reference so subsequent calls don't emit
			// single fact groups.
			g.last = nil
			break
		}

		// Compare against the last fact if one is defined. If they don't
		// match, break and return the buffered facts.
		if g.last == nil {
			g.buffer.Write(f)
		} else {
			if ok = g.cmp(g.last, f); ok {
				g.buffer.Write(f)
			}
		}

		// Maintain a reference to the last fact for the next iteration.
		g.last = f

		if !ok {
			break
		}
	}

	// Check for an error.
	if err := g.iter.Err(); err != nil {
		return nil
	}

	if g.buffer.Len() > 0 {
		return g.buffer.Facts()
	}

	return nil
}

func (g *grouper) Err() error {
	return g.iter.Err()
}

func Groupby(iter Iterator, cmp func(f1, f2 *Fact) bool) FactsIterator {
	return &grouper{
		iter:   iter,
		cmp:    cmp,
		buffer: NewBuffer(nil),
	}
}

func MapFacts(iter FactsIterator, proc func(facts Facts) error) error {
	var (
		err   error
		facts Facts
	)

	for {
		if facts = iter.Next(); facts == nil {
			break
		}

		if err = proc(facts); err != nil {
			return err
		}
	}

	return iter.Err()
}
