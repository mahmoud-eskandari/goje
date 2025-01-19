package goje

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Default Database Connection that fill by
var DefatultDB *sql.DB

// ConnectDB Connect and set default database
func InitDB(conn *DBConfig) error {
	db, err := NewDBConnection(conn)
	if err != nil {
		return err
	}
	DefatultDB = db
	return nil
}

// ConnectDB Connect and return database
func NewDBConnection(conn *DBConfig) (*sql.DB, error) {
	if conn.Driver != "mysql" {
		return nil, ErrUnknownDBDriver
	}

	db, err := sql.Open(conn.Driver, conn.String())

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(conn.MaxIdleTime)
	db.SetMaxIdleConns(conn.MaxIdleConns)
	db.SetMaxOpenConns(conn.MaxOpenConns)
	db.SetConnMaxLifetime(conn.ConnMaxLifetime)

	return db, nil
}

// GetHandler make a handler from default database and a TODO context
func GetHandler() *Context {
	return &Context{
		Ctx: context.TODO(),
		DB:  DefatultDB,
	}
}

// MakeHandler make a handler from default database
func MakeHandler(ctx context.Context) *Context {
	return &Context{
		Ctx: ctx,
		DB:  DefatultDB,
	}
}

// DefaultHandler make a handler from default database
func MakeTxHandler(ctx context.Context, options *sql.TxOptions) (*Context, error) {
	tx, err := DefatultDB.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}

	return &Context{
		Ctx: ctx,
		DB:  tx,
		Tx:  true,
	}, nil
}

/*
Goje database config schema
# yaml example
driver: mysql
host: 127.0.0.1
port: 3306
user: root
password:
schema: mydbname
*/
type DBConfig struct {
	Driver   string            `json:"driver" yaml:"driver"`
	Host     string            `json:"host" yaml:"host"`
	Port     int               `json:"port" yaml:"port"`
	User     string            `json:"user" yaml:"user"`
	Password string            `json:"password" yaml:"password"`
	Schema   string            `json:"schema" yaml:"schema"`
	Flags    map[string]string `json:"flags" yaml:"flags"`

	MaxIdleTime     time.Duration `json:"MaxIdleTime" yaml:"MaxIdleTime"`
	MaxOpenConns    int           `json:"MaxOpenConns" yaml:"MaxOpenConns"`
	MaxIdleConns    int           `json:"MaxIdleConns" yaml:"MaxIdleConns"`
	ConnMaxLifetime time.Duration `json:"ConnMaxLifetime" yaml:"ConnMaxLifetime"`
}

func (db DBConfig) String() string {
	flags := ""
	for k, v := range db.Flags {
		flags += "&" + k + "=" + url.QueryEscape(v)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", db.User, db.Password, db.Host, db.Port, db.Schema) + flags
}
