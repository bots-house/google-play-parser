package ramda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Path(t *testing.T) {
	testValue := map[string]any{
		"a": 11,
		"b": []any{
			map[string][]any{
				"c": {3},
			},
		},
		"a:1": map[string]any{
			"11": 15,
		},
	}

	testCases := []struct {
		name     string
		path     []any
		expected any
	}{
		{
			name:     "OneKey",
			expected: 11,
			path:     []any{"a"},
		},

		{
			name:     "...",
			expected: 15,
			path:     []any{"a:1", "11"},
		},

		{
			name:     "Deeper",
			expected: 3,
			path:     []any{"b", 0, "c", 0},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := Path(test.path, testValue)

			assert.Equal(t, test.expected, result)
		})
	}
}

func Test_Map(t *testing.T) {
	fn := func(i int) int { return i * 2 }

	t.Run("Slice", func(t *testing.T) {
		arr := []int{2, 3, 4}

		result := Map(arr, fn)

		t.Log(result)
	})

	t.Run("Map", func(t *testing.T) {
		m := map[string]int{
			"x": 2,
			"y": 3,
		}

		result := Map(m, fn)

		t.Log(result)
	})

	t.Run("Struct", func(t *testing.T) {
		input := struct {
			X int
			Y int
		}{
			X: 10,
			Y: 20,
		}

		result := Map(&input, fn)

		t.Log(result)
	})
}
