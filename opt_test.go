package opt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Int    int
	String string
}

func TestOpt(t *testing.T) {
	t.Run("copies are deep", func(t *testing.T) {
		type Inner struct {
			Bool bool
		}

		type Outer struct {
			Inner Opt[Inner]
		}

		// Arrange
		expected := Outer{
			Inner: New(Inner{}),
		}
		actual := Outer{
			Inner: New(Inner{}),
		}

		// Act
		copy := actual
		copy.Inner.Value.Bool = true

		// Assert
		require.NotEqual(t, copy, actual)
		require.Equal(t, expected, actual)
	})

	t.Run("pointers are not", func(t *testing.T) {
		type Inner struct {
			Bool bool
		}

		type Outer struct {
			Inner *Inner
		}

		// Arrange
		expected := Outer{
			Inner: &Inner{},
		}
		actual := Outer{
			Inner: &Inner{},
		}

		// Act
		copy := actual
		copy.Inner.Bool = true

		// Assert
		require.Equal(t, copy, actual)
		require.NotEqual(t, expected, actual)
	})
}

func TestNew(t *testing.T) {
	// Arrange
	value := TestStruct{
		Int:    42,
		String: "Hello, World!",
	}

	// Act
	opt := New(value)

	// Assert
	require.True(t, opt.IsSet)
	require.Equal(t, value, opt.Value)
}

func TestOf(t *testing.T) {
	t.Run("with nil pointer", func(t *testing.T) {
		// Arrange
		var ptr *TestStruct

		// Act
		opt := Of(ptr)

		// Assert
		require.False(t, opt.IsSet)
		require.Zero(t, opt.Value)
	})

	t.Run("with non-nil pointer", func(t *testing.T) {
		// Arrange
		ptr := &TestStruct{
			Int:    42,
			String: "Hello, World!",
		}

		// Act
		opt := Of(ptr)

		// Assert
		require.True(t, opt.IsSet)
		require.Equal(t, *ptr, opt.Value)
	})
}

func TestOpt_Get(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// Arrange
		expected := TestStruct{
			Int:    42,
			String: "Hello, World!",
		}
		opt := New(expected)

		// Act
		actual, ok := opt.Get()

		// Assert
		require.True(t, ok)
		require.Equal(t, expected, actual)
	})

	t.Run("without value", func(t *testing.T) {
		// Arrange
		var opt Opt[TestStruct]

		// Act
		actual, ok := opt.Get()

		// Assert
		require.False(t, ok)
		require.Zero(t, actual)
	})
}

func TestOpt_Or(t *testing.T) {
	t.Run("filled value empty other", func(t *testing.T) {
		// Arrange
		expected := New(42)
		opt := New(42)
		var other Opt[int]

		// Act
		actual := opt.Or(other)

		// Assert
		require.Equal(t, expected, actual)
	})

	t.Run("filled value filled other", func(t *testing.T) {
		// Arrange
		expected := New(42)
		opt := New(42)
		other := New(24)

		// Act
		actual := opt.Or(other)

		// Assert
		require.Equal(t, expected, actual)
	})

	t.Run("empty value empty other", func(t *testing.T) {
		// Arrange
		var expected Opt[int]
		var opt Opt[int]
		var other Opt[int]

		// Act
		actual := opt.Or(other)

		// Assert
		require.Equal(t, expected, actual)
	})

	t.Run("empty value filled other", func(t *testing.T) {
		// Arrange
		expected := New(24)
		var opt Opt[int]
		other := New(24)

		// Act
		actual := opt.Or(other)

		// Assert
		require.Equal(t, expected, actual)
	})
}

func TestOpt_OrValue(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// Arrange
		expected := 42
		opt := New(42)

		// Act
		actual := opt.OrValue(24)

		// Assert
		require.Equal(t, expected, actual)
	})

	t.Run("without value", func(t *testing.T) {
		// Arrange
		expected := 24
		var opt Opt[int]

		// Act
		actual := opt.OrValue(24)

		// Assert
		require.Equal(t, expected, actual)
	})
}

