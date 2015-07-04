package origins

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/chop-dbhi/origins/chrono"
	"github.com/sirupsen/logrus"
)

// untitle takes a title-cased field and lowercases the title portion. This
// applies to field names.
func untitle(s string) string {
	b := []byte(s)

	for i, c := range b {
		if c >= 65 && c <= 90 {
			b[i] = c + 32
		} else {
			break
		}
	}

	return string(b)
}

// Reflect takes a value and returns a set of partially defined facts
// containing attribute and value components. Currently, struct values or
// pointers to struct values are supported. Struct fields with a primitive
// type are included (support for pointers coming soon). The `origins` tag
// can be used to specify an alternate identity name, an attribute domain
// name or omit the field all together.
func Reflect(v interface{}) (Facts, error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	// Deference pointers.
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = reflect.Indirect(val)
	}

	// Only structs are supported.
	if typ.Kind() != reflect.Struct {
		logrus.Panicf("origins: facts cannot be generated from %s types", typ.Kind())
	}

	// Number of fields on the struct.
	num := typ.NumField()

	var (
		sf    reflect.StructField
		fv    reflect.Value
		toks  []string
		facts []*Fact
		err   error

		attrDomain, attrName   string
		valueDomain, valueName string
		attr, value            *Ident
	)

	// Iterate fields by index.
	for i := 0; i < num; i++ {
		sf = typ.Field(i)

		// Non-empty package denotes an unexported field.
		if sf.PkgPath != "" {
			logrus.Debugf("origins: skipping unexported field %s", sf.Name)
			continue
		}

		attrDomain = ""
		attrName = untitle(sf.Name)
		valueDomain = ""
		valueName = ""

		// Evaluate tag values.
		toks = strings.Split(sf.Tag.Get("origins"), ",")

		if len(toks) > 0 {
			// Omit the field
			if toks[0] == "-" {
				continue
			}

			// First token is attribute name.
			if toks[0] != "" {
				attrName = toks[0]
			}

			// Explicit attribute domain.
			if len(toks) > 1 && toks[1] != "" {
				attrDomain = toks[1]
			}

			// Explicit value domain for reference.
			if len(toks) > 2 && toks[2] != "" {
				valueDomain = toks[2]
			}

			if len(toks) > 3 {
				logrus.Panic("attribute name, attribute domain, value domain tags are supported")
			}
		}

		fv = val.Field(i)

		// If the value implements the stringer inteface use it.
		switch x := fv.Interface().(type) {
		case time.Time:
			valueName = chrono.Format(x)

		case fmt.Stringer:
			valueName = x.String()

		// Evaluate primitive types.
		default:
			// Only pritimive types are supported.
			switch sf.Type.Kind() {
			case reflect.String:
				valueName = fv.String()

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				valueName = fmt.Sprint(fv.Int())

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				valueName = fmt.Sprint(fv.Uint())

			case reflect.Float32, reflect.Float64:
				valueName = fmt.Sprint(fv.Float())

			case reflect.Bool:
				valueName = fmt.Sprint(fv.Bool())

			case reflect.Complex64, reflect.Complex128:
				valueName = fmt.Sprint(fv.Complex())

			default:
				logrus.Debugf("origins: skipping unsupported field %s (%s type)", sf.Name, sf.Type.Kind())
				continue
			}
		}

		if attr, err = NewIdent(attrDomain, attrName); err != nil {
			return nil, err
		}

		if value, err = NewIdent(valueDomain, valueName); err != nil {
			return nil, err
		}

		facts = append(facts, &Fact{
			Attribute: attr,
			Value:     value,
		})
	}

	return facts, nil
}
