package sql_builder

import (
	"context"
	"database/sql"
)

type Sqlable interface {
	ExecContext(ctx context.Context, query string, args ...interface{})(sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type DbMock struct{
	query     string
	args      []interface{}
	affected  int64
	lastId    int64
	lastError error
}

func NewDbMock() *DbMock {
	return &DbMock{}
}

func (dm *DbMock) RowsAffected() (int64, error){ return dm.affected, dm.lastError }
func (dm *DbMock) LastInsertId() (int64, error){ return dm.lastId, dm.lastError }

func (dm *DbMock) ExecContext(ctx context.Context, query string, args ...interface{})(sql.Result, error){
	dm.query = query
	dm.args  = args
	dm.lastId = 1
	dm.affected = int64(len(args))
	dm.lastError = nil

	return dm,nil
}

func (dm *DbMock) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	dm.query = query
	dm.args  = args
	dm.lastId = 1
	dm.affected = int64(len(args))
	dm.lastError = nil

	return nil
}
