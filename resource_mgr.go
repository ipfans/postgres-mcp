package postgresmcp

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

// ResourceManager 处理资源相关操作
type ResourceManager struct {
	db      DatabaseQuerier
	baseURL string
}

func NewResourceManager(db DatabaseQuerier, baseURL string) *ResourceManager {
	return &ResourceManager{
		db:      db,
		baseURL: baseURL,
	}
}

func (rm *ResourceManager) ListResources(ctx context.Context) ([]*mcp_golang.Content, error) {
	rows, err := rm.db.Query(ctx, `
		SELECT 
			t.table_name,
			array_agg(
				c.column_name || ' ' || c.data_type
				ORDER BY c.ordinal_position
			) as columns
		FROM information_schema.tables t
		JOIN information_schema.columns c 
			ON c.table_name = t.table_name 
			AND c.table_schema = t.table_schema
		WHERE t.table_schema = 'public'
		GROUP BY t.table_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*mcp_golang.Content
	for rows.Next() {
		var tableName string
		var columns []string
		if err := rows.Scan(&tableName, &columns); err != nil {
			return nil, err
		}

		schemaInfo := fmt.Sprintf("表名: %s\n字段:\n", tableName)
		for _, col := range columns {
			schemaInfo += fmt.Sprintf("- %s\n", col)
		}
		content := mcp_golang.NewTextContent(schemaInfo)
		resources = append(resources, content)
	}
	return resources, rows.Err()
}
