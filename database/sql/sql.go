package sql

import (
	"database/sql/driver"
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



















