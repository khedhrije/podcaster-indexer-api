package api

import (
	"context"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// Indexer interface defines a method for handling categories
type Indexer interface {
	Categories(ctx context.Context) error
	Tags(ctx context.Context) error
	Walls(ctx context.Context) error
	Blocks(ctx context.Context) error
	Programs(ctx context.Context) error
	Episodes(ctx context.Context) error
	Medias(ctx context.Context) error
	All(ctx context.Context) error
}

// indexerApi struct implements the Indexer interface
type indexerApi struct {
	indexer        port.Indexer
	wallAdapter    port.Wall
	catAdapter     port.Category
	tagAdapter     port.Tag
	blockAdapter   port.Block
	programAdapter port.Program
	episodeAdapter port.Episode
	mediaAdapter   port.Media
}

// NewIndexerApi returns a new instance of indexerApi
func NewIndexerApi(
	indexer port.Indexer,
	wallAdapter port.Wall,
	catAdapter port.Category,
	tagAdapter port.Tag,
	blockAdapter port.Block,
	programAdapter port.Program,
	episodeAdapter port.Episode,
	mediaAdapter port.Media) Indexer {
	return &indexerApi{
		indexer:        indexer,
		catAdapter:     catAdapter,
		tagAdapter:     tagAdapter,
		wallAdapter:    wallAdapter,
		blockAdapter:   blockAdapter,
		programAdapter: programAdapter,
		episodeAdapter: episodeAdapter,
		mediaAdapter:   mediaAdapter,
	}
}

// mapArrayToInterface converts an array of any type to an array of empty interfaces
func mapArrayToInterface[V any](array []V) []any {
	if len(array) == 1 {
		return []any{array[0]}
	}

	var res []interface{}
	for _, el := range array {
		res = append(res, el)
	}
	return res
}

const (
	Cats     = "cats"
	Tags     = "tags"
	Walls    = "walls"
	Blocks   = "blocks"
	Programs = "programs"
	Episodes = "episodes"
	Medias   = "medias"
)
