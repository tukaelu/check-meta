package checkmeta

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func initialize() {
	opts.Namespace = "dummy"
	opts.MetaKey = "dummy"
	opts.Expected = ""
	opts.IsRegex = false
	opts.GreaterThan = false
	opts.LessThan = false
	opts.GreaterOrEqual = false
	opts.LessOrEqual = false
	opts.CompareNamespace = ""
	opts.CompareMetaKey = ""
	opts.compareMetaValue = nil
}

func TestCheckStringTypeValue(t *testing.T) {
	testCase := []struct {
		expected         string
		actual           interface{}
		status           checkers.Status
		regex            bool
		compareMeta      bool
		compareMetaValue string
	}{
		{
			expected:         "foobar",
			actual:           "foobar",
			status:           checkers.OK,
			regex:            false,
			compareMeta:      false,
			compareMetaValue: "",
		},
		{
			expected:         "hoge",
			actual:           "fuga",
			status:           checkers.CRITICAL,
			regex:            false,
			compareMeta:      false,
			compareMetaValue: "",
		},
		{
			expected:         "foo.*",
			actual:           "foobar",
			status:           checkers.CRITICAL,
			regex:            false,
			compareMeta:      false,
			compareMetaValue: "",
		},
		{
			expected:         "foo.*",
			actual:           "foobar",
			status:           checkers.OK,
			regex:            true,
			compareMeta:      false,
			compareMetaValue: "",
		},
		{
			expected:         "foo[a-z]{3}",
			actual:           "foobar",
			status:           checkers.OK,
			regex:            true,
			compareMeta:      false,
			compareMetaValue: "",
		},
		{
			expected:         "",
			actual:           "foobar",
			status:           checkers.OK,
			regex:            false,
			compareMeta:      true,
			compareMetaValue: "foobar",
		},
	}

	for i, tc := range testCase {
		initialize()
		opts.Expected = tc.expected
		opts.IsRegex = tc.regex
		if tc.compareMeta {
			opts.CompareMetaKey = opts.MetaKey
			opts.compareMetaValue = tc.compareMetaValue
		}
		chk := checkMetadata(tc.actual)
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}

func TestCheckNumberTypeValue(t *testing.T) {
	testCase := []struct {
		expected         string
		actual           interface{}
		status           checkers.Status
		gt               bool
		lt               bool
		ge               bool
		le               bool
		compareMeta      bool
		compareMetaValue float64
	}{
		{
			expected:         "1000",
			actual:           float64(1000),
			status:           checkers.OK,
			gt:               false,
			lt:               false,
			ge:               false,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1001),
			status:           checkers.OK,
			gt:               true,
			lt:               false,
			ge:               false,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(999),
			status:           checkers.OK,
			gt:               false,
			lt:               true,
			ge:               false,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1000),
			status:           checkers.OK,
			gt:               false,
			lt:               false,
			ge:               true,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1000),
			status:           checkers.OK,
			gt:               false,
			lt:               false,
			ge:               false,
			le:               true,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1001),
			status:           checkers.CRITICAL,
			gt:               false,
			lt:               false,
			ge:               false,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1000),
			status:           checkers.CRITICAL,
			gt:               true,
			lt:               false,
			ge:               false,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1000),
			status:           checkers.CRITICAL,
			gt:               false,
			lt:               true,
			ge:               false,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(999),
			status:           checkers.CRITICAL,
			gt:               false,
			lt:               false,
			ge:               true,
			le:               false,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1001),
			status:           checkers.CRITICAL,
			gt:               false,
			lt:               false,
			ge:               false,
			le:               true,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "1000",
			actual:           float64(1001),
			status:           checkers.UNKNOWN,
			gt:               false,
			lt:               false,
			ge:               true,
			le:               true,
			compareMeta:      false,
			compareMetaValue: 0,
		},
		{
			expected:         "",
			actual:           float64(1000),
			status:           checkers.OK,
			gt:               false,
			lt:               false,
			ge:               false,
			le:               false,
			compareMeta:      true,
			compareMetaValue: 1000,
		},
	}

	for i, tc := range testCase {
		initialize()
		opts.Expected = tc.expected
		opts.GreaterThan = tc.gt
		opts.LessThan = tc.lt
		opts.GreaterOrEqual = tc.ge
		opts.LessOrEqual = tc.le
		if tc.compareMeta {
			opts.CompareMetaKey = opts.MetaKey
			opts.compareMetaValue = tc.compareMetaValue
		}
		chk := checkMetadata(tc.actual)
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}

func TestCheckBooleanTypeValue(t *testing.T) {
	testCase := []struct {
		expected         string
		actual           interface{}
		status           checkers.Status
		compareMeta      bool
		compareMetaValue bool
	}{
		{
			expected:         "true",
			actual:           bool(true),
			status:           checkers.OK,
			compareMeta:      false,
			compareMetaValue: false,
		},
		{
			expected:         "false",
			actual:           bool(false),
			status:           checkers.OK,
			compareMeta:      false,
			compareMetaValue: false,
		},
		{
			expected:         "true",
			actual:           bool(false),
			status:           checkers.CRITICAL,
			compareMeta:      false,
			compareMetaValue: false,
		},
		{
			expected:         "true",
			actual:           nil,
			status:           checkers.UNKNOWN,
			compareMeta:      false,
			compareMetaValue: false,
		},
		{
			expected:         "",
			actual:           bool(true),
			status:           checkers.OK,
			compareMeta:      true,
			compareMetaValue: true,
		},
	}

	for i, tc := range testCase {
		initialize()
		opts.Expected = tc.expected
		if tc.compareMeta {
			opts.CompareMetaKey = opts.MetaKey
			opts.compareMetaValue = tc.compareMetaValue
		}
		chk := checkMetadata(tc.actual)
		assert.Equal(t, chk.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}
