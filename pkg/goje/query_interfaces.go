package goje

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
	Group Query
**/

type QueryGroup struct {
	query string
	args  []interface{}
}

func (h QueryGroup) GetType() string {
	return QueryTypeGroup
}

func (h QueryGroup) GetQuery() string {
	return h.query
}

func (h QueryGroup) GetArgs() []interface{} {
	return h.args
}

func GroupBy(query string, args ...interface{}) QueryGroup {
	return QueryGroup{
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
	return QueryTypeJoin
}

func (q QueryJoin) GetQuery() string {
	var on string
	if q.on != "" {
		on = " ON " + q.on
	}
	return " " + q.joinType + " JOIN " + q.table + on + " "
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
