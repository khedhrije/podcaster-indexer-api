package api

import (
	"context"
)

func (api indexerApi) All(ctx context.Context) error {

	go func() {
		api.Walls(ctx)
	}()

	go func() {
		api.Blocks(ctx)
	}()

	go func() {
		api.Programs(ctx)
	}()

	go func() {
		api.Episodes(ctx)
	}()

	go func() {
		api.Medias(ctx)
	}()

	go func() {
		api.Categories(ctx)
	}()

	go func() {
		api.Tags(ctx)
	}()

	return nil
}
