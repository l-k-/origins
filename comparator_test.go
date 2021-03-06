package origins

import (
	"testing"

	"github.com/chop-dbhi/origins/chrono"
	"github.com/stretchr/testify/assert"
)

var (
	compFacts                            Facts
	bobName, sueName, bobColor, sueColor *Fact
)

func init() {
	e1 := &Ident{"people", "bob"}
	e2 := &Ident{"people", "sue"}

	a1 := &Ident{"people", "name"}
	a2 := &Ident{"people", "color"}

	v1 := &Ident{"", "Bob"}
	v2 := &Ident{"", "Sue"}
	v3 := &Ident{"", "red"}
	v4 := &Ident{"", "blue"}

	t0 := chrono.MustParse("2015-01-01")
	t1 := chrono.MustParse("2015-01-02")

	bobName = &Fact{
		Time:        t0,
		Entity:      e1,
		Attribute:   a1,
		Value:       v1,
		Transaction: 1,
	}

	bobColor = &Fact{
		Time:        t0,
		Entity:      e1,
		Attribute:   a2,
		Value:       v3,
		Transaction: 1,
	}

	sueName = &Fact{
		Time:        t1,
		Entity:      e2,
		Attribute:   a1,
		Value:       v2,
		Transaction: 2,
	}

	sueColor = &Fact{
		Time:        t1,
		Entity:      e2,
		Attribute:   a2,
		Value:       v4,
		Transaction: 2,
	}

	compFacts = Facts{
		bobName,
		bobColor,
		sueName,
		sueColor,
	}
}

func testComparator(t *testing.T, comp Comparator) Facts {
	facts := make(Facts, 4)
	copy(facts, compFacts)

	Timsort(facts, comp)

	return facts
}

func TestEntityComparator(t *testing.T) {
	facts := testComparator(t, EntityComparator)

	// bob, bob, sue, sue
	exp := Facts{
		bobName,
		bobColor,
		sueName,
		sueColor,
	}

	assert.Equal(t, exp, facts)
}

func TestAttributeComparator(t *testing.T) {
	facts := testComparator(t, AttributeComparator)

	// color, color, name, name
	exp := Facts{
		bobColor,
		sueColor,
		bobName,
		sueName,
	}

	assert.Equal(t, exp, facts)
}

func TestValueComparator(t *testing.T) {
	facts := testComparator(t, ValueComparator)

	// blue, Bob, Sue, red
	exp := Facts{
		bobName,
		sueName,
		sueColor,
		bobColor,
	}

	assert.Equal(t, exp, facts)
}

func TestTimeComparator(t *testing.T) {
	facts := testComparator(t, TimeComparator)

	// t0, t0, t1, t1
	exp := Facts{
		bobName,
		bobColor,
		sueName,
		sueColor,
	}

	assert.Equal(t, exp, facts)
}

func TestEAVTComparator(t *testing.T) {
	facts := testComparator(t, EAVTComparator)

	// bob:color:red, bob:name:Bob, sue:color:blue, sue:name:Sue
	exp := Facts{
		bobColor,
		bobName,
		sueColor,
		sueName,
	}

	assert.Equal(t, exp, facts)
}

func TestAVETComparator(t *testing.T) {
	facts := testComparator(t, AVETComparator)

	// color:blue:sue, color:red:bob, name:Bob:bob, name:Sue:sue
	exp := Facts{
		sueColor,
		bobColor,
		bobName,
		sueName,
	}

	assert.Equal(t, exp, facts)
}

func TestAEVTComparator(t *testing.T) {
	facts := testComparator(t, AEVTComparator)

	// color:bob:red, color:sue:blue, name:bob:Bob, name:sue:Sue
	exp := Facts{
		bobColor,
		sueColor,
		bobName,
		sueName,
	}

	assert.Equal(t, exp, facts)
}

func TestVAETComparator(t *testing.T) {
	facts := testComparator(t, VAETComparator)

	// Bob:name:bob, Sue:name:sue, blue:color:sue, red:color:bob,
	exp := Facts{
		bobName,
		sueName,
		sueColor,
		bobColor,
	}

	assert.Equal(t, exp, facts)
}
