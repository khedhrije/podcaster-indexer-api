package port

import (
	"context"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
)

type Wall interface {
	FindAll(ctx context.Context) ([]*model.Wall, error)
}

type Category interface {
	FindAll(ctx context.Context) ([]*model.Category, error)
}
type Tag interface {
	FindAll(ctx context.Context) ([]*model.Tag, error)
}

type Block interface {
	FindAll(ctx context.Context) ([]*model.Block, error)
}

type Program interface {
	FindAll(ctx context.Context) ([]*model.Program, error)
}

type Episode interface {
	FindAll(ctx context.Context) ([]*model.Episode, error)
}

type Media interface {
	FindAll(ctx context.Context) ([]*model.Media, error)
}
