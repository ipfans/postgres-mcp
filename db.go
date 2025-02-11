package postgresmcp

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DatabaseQuerier 定义数据库查询接口
type DatabaseQuerier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}
