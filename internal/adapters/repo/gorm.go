package repo

import (
	"context"
	"database/sql"
)

type GSQL interface {
	AutoMigrate(models ...interface{})
	Create(ctx context.Context, data interface{}) error
	Update(ctx context.Context, data interface{}) error
	RawQuery(ctx context.Context, scanner interface{}, query string, args ...interface{}) error
	UpdateMany(ctx context.Context, data, query interface{}, args ...interface{}) error
	UpdateOne(ctx context.Context, data, query interface{}, args ...interface{}) error
	BeginFind(ctx context.Context, value interface{}) Find
	Delete(ctx context.Context, data interface{}, condition interface{}, args ...interface{}) error
	DeleteByQuery(ctx context.Context, data, query interface{}, args ...interface{}) error
	GetInstance() interface{}
}
type Find interface {
	Where(query interface{}, args ...interface{}) Find
	Having(query interface{}, args ...interface{}) Find
	Page(current int, limit int) Find
	Join(query string, args ...interface{}) Find
	Or(query interface{}, args ...interface{}) Find
	Not(query interface{}, args ...interface{}) Find
	Count(total *int64) error
	Find(result interface{}, args ...interface{}) error
	First(result interface{}, args ...interface{}) error
	Select(query interface{}, args ...interface{}) Find
	Scan(result interface{}) error
	OrderBy(query string) Find
	Group(query string) Find
	Limit(limit int) Find
	Rows() (*sql.Rows, error)
}
type FindAllInput struct {
	PageInput PageInput
	JoinInput []JoinInput
	OrderBy   string
}
type PageInput struct {
	Current int
	Limit   int
}
type JoinInput struct {
	Table     string
	Condition string
	JoinType  string
}
type GTransaction interface {
	Begin(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
type Tx interface {
	WithTransaction(context.Context, func(context.Context) error) error
}
