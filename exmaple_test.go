package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Add(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		val1           int
		val2           int
		expectedResult int
	}{
		{
			name:           "add correct",
			val1:           2,
			val2:           4,
			expectedResult: 6,
		},
		{
			name:           "add negative",
			val1:           -2,
			val2:           -4,
			expectedResult: -6,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expectedResult, Add(tc.val1, tc.val2))
		})
	}
}
