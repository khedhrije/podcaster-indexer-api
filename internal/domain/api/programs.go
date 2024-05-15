package api

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Programs handles the indexing of programs
func (api indexerApi) Programs(ctx context.Context) error {
	// Initialize indexName with current Unix timestamp
	indexName := fmt.Sprintf("programs-%v", time.Now().Unix())

	// Create a new index
	if err := api.indexer.CreateIndex(ctx, indexName, Programs); err != nil {
		return err
	}

	// Assign the in-progress alias to the new index
	if err := api.indexer.CreateAlias(ctx, indexName, "in-progress"); err != nil {
		return err
	}

	// Retrieve all programs
	programs, err := api.programAdapter.FindAll(ctx)
	if err != nil {
		return err
	}

	// Bulk index programs
	if err := api.indexer.RecordBulkItems(ctx, indexName, mapArrayToInterface(programs), 5, 5); err != nil {
		return err
	}

	// Get all indexes associated with the latest alias
	latestIndexesNames := api.indexer.IndexByAlias(ctx, "latest")

	// Find the latest index name for programs
	var programsLatestIndexName string
	for _, latestIndexName := range latestIndexesNames {
		if strings.Contains(latestIndexName, "programs") {
			programsLatestIndexName = latestIndexName
		}
	}

	// If there is no latest programs index
	if programsLatestIndexName == "" {
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

		// Find the previous index name for programs
		var programsPreviousIndexName string
		for _, iName := range previousIndexesNames {
			if strings.Contains(iName, "programs") {
				programsPreviousIndexName = iName
			}
		}

		if programsPreviousIndexName == "" {
			// Make the latest index the previous index
			if err := api.indexer.DeleteAlias(ctx, programsLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, programsLatestIndexName, "previous"); err != nil {
				return err
			}
		} else {
			// Delete the previous index
			if err := api.indexer.DeleteIndexes(ctx, []string{programsPreviousIndexName}); err != nil {
				return err
			}
			// Update aliases accordingly
			if err := api.indexer.DeleteAlias(ctx, programsLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, programsLatestIndexName, "previous"); err != nil {
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
