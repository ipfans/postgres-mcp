package postgresmcp

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

// MCPServer 处理MCP服务器配置和启动
type MCPServer struct {
	server    *mcp_golang.Server
	transport *http.GinTransport
	resources *ResourceManager
	queryExec *QueryExecutor
}

type QueryArguments struct {
	Query string `json:"query" jsonschema:"required,description=The sql query to execute"`
}

func NewMCPServer(dbString string, baseURL string) (*MCPServer, error) {
	dbpool, err := pgxpool.New(context.Background(), dbString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	transport := http.NewGinTransport()
	server := mcp_golang.NewServer(transport)
	resources := NewResourceManager(dbpool, baseURL)
	queryExec := NewQueryExecutor(dbpool)

	return &MCPServer{
		server:    server,
		transport: transport,
		resources: resources,
		queryExec: queryExec,
	}, nil
}

func (s *MCPServer) registerHandlers() error {
	// 注册资源处理器
	if err := s.server.RegisterResource(
		"resources",
		"Database Resources",
		"List all database tables",
		"application/json",
		s.resources.ListResources,
	); err != nil {
		return err
	}

	// 注册查询工具
	if err := s.server.RegisterTool(
		"query",
		"Run a read-only SQL query",
		func(ctx context.Context, req QueryArguments) (*mcp_golang.ToolResponse, error) {
			return s.queryExec.ExecuteReadOnlyQuery(ctx, req.Query)
		},
	); err != nil {
		return err
	}

	return nil
}

func (s *MCPServer) Start() error {
	if err := s.registerHandlers(); err != nil {
		return err
	}

	router := gin.Default()
	router.POST("/mcp", s.transport.Handler())
	return router.Run(":8080")
}

// Server 是主入口函数
func Server(dbString string) {
	baseURL := "http://localhost:8080/resources"
	server, err := NewMCPServer(dbString, baseURL)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
		return
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
