package main

import (
	"fmt"
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
		fmt.Println(tc.name)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expectedResult, Add(tc.val1, tc.val2))
		})
	}
}

func Test_Add2(p *testing.T) {
	p.Parallel()

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

	p.Run(testCases[0].name, func(t *testing.T) {
		assert.Equal(t, testCases[0].expectedResult, Add(testCases[0].val1, testCases[0].val2))
	})

	p.Run(testCases[1].name, func(t *testing.T) {
		assert.Equal(t, testCases[1].expectedResult, Add(testCases[1].val1, testCases[1].val2))
	})
}