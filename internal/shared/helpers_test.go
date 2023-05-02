package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterMap(t *testing.T) {
	testValue := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
	}

	testCases := []struct {
		name     string
		filter   func(string, int) bool
		expected map[string]int
	}{
		{
			name: "Similar",
			filter: func(s string, i int) bool {
				return s == "1" || s == "5"
			},
			expected: map[string]int{
				"1": 1,
				"5": 5,
			},
		},

		{
			name: "EvenValues",
			filter: func(s string, i int) bool {
				return i%2 == 0
			},
			expected: map[string]int{
				"2": 2,
				"4": 4,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := FilterMap(testValue, test.filter)

			assert.Equal(t, test.expected, result)
		})
	}
}

func Test_Merge(t *testing.T) {
	type user struct {
		ID   int64
		Name string
	}

	user1 := user{
		ID: 11,
	}

	user2 := user{
		Name: "Steve",
	}

	result := Assign(&user1, &user2)

	assert.Equal(t, user1.ID, result.ID)
	assert.Equal(t, user2.Name, result.Name)
}
