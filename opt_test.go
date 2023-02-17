package opt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Int    int
	String string
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
