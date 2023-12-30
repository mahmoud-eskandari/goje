package goje

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DefatultDB *sql.DB

// DefaultHandler make a handler from default database
func DefaultHandler(ctx context.Context) *ContextHandler {
	return &ContextHandler{
		Ctx: ctx,
		DB:  DefatultDB,
	}
}

// ConnectDB Connect and set default database
func ConnectDB(conn *DBConfig) error {
	if conn.Driver != "mysql" {
		return ErrUnknownDBDriver
	}

	db, err := sql.Open(conn.Driver, conn.String())

	if err != nil {
		return err
	}

	db.SetConnMaxIdleTime(conn.MaxIdleTime)
	db.SetMaxIdleConns(conn.MaxIdleConns)
	db.SetMaxOpenConns(conn.MaxOpenConns)
	db.SetConnMaxLifetime(conn.ConnMaxLifetime)

	DefatultDB = db
	return nil
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True"+flags, db.User, db.Password, db.Host, db.Port, db.Schema)
}
