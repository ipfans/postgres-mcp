package postgresmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	mcp_golang "github.com/metoro-io/mcp-golang"
)

// QueryExecutor 处理数据库查询操作
type QueryExecutor struct {
	db DatabaseQuerier
}

func NewQueryExecutor(db DatabaseQuerier) *QueryExecutor {
	return &QueryExecutor{db: db}
}

func (qe *QueryExecutor) ExecuteReadOnlyQuery(ctx context.Context, query string) (*mcp_golang.ToolResponse, error) {
	tx, err := qe.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("Could not rollback transaction: %v", err)
		}
	}()

	if _, err := tx.Exec(ctx, "BEGIN TRANSACTION READ ONLY"); err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	columnNames := rows.FieldDescriptions()
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columnNames {
			rowMap[col.Name] = values[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	jsonResult, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query results to JSON: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	content := mcp_golang.NewTextContent(string(jsonResult))
	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{content},
	}, nil
}
