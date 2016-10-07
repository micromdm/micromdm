package certificate

import (
	"testing"
)

var whereTests = []struct {
	when Where
	then string
}{
	{Where{"field", "value", "="}, "field = 'value'"},
	{Where{"field", 1, "="}, "field = 1"},
	{Where{"field", false, "="}, "field = false"},
	{Where{"field", []string{"foo", "bar"}, "IN"}, "field IN ('foo','bar')"},
	{Where{"field", "%foo%", "LIKE"}, "field LIKE '%foo%'"},
	{Where{"field", "bar", "!="}, "field != 'bar'"},
	{Where{"field", nil, "IS"}, "field IS NULL"},
}

func TestWhere_String(t *testing.T) {
	for _, test := range whereTests {
		v := test.when.String()
		if v != test.then {
			t.Error(
				"Expected", test.then,
				"got", v,
			)
		}
	}
}

var waTests = []struct {
	when WhereAnd
	then string
}{
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

var woTests = []struct {
	when WhereOr
	then string
}{
	{WhereOr{Where{"field", "value", "="}, Where{"field", "bar", "="}}, "field = 'value' OR field = 'bar'"},
}

func TestWhereOr_String(t *testing.T) {
	for _, test := range woTests {
		v := test.when.String()
		if v != test.then {
			t.Error(
				"Expected", test.then,
				"got", v,
			)
		}
	}
}
