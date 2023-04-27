package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {
	raw := []byte(`{'ds:0' : {id:'T6C6Se',request:[]},'ds:1' : {id:'Qu3lde',request:[null,true]},'ds:2' : {id:'X5bnOb',ext: 1.59168981E8 ,request:[]}}`)

	result, err := parseServiceData(raw)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(
		t,
		containsAll(
			keys(result),
			[]string{"ds:0", "ds:1", "ds:2"}...,
		),
	)
}

func keys[K comparable, V any, M ~map[K]V](m M) []K {
	keys := make([]K, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func containsAll[V comparable, S ~[]V](slice S, entries ...V) (contains bool) {
	for _, value := range slice {
		contains = eq(value, entries...)
	}

	return
}

func eq[V comparable](value V, entries ...V) bool {
	for _, entry := range entries {
		if value == entry {
			return true
		}
	}

	return false
}
