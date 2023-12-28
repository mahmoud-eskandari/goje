package goje

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type ContextHandler struct {
	DB  *sql.DB
	Ctx context.Context
}

type Entity interface {
	TableName() string
	Columns() []string
	GetCtx() *ContextHandler
	GetParent() *Entity
}

var (
	ErrHandlerIsNil  = errors.New("context handler dosen't set properly")
	ErrRecursiveLoad = errors.New("recursive load is forbidden")
)

const (
	QueryTypeLimit   = "limit"
	QueryTypeOffset  = "offset"
	QueryTypeWhere   = "where"
	QueryTypeWhereIn = "where in"
)

type QueryInterface interface {
	GetType() string
	GetQuery() string
	GetArgs() []interface{}
}

func QueryBuilder(Tablename string, Columns []string, Queries []QueryInterface) (string, []interface{}, error) {
	query := "SELECT " + strings.Join(Columns, ",") + " FROM " + Tablename
	var args []interface{}
	var where []string

	for _, q := range Queries {
		if q.GetType() == QueryTypeWhere {
			if strings.Count(q.GetQuery(), "?") != len(q.GetArgs()) {
				return "", nil, errors.New(q.GetQuery() + "; args dosen't match with binds `?`")
			}
			where = append(where, q.GetQuery())
			args = append(args, q.GetArgs()...)
		}
	}

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
