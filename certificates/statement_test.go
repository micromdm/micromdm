package certificates

import (
	"testing"
)

type testpair struct {
	when Where
	then string
}

var tests = []testpair{
	{Where{"field", "value", "="}, "field = 'value'"},
	{Where{"field", 1, "="}, "field = 1"},
	{Where{"field", false, "="}, "field = false"},
	{Where{"field", []string{"foo", "bar"}, "IN"}, "field IN ('foo','bar')"},
	{Where{"field", "%foo%", "LIKE"}, "field LIKE '%foo%'"},
	{Where{"field", "bar", "!="}, "field != 'bar'"},
	{Where{"field", nil, "IS"}, "field IS NULL"},
}

func TestWhere(t *testing.T) {
	for _, test := range tests {
		v := test.when.String()
		if v != test.then {
			t.Error(
				"Expected", test.then,
				"got", v,
			)
		}
	}
}

var waTests = []testpair{
	{WhereAnd{Where{"field", "value", "="}, Where{"field", "bar", "="}}, "field = 'value' AND field = 'bar'"},
}

func TestWhereAnd_String(t *testing.T) {
	for _, test := range waTests {
		v := test.when.String()
		if v != test.then {
			t.Error(
				"Expected", test.then,
				"got", v,
			)
		}
	}
}
