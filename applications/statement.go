package applications

import (
	"fmt"
	"strings"
)

// UUID is a filter that can be added as a parameter to narrow down the list of returned results
type UUID struct {
	UUID string
}

func (p UUID) where() string {
	return fmt.Sprintf("application_uuid = '%s'", p.UUID)
}

type Name struct {
	Name string
}

func (p Name) where() string {
	return fmt.Sprintf("name = '%s'", p.Name)
}

type Version struct {
	Version string
}

func (p Version) where() string {
	return fmt.Sprintf("version = '%s'", p.Version)
}

// whereer is for building args passed into a method which finds resources
type whereer interface {
	where() string
}

// add WHERE clause from params
func addWhereFilters(stmt string, separator string, params ...interface{}) string {
	var where []string
	for _, param := range params {
		if f, ok := param.(whereer); ok {
			where = append(where, f.where())
		}
	}

	if len(where) != 0 {
		whereFilter := strings.Join(where, " "+separator+" ")
		stmt = fmt.Sprintf("%s WHERE %s", stmt, whereFilter)
	}
	return stmt
}

// boolean operators are applied to where conditions which are part of a whereClauseGroup
type booleanOperator string

const (
	OR  = "OR"
	AND = "AND"
)

type whereClauseGroup struct {
	Operator booleanOperator
	Clauses  []whereClause
}

// Get a string representing the where clause
// Second return value is an array of arguments to give to db.Exec etc.
func (cg whereClauseGroup) String() (string, []string) {
	var clauses []string
	var values []string = make([]string, len(cg.Clauses))

	for i, c := range cg.Clauses {
		c.Placeholder = fmt.Sprintf("$%d", i)
		clauses = append(clauses, c.String())
		values = append(values, c.Value)
	}

	return strings.Join(clauses, string(cg.Operator)), values
}

// Struct representation of a where clause. Does not deal with field name escaping or any inference of the value.
// I.E Do your own quoting.
type whereClause struct {
	Operator    string
	Field       string
	Value       string
	Placeholder string
}

func (c whereClause) String() string {
	return fmt.Sprintf(`%s %s %s`, c.Field, c.Operator, c.Value)
}

func Where(field string, operator string, value string) whereClause {
	return whereClause{
		Operator:    operator,
		Field:       field,
		Value:       value,
		Placeholder: "$1",
	}
}

func WhereAnd(clauses ...whereClause) whereClauseGroup {
	return whereClauseGroup{
		Operator: "AND",
		Clauses:  clauses,
	}
}

func WhereOr(clauses ...whereClause) whereClauseGroup {
	return whereClauseGroup{
		Operator: "OR",
		Clauses:  clauses,
	}
}