func TestOpt_Ptr(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// Arrange
		expected := TestStruct{
			Int:    42,
			String: "Hello, World!",
		}
		opt := New(expected)

		// Act
		actual := opt.Ptr()

		// Assert
		require.NotNil(t, actual)
		require.Equal(t, expected, *actual)
	})

	t.Run("without value", func(t *testing.T) {
		// Arrange
		var opt Opt[TestStruct]

		// Act
		actual := opt.Ptr()

		// Assert
		require.Nil(t, actual)
	})
}

func TestOpt_Set(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// Arrange
		expected := TestStruct{
			Int:    24,
			String: "!dlorW ,olleH",
		}
		opt := New(TestStruct{
			Int:    42,
			String: "Hello, World!",
		})

		// Act
		opt.Set(expected)

		// Assert
		require.True(t, opt.IsSet)
		require.Equal(t, expected, opt.Value)
	})

	t.Run("without value", func(t *testing.T) {
		// Arrange
		expected := TestStruct{
			Int:    42,
			String: "Hello, World!",
		}
		var opt Opt[TestStruct]

		// Act
		opt.Set(expected)

		// Assert
		require.True(t, opt.IsSet)
		require.Equal(t, expected, opt.Value)
	})
}

func TestOpt_Reset(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// Arrange
		opt := New(TestStruct{
			Int:    42,
			String: "Hello, World!",
		})

		// Act
		opt.Reset()

		// Assert
		require.False(t, opt.IsSet)
		require.Zero(t, opt.Value)
	})

	t.Run("without value", func(t *testing.T) {
		// Arrange
		var opt Opt[TestStruct]

		// Act
		opt.Reset()

		// Assert
		require.False(t, opt.IsSet)
		require.Zero(t, opt.Value)
	})
}

func TestMap(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// Arrange
		expected := 42
		opt := New(TestStruct{
			Int:    expected,
			String: "Hello, World",
		})

		// Act
		actual := Map(opt, func(v TestStruct) int {
			return v.Int
		})

		// Assert
		require.True(t, actual.IsSet)
		require.Equal(t, expected, actual.Value)
	})

	t.Run("without value", func(t *testing.T) {
		// Arrange
		var opt Opt[TestStruct]

		// Act
		actual := Map(opt, func(v TestStruct) int {
			return v.Int
		})

		// Assert
		require.False(t, actual.IsSet)
		require.Zero(t, actual.Value)
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("with value returns value", func(t *testing.T) {
		// Arrange
		expected := 42
		opt := New(TestStruct{
			Int:    expected,
			String: "Hello, World!",
		})

		// Act
		actual := FlatMap(opt, func(v TestStruct) Opt[int] {
			return New(v.Int)
		})

		// Assert
		require.True(t, actual.IsSet)
		require.Equal(t, expected, actual.Value)
	})

	t.Run("without value returns value", func(t *testing.T) {
		// Arrange
		var opt Opt[TestStruct]

		// Act
		actual := FlatMap(opt, func(v TestStruct) Opt[int] {
			return New(v.Int)
		})

		// Assert
		require.False(t, actual.IsSet)
		require.Zero(t, actual.Value)
	})

	t.Run("with value returns nothing", func(t *testing.T) {
		// Arrange
		opt := New(TestStruct{
			Int:    42,
			String: "Hello, World!",
		})

		// Act
		actual := FlatMap(opt, func(v TestStruct) Opt[int] {
			return Opt[int]{}
		})

		// Assert
		require.False(t, actual.IsSet)
		require.Zero(t, actual.Value)
	})

	t.Run("without value returns nothing", func(t *testing.T) {
		// Arrange
		var opt Opt[TestStruct]

		// Act
		actual := FlatMap(opt, func(v TestStruct) Opt[int] {
			return Opt[int]{}
		})

		// Assert
		require.False(t, actual.IsSet)
		require.Zero(t, actual.Value)
	})
}
