// Package opt provides optional values sans pointers.
package opt

// Opt is an optional value.
// The zero value represents an empty optional value.
type Opt[T any] struct {
	// IsSet is true when the Opt contains a value.
	IsSet bool
	// Value is the value contained in the Opt, when IsSet is true.
	Value T
}

// New returns an Opt containing a value.
func New[T any](value T) Opt[T] {
	return Opt[T]{
		IsSet: true,
		Value: value,
	}
}

// Of turns a pointer into an Opt.
// Passing a nil pointer returns an empty value.
func Of[T any](ptr *T) Opt[T] {
	var res Opt[T]

	if ptr != nil {
		res.Set(*ptr)
	}

	return res
}

// Get returns the inner value, and if the Opt contains a value.
func (o Opt[T]) Get() (value T, ok bool) {
	return o.Value, o.IsSet
}

// Ptr turns this Opt into a pointer.
func (o Opt[T]) Ptr() *T {
	if o.IsSet {
		return &o.Value
	}

	return nil
}

// Set places a value in the Opt.
func (o *Opt[T]) Set(value T) {
	o.Value = value
	o.IsSet = true
}

// Reset empties the Opt.
func (o *Opt[T]) Reset() {
	var zero T
	o.Value = zero
	o.IsSet = false
}

// Map performs an operation on the contained value.
func Map[T, U any](in Opt[T], f func(v T) U) Opt[U] {
	var out Opt[U]

	if in.IsSet {
		out.Set(f(in.Value))
	}

	return out
}

// FlatMap performs an operation on the contained value, which returns an Opt.
func FlatMap[T, U any](in Opt[T], f func(v T) Opt[U]) Opt[U] {
	if in.IsSet {
		return f(in.Value)
	}

	return Opt[U]{}
}
