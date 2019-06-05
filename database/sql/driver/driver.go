package driver

import (
	"context"
	"errors"
	"reflect"
)

type Value interface{}

type NamedValue struct {
	Name string
	Ordinal int
	Value Value
}

type Driver interface {
	Open(name string) (Conn, error)
}
type DriverContext interface {
	OpenConnector(name string)  (Connector, error)
}
type Connector interface {
	Connect(context.Context) (Conn, error)
	Driver() Driver
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Execer interface {
	Exec(query string, args []Value) (Result, error)
}
type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args []NamedValue) (Result, error)
}
type Queryer interface {
	Query(query string, args []Value) (Rows, error)
}
type QueryerContext interface {
	QueryContext(ctx context.Context, query string, args []NamedValue) (Rows, error)
}

type Conn interface {
	Prepare(query string) (Stmt, error)
	Close() error
	Begin() (Tx, error)
}
type ConnPrepareContext interface {
	PrepareContext(ctx context.Context, query string) (Stmt, error)
}

type IsolationLevel int
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly bool
}

type ConnBeginTx interface {
	BeginTx(ctx context.Context, opts TxOptions) (Tx, error)
}

type SessionResetter interface {
	ResetSession(ctx context.Context) error
}

type Stmt interface {
	Close() error
	NumInput() int
	Exec(args []Value) (Result, error)
	Query(args []Value) (Rows,error)
}
type StmtExecContext interface {
	ExecContext(ctx context.Context, args []NamedValue) (Result, error)
}
type StmtQueryContext interface {
	QueryContext(ctx context.Context, args []NamedValue) (Rows, error)
}

type NamedValueChecker interface {
	checkNamedValue(*NamedValue) error
}

type ColumnConverter interface {
	ColumnConverter(idx int) ValueConverter
}

type Rows interface {
	Columns() []string
	Next(dest []Value) error
	Close() error
}

type RowsNextResultSet interface {
	Rows
	HasNextResultSet() bool
	NextResultSet() error
}

type RowsColumnTypeScanType interface {
	Rows
	ColumnTypeScanType(index int) reflect.Type
}

type RowsColumnTypeDatabaseTypeName interface {
	Rows
	ColumnTypeDatabaseTypeName(index int) string
}

type RowsColumnTypeLength interface {
	Rows
	ColumnTypeLength(index int) (length int64, ok bool)
}

type RowsColumnTypeNullable interface {
	Rows
	ColumnTypeNullable(index int) (nullable, ok bool)
}

type RowsColumnTypePrecisionScale interface {
	Rows
	ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool)
}

type Tx interface {
	Commit() error
	Rollback() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type RowsAffected int64
var _ Result = RowsAffected(0)
func (RowsAffected) LastInsertId() (int64, error) {
	return 0, errors.New("LastInsertId is not supported by this driver")
}
func (v RowsAffected) RowsAffected() (int64, error) {
	return int64(v), nil
}

type noRows struct{}
var _ Result = noRows{}
func (noRows) LastInsertId() (int64, error) {
	return 0, errors.New("no LastInsertId available after DDL statement")
}
func (noRows) RowsAffected() (int64, error) {
	return 0, errors.New("no RowsAffected available after DDL statement")
}




