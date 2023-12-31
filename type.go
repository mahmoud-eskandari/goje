package goje

import (
	"context"
	"database/sql"
	"errors"
)

type QueryAble interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Context and Database Handler struct
type Context struct {
	DB  QueryAble
	Ctx context.Context
	Tx  bool
}

func (c Context) Commit() error {
	if !c.Tx {
		return ErrIsntATx
	}
	tx, ok := c.DB.(*sql.Tx)
	if !ok {
		return ErrTxIsntSet
	}

	err := tx.Commit()
	if err != nil {
		return err
	}
	c.Ctx.Done()
	return nil
}

func (c Context) Rollback() error {
	if !c.Tx {
		return ErrIsntATx
	}
	tx, ok := c.DB.(*sql.Tx)
	if !ok {
		return ErrTxIsntSet
	}

	err := tx.Rollback()
	if err != nil {
		return err
	}
	c.Ctx.Done()
	return nil
}

// Entity All generated structs implemented from Entity
type Entity interface {
	TableName() string
	Columns() []string
	GetCtx() *Context
	SetCtx(*Context)
	GetParent() *Entity
}

var (
	ErrHandlerIsNil        = errors.New("context handler dosen't set properly")
	ErrRecursiveLoad       = errors.New("recursive load is forbidden")
	ErrNoColsSetForUpdate  = errors.New("cols should have at least one proprty for update")
	ErrNoRowsForInsert     = errors.New("there isn't any row for insert into database")
	ErrNoRowsColsForInsert = errors.New("cols should have at least one proprty for update")
	ErrUnknownDBDriver     = errors.New("goje doesn't support this driver")
	ErrIsntATx             = errors.New("it isn't a transactional context")
	ErrTxIsntSet           = errors.New("there is not any transaction context")
)
