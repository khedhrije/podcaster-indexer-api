package api

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Blocks handles the indexing of blocks
func (api indexerApi) Blocks(ctx context.Context) error {
	// Initialize indexName with current Unix timestamp
	indexName := fmt.Sprintf("blocks-%v", time.Now().Unix())

	// Create a new index
	if err := api.indexer.CreateIndex(ctx, indexName, Blocks); err != nil {
		return err
	}

	// Assign the in-progress alias to the new index
	if err := api.indexer.CreateAlias(ctx, indexName, "in-progress"); err != nil {
		return err
	}

	// Retrieve all blocks
	blocks, err := api.blockAdapter.FindAll(ctx)
	if err != nil {
		return err
	}

	// Bulk index blocks
	if err := api.indexer.RecordBulkItems(ctx, indexName, mapArrayToInterface(blocks), 5, 5); err != nil {
		return err
	}

	// Get all indexes associated with the latest alias
	latestIndexesNames := api.indexer.IndexByAlias(ctx, "latest")

	// Find the latest index name for blocks
	var blocksLatestIndexName string
	for _, latestIndexName := range latestIndexesNames {
		if strings.Contains(latestIndexName, "blocks") {
			blocksLatestIndexName = latestIndexName
		}
	}

	// If there is no latest blocks index
	if blocksLatestIndexName == "" {
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

		// Find the previous index name for blocks
		var blocksPreviousIndexName string
		for _, iName := range previousIndexesNames {
			if strings.Contains(iName, "blocks") {
				blocksPreviousIndexName = iName
			}
		}

		if blocksPreviousIndexName == "" {
			// Make the latest index the previous index
			if err := api.indexer.DeleteAlias(ctx, blocksLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, blocksLatestIndexName, "previous"); err != nil {
				return err
			}
		} else {
			// Delete the previous index
			if err := api.indexer.DeleteIndexes(ctx, []string{blocksPreviousIndexName}); err != nil {
				return err
			}
			// Update aliases accordingly
			if err := api.indexer.DeleteAlias(ctx, blocksLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, blocksLatestIndexName, "previous"); err != nil {
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
