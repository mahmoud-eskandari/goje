package goje

import (
	"context"
	"database/sql"
	"errors"
)

// Context and Database Handler struct
type ContextHandler struct {
	DB  *sql.DB
	Ctx context.Context
}

// Entity All generated structs implemented from Entity
type Entity interface {
	TableName() string
	Columns() []string
	GetCtx() *ContextHandler
	SetCtx(*ContextHandler)
	GetParent() *Entity
}

var (
	ErrHandlerIsNil        = errors.New("context handler dosen't set properly")
	ErrRecursiveLoad       = errors.New("recursive load is forbidden")
	ErrNoColsSetForUpdate  = errors.New("cols should have at least one proprty for update")
	ErrNoRowsForInsert     = errors.New("there isn't any row for insert into database")
	ErrNoRowsColsForInsert = errors.New("cols should have at least one proprty for update")
)
