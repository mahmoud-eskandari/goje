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

	if Action == Select {
		query += " " + strings.Join(Columns, ",")
	}

	query += " FROM " + Tablename

	conditions, args, err := SQLConditionBuilder(Tablename, Queries)

	return query + conditions, args, err
}

// SQLConditionBuilder [JOIN WHERE LIMIT OFFSET] ...builder
func SQLConditionBuilder(Tablename string, Queries []QueryInterface) (string, []interface{}, error) {
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

	//Add having
	for _, q := range Queries {
		if q.GetType() == QueryTypeHaving {

			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}

			query += " HAVING " + q.GetQuery()
			args = append(args, q.GetArgs()...)
		}
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

/*
*

	Limit Query

*
*/
type QueryLimit int

func (l QueryLimit) GetType() string {
	return QueryTypeLimit
}

func (l QueryLimit) GetQuery() string {
	return "LIMIT ?"
}

func (l QueryLimit) GetArgs() []interface{} {
	return []interface{}{l}
}

func Limit(limit int) QueryLimit {
	return QueryLimit(limit)
}

/**
	Offset Query
**/

type QueryOffset int

func (l QueryOffset) GetType() string {
	return QueryTypeOffset
}

func (l QueryOffset) GetQuery() string {
	return "OFFSET ?"
}

func (l QueryOffset) GetArgs() []interface{} {
	return []interface{}{l}
}

func Offset(offset int) QueryOffset {
	return QueryOffset(offset)
}

/**
	order Query
**/

type QueryOrder struct {
	query string
	args  []interface{}
}

func (o QueryOrder) GetType() string {
	return QueryTypeOrder
}

func (o QueryOrder) GetQuery() string {
	return o.query
}

func (o QueryOrder) GetArgs() []interface{} {
	return o.args
}

func Order(query string, args ...interface{}) QueryOrder {
	return QueryOrder{
		query: query,
		args:  args,
	}
}

/**
	having Query
**/

type QueryHaving struct {
	query string
	args  []interface{}
}

func (h QueryHaving) GetType() string {
	return QueryTypeHaving
}

func (h QueryHaving) GetQuery() string {
	return h.query
}

func (h QueryHaving) GetArgs() []interface{} {
	return h.args
}

func Having(query string, args ...interface{}) QueryHaving {
	return QueryHaving{
		query: query,
		args:  args,
	}
}

/**
	Where Query
**/

type QueryWhere struct {
	query string
	args  []interface{}
}

func (q QueryWhere) GetType() string {
	return QueryTypeWhere
}

func (q QueryWhere) GetQuery() string {
	return q.query
}

func (q QueryWhere) GetArgs() []interface{} {
	return q.args
}

func Where(query string, args ...interface{}) QueryWhere {
	return QueryWhere{
		query: query,
		args:  args,
	}
}

/**
	Where in Query
**/

type QueryWhereIn struct {
	query string
	args  []interface{}
}

func (q QueryWhereIn) GetType() string {
	return QueryTypeWhereIn
}

func (q QueryWhereIn) GetQuery() string {
	return q.query
}

func (q QueryWhereIn) GetArgs() []interface{} {
	return q.args
}

func WhereIn(columnName string, args ...interface{}) QueryWhereIn {
	return QueryWhereIn{
		query: columnName,
		args:  args,
	}
}

/**
	Join Query Builder
**/

type QueryJoin struct {
	joinType string
	table    string
	on       string
	args     []interface{}
}

func (q QueryJoin) GetType() string {
	return QueryTypeWhere
}

func (q QueryJoin) GetQuery() string {
	var on string
	if q.on != "" {
		on = " ON " + q.on
	}
	return q.joinType + " JOIN " + q.table + on
}

func (q QueryJoin) GetArgs() []interface{} {
	return q.args
}

func InnerJoin(table string, on string, args ...interface{}) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Inner,
		table:    table,
		args:     args,
	}
}

func OuterJoin(table string, on string, args ...interface{}) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Outer,
		table:    table,
		args:     args,
	}
}

func NaturalJoin(table string, on string, args ...interface{}) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Natural,
		table:    table,
		args:     args,
	}
}

func RightJoin(table string, on string, args ...interface{}) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Right,
		table:    table,
		args:     args,
	}
}

func LeftJoin(table string, on string, args ...interface{}) QueryJoin {
	return QueryJoin{
		on:       on,
		joinType: Left,
		table:    table,
		args:     args,
	}
}
