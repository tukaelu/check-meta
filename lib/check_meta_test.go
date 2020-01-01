package checkmeta

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestCheckMetaValue(t *testing.T) {

	testCase := []struct {
		expected string
		actual   interface{}
		status   checkers.Status
	}{
		{
			expected: "foobar",
			actual:   "foobar",
			status:   checkers.OK,
		},
		{
			expected: "hoge",
			actual:   "fuga",
			status:   checkers.CRITICAL,
		},
		{
			expected: "1000",
			actual:   float64(1000),
			status:   checkers.OK,
		},
		{
			expected: "1000",
			actual:   float64(1002),
			status:   checkers.CRITICAL,
		},
		{
			expected: "true",
			actual:   bool(true),
			status:   checkers.OK,
		},
		{
			expected: "false",
			actual:   bool(false),
			status:   checkers.OK,
		},
		{
			expected: "true",
			actual:   bool(false),
			status:   checkers.CRITICAL,
		},
		{
			expected: "true",
			actual:   nil,
			status:   checkers.UNKNOWN,
		},
	}

	for i, tc := range testCase {
		chk := checkMetaValue(tc.expected, tc.actual, "dummy")
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}
