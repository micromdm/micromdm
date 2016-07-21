package certificates

import (
	"fmt"
	"strconv"
	"strings"
)

// Where struct wraps a where clause.
// Basic usage is Where{"field","value"} for WHERE field = value
// The Operator field can be used if you need < > = != LIKE / NOT LIKE operators
type Where struct {
	Field    string
	Value    interface{}
	Operator string
}

// Stringer produces the WHERE condition
func (w Where) String() string {
	var operator string = w.Operator
	if w.Operator == "" {
		operator = "="
	}

	var quotedValue string
	switch w.Value.(type) {
	case string:
		quotedValue = fmt.Sprintf("'%s'", w.Value)
	case nil:
		operator = "IS"
		quotedValue = "NULL"
	case bool:
		if w.Value.(bool) == true {
			quotedValue = "true"
		} else {
			quotedValue = "false"
		}
	case []string: // IN("strings...")
		operator = "IN"
		inValues := w.Value.([]string)
		quotedValue = "('" + strings.Join(inValues, "','") + "')"
	case int:
		quotedValue = strconv.Itoa(w.Value.(int))
	}

	return fmt.Sprintf("%s %s %s", w.Field, operator, quotedValue)
}

// whereer is for building args passed into a method which finds resources
type whereer interface {
	where() string
}

type WhereAnd []Where

func (wa WhereAnd) String() {
	return strings.Join(wa, " AND ")
}

type WhereOr []Where

func (wo WhereOr) String() {
	return strings.Join(wo, " OR ")
}
