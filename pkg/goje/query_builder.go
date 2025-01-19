package goje

import (
	"errors"
	"strings"
)

const (
	QueryTypeLimit   = "limit"
	QueryTypeOffset  = "offset"
	QueryTypeOrder   = "order"
	QueryTypeHaving  = "having"
	QueryTypeGroup   = "group"
	QueryTypeWhere   = "where"
	QueryTypeWhereIn = "where in"
	QueryTypeJoin    = "join"
)

type QueryInterface interface {
	GetType() string
	GetQuery() string
	GetArgs() []interface{}
}

const (
	Select string = "SELECT"
	Insert string = "INSERT"
	Update string = "UPDATE"
	Delete string = "DELETE"

	Left    string = "LEFT"
	Right   string = "RIGHT"
	Inner   string = "INNER"
	Outer   string = "OUTER"
	Natural string = "NATURAL"
)

// make a select ... FROM query
func SelectQueryBuilder(Tablename string, Columns []string, Queries []QueryInterface) (string, []interface{}, error) {
	return ArgumentLessQueryBuilder(Select, Tablename, Columns, Queries)
}

// ArgumentLessQueryBuilder (Select, Delete) query builder
func ArgumentLessQueryBuilder(Action, Tablename string, Columns []string, Queries []QueryInterface) (string, []interface{}, error) {

	if Action != Select && Action != Delete {
		return "", nil, errors.New("this function dosen't support: " + Action)
	}

	query := Action

	Columns = columnsFilter(Columns)
	if Action == Select {
		query += " " + strings.Join(Columns, ",") + " "
	}

	query += " FROM " + Tablename

	conditions, args, err := SQLConditionBuilder(Queries)

	return query + conditions, args, err
}

// SQLConditionBuilder [JOIN WHERE LIMIT OFFSET] ...builder
func SQLConditionBuilder(Queries []QueryInterface) (string, []interface{}, error) {
	query := " "
	var args []interface{}
	var where []string

	//Produce Joins
	for _, q := range Queries {
		if q.GetType() == QueryTypeJoin {
			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}
			query += q.GetQuery()
			args = append(args, q.GetArgs()...)
		}
	}

	//Produce Where condition
	for _, q := range Queries {
		if q.GetType() == QueryTypeWhere {
			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}
			where = append(where, q.GetQuery())
			args = append(args, q.GetArgs()...)
		}
	}

	//Produce Where in condition
	for _, q := range Queries {
		if q.GetType() == QueryTypeWhereIn {
			if len(q.GetArgs()) == 0 {
				return "", nil, errors.New(q.GetQuery() + "; should have at least one argument")
			}

			bindParams := strings.Repeat(",?", len(q.GetArgs()))

			where = append(where, q.GetQuery()+" IN("+bindParams[1:]+")")
			args = append(args, q.GetArgs()...)
		}
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	//Add group by
	var groupbys []string
	for _, q := range Queries {
		if q.GetType() == QueryTypeGroup {
			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}
			groupbys = append(groupbys, columnFilter(q.GetQuery()))
			args = append(args, q.GetArgs()...)
		}

	}
	if len(groupbys) > 0 {
		query += " GROUP BY " + strings.Join(groupbys, ",")
	}

	//Add having
	var havings []string
	for _, q := range Queries {
		if q.GetType() == QueryTypeHaving {
			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}
			havings = append(havings, columnFilter(q.GetQuery()))
			args = append(args, q.GetArgs()...)
		}
	}

	if len(havings) > 0 {
		query += " HAVING " + strings.Join(havings, " AND ")
	}

	//Add orders
	for _, q := range Queries {
		if q.GetType() == QueryTypeOrder {

			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}

			query += " ORDER BY " + q.GetQuery()
			args = append(args, q.GetArgs()...)
		}
	}

	//Add limitations
	for _, q := range Queries {
		if q.GetType() == QueryTypeLimit || q.GetType() == QueryTypeOffset {
			query += " " + q.GetQuery()
			args = append(args, q.GetArgs()...)
		}
	}

	return query, args, nil
}

// filter multiple columns
func columnsFilter(in []string) []string {
	for i := range in {
		in[i] = columnFilter(in[i])
	}
	return in
}

// filter column: add backtik if needed
func columnFilter(v string) string {
	if strings.Contains(v, "`") ||
		strings.Contains(v, " ") ||
		strings.Contains(v, "(") ||
		strings.Contains(v, "*") ||
		strings.Contains(v, "%") ||
		strings.Contains(v, ".") {
		return v
	}
	return "`" + v + "`"
}
