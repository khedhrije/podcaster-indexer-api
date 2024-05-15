package api

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Walls handles the indexing of walls
func (api indexerApi) Walls(ctx context.Context) error {
	// Initialize indexName with current Unix timestamp
	indexName := fmt.Sprintf("walls-%v", time.Now().Unix())

	// Create a new index
	if err := api.indexer.CreateIndex(ctx, indexName, Walls); err != nil {
		return err
	}

	// Assign the in-progress alias to the new index
	if err := api.indexer.CreateAlias(ctx, indexName, "in-progress"); err != nil {
		return err
	}

	// Retrieve all walls
	walls, err := api.wallAdapter.FindAll(ctx)
	if err != nil {
		return err
	}

	// Bulk index walls
	if err := api.indexer.RecordBulkItems(ctx, indexName, mapArrayToInterface(walls), 5, 5); err != nil {
		return err
	}

	// Get all indexes associated with the latest alias
	latestIndexesNames := api.indexer.IndexByAlias(ctx, "latest")

	// Find the latest index name for walls
	var wallsLatestIndexName string
	for _, latestIndexName := range latestIndexesNames {
		if strings.Contains(latestIndexName, "walls") {
			wallsLatestIndexName = latestIndexName
		}
	}

	// If there is no latest walls index
	if wallsLatestIndexName == "" {
		// Make the in-progress index the latest index
		if err := api.indexer.DeleteAlias(ctx, indexName, "in-progress"); err != nil {
			return err
		}
		if err := api.indexer.CreateAlias(ctx, indexName, "latest"); err != nil {
			return err
		}
	} else {
		// Handle existing latest index
		previousIndexesNames := api.indexer.IndexByAlias(ctx, "previous")

		// Find the previous index name for walls
		var wallsPreviousIndexName string
		for _, iName := range previousIndexesNames {
			if strings.Contains(iName, "walls") {
				wallsPreviousIndexName = iName
			}
		}

		if wallsPreviousIndexName == "" {
			// Make the latest index the previous index
			if err := api.indexer.DeleteAlias(ctx, wallsLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, wallsLatestIndexName, "previous"); err != nil {
				return err
			}
		} else {
			// Delete the previous index
			if err := api.indexer.DeleteIndexes(ctx, []string{wallsPreviousIndexName}); err != nil {
				return err
			}
			// Update aliases accordingly
			if err := api.indexer.DeleteAlias(ctx, wallsLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, wallsLatestIndexName, "previous"); err != nil {
				return err
			}
		}

		// Always remove the in-progress alias from the new index
		if err := api.indexer.DeleteAlias(ctx, indexName, "in-progress"); err != nil {
			return err
		}

		// Create latest alias for the new index
		if err := api.indexer.CreateAlias(ctx, indexName, "latest"); err != nil {
			return err
		}
	}

	return nil
}
