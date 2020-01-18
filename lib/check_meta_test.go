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

	opts.MetaKey = "dummy"

	for i, tc := range testCase {
		opts.Expected = tc.expected
		chk := checkMetaValue(tc.actual)
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}

func TestCheckStringTypeValue(t *testing.T) {
	testCase := []struct {
		expected string
		actual   interface{}
		status   checkers.Status
		regex    bool
	}{
		{
			expected: "foobar",
			actual:   "foobar",
			status:   checkers.OK,
			regex:    false,
		},
		{
			expected: "hoge",
			actual:   "fuga",
			status:   checkers.CRITICAL,
			regex:    false,
		},
		{
			expected: "foo.*",
			actual:   "foobar",
			status:   checkers.CRITICAL,
			regex:    false,
		},
		{
			expected: "foo.*",
			actual:   "foobar",
			status:   checkers.OK,
			regex:    true,
		},
		{
			expected: "foo[a-z]{3}",
			actual:   "foobar",
			status:   checkers.OK,
			regex:    true,
		},
	}

	opts.MetaKey = "dummy"

	for i, tc := range testCase {
		opts.Expected = tc.expected
		opts.IsRegex = tc.regex
		chk := checkMetaValue(tc.actual)
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}

func TestCheckNumberTypeValue(t *testing.T) {
	testCase := []struct {
		expected string
		actual   interface{}
		status   checkers.Status
		gt       bool
		lt       bool
		ge       bool
		le       bool
	}{
		{
			expected: "1000",
			actual:   float64(1000),
			status:   checkers.OK,
			gt:       false,
			lt:       false,
			ge:       false,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(1001),
			status:   checkers.OK,
			gt:       true,
			lt:       false,
			ge:       false,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(999),
			status:   checkers.OK,
			gt:       false,
			lt:       true,
			ge:       false,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(1000),
			status:   checkers.OK,
			gt:       false,
			lt:       false,
			ge:       true,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(1000),
			status:   checkers.OK,
			gt:       false,
			lt:       false,
			ge:       false,
			le:       true,
		},
		{
			expected: "1000",
			actual:   float64(1001),
			status:   checkers.CRITICAL,
			gt:       false,
			lt:       false,
			ge:       false,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(1000),
			status:   checkers.CRITICAL,
			gt:       true,
			lt:       false,
			ge:       false,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(1000),
			status:   checkers.CRITICAL,
			gt:       false,
			lt:       true,
			ge:       false,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(999),
			status:   checkers.CRITICAL,
			gt:       false,
			lt:       false,
			ge:       true,
			le:       false,
		},
		{
			expected: "1000",
			actual:   float64(1001),
			status:   checkers.CRITICAL,
			gt:       false,
			lt:       false,
			ge:       false,
			le:       true,
		},
	}

	opts.MetaKey = "dummy"

	for i, tc := range testCase {
		opts.Expected = tc.expected
		opts.GreaterThan = tc.gt
		opts.LessThan = tc.lt
		opts.GreaterOrEqual = tc.ge
		opts.LessOrEqual = tc.le
		chk := checkMetaValue(tc.actual)
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}
