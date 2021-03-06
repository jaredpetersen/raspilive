package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestTransformError(t *testing.T) {
	testCases := []struct {
		err      error
		expected error
	}{
		{
			err:      errors.New("required key RASPILIVE_HLS_PORT missing value"),
			expected: errors.New("Required config key RASPILIVE_HLS_PORT is missing"),
		},
		{
			err:      errors.New("something went wrong"),
			expected: errors.New("something went wrong"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.err), func(t *testing.T) {
			actual := transformError(tc.err)

			if errors.Is(actual, tc.expected) {
				t.Error("Transformed error is not correct, got", actual)
			}
		})
	}

}
