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
