package protocol

import (
	"bytes"
	"fmt"
	"reflect"
)

// A Rule represents a particular decoding method, usually bound to a type. It
// gives information regarding how to decode and encode a particular piece
// of data to and from a `*bytes.Buffer`.
//
// Only after a call to the `AppliesTo` method returns true can it be said that
// the Rule is allowed to work on a given value.
type Rule interface {
	// AppliesTo predicates which type(s) a particular rule can work upon.
	//
	// Typically AppliesTo only passes for a single type, and it is usually
	// implemented as:
	//
	// ```go
	// func (r RuleImpl) AppliesTo(typ reflect.Type) bool {
	//     return typ.Kind == reflect.SomeKind
	// }
	// ```
	AppliesTo(typ reflect.Type) bool

	// Decode reads from the given *bytes.Buffer and returns the decoded
	// contents of a single field of data. If an error is encountered during
	// read-time, or if the data is invalid post read-time, an error will be
	// thrown.
	Decode(r *bytes.Buffer) (interface{}, error)

	// Encode turns a Go type instance into some bytes packed together in
	// a []byte. If the data is the wrong type, an error will be thrown. If
	// the data is un-encodable, or an error occurs while encoding, the
	// error will be thrown.
	Encode(v interface{}) ([]byte, error)
}

func ErrorMismatchedType(expected string, actual interface{}) error {
	return fmt.Errorf(
		"cannot encode mismatched type (expected: %s, got: %s)",
		expected, reflect.TypeOf(actual),
	)
}
