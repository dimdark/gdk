package sql

import (
	"github.com/dimdark/gdk/database/sql/driver"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	driversMu sync.RWMutex
	drivers = make(map[string]driver.Driver)
)

var nowFunc = time.Now()

func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	drivers = make(map[string]driver.Driver)
}

func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

type NamedArg struct {
	_Named_Fields_Required struct{}
	Name string
	Value interface{}
}
func Named(name string, value interface{}) NamedArg {
	return NamedArg{Name: name, Value: value}
}

type IsolationLevel int
const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)
func (i IsolationLevel) String() string {
	switch i {
	case LevelDefault:
		return "Default"
	case LevelReadUncommitted:
		return "Read Uncommitted"
	case LevelReadCommitted:
		return "Read Committed"
	case LevelWriteCommitted:
		return "Write Committed"
	case LevelRepeatableRead:
		return "Repeatable Read"
	case LevelSnapshot:
		return "Snapshot"
	case LevelSerializable:
		return "Serializable"
	case LevelLinearizable:
		return "Lineariable"
	default:
		return "IsolationLevel(" + strconv.Itoa(int(i)) + ")"
	}
}

type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly bool
}

type RawBytes []byte

type Scanner interface {
	Scan(src interface{}) error
}

type NullString struct {
	String string
	Valid bool
}
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return convertAssign(&ns.String, value)
}

func (ns NullString) Value() (driver.Value, error){
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

type NullInt64 struct {
	Int64 int64
	Valid bool
}
func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Int64, value)
}
func (n *NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

type NullFloat64 struct {
	Float64 float64
	Valid bool
}
func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Float64, value)
}
func (n *NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

type NullBool struct {
	Bool bool
	Valid bool
}
func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Bool, value)
}
func (n *NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

type DB struct {
	waitDuration int64
	connector driver.Connector
	numClosed uint64

	mu sync.Mutex
	freeConn []*driverConn
}












