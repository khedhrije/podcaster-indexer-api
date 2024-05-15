package port

import "context"

// Indexer is an interface that abstracts the functionalities of the Elasticsearch client.
type Indexer interface {
	CreateIndex(ctx context.Context, indexName string, docsType string) error
	DeleteIndexes(ctx context.Context, indexNames []string) error
	CreateAlias(ctx context.Context, indexName, aliasName string) error
	DeleteAlias(ctx context.Context, indexName, aliasName string) error
	IndexByAlias(ctx context.Context, aliasName string) []string
	MoveIndex(ctx context.Context, indexationName string) error
	RecordBulkItems(ctx context.Context, indexName string, items []interface{}, backoffRetryCount, backoffTimeSeconds int) error
}
