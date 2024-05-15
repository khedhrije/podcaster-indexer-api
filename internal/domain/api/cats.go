package api

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Categories handles the indexing of categories
func (api indexerApi) Categories(ctx context.Context) error {
	// Initialize indexName with current Unix timestamp
	indexName := fmt.Sprintf("cats-%v", time.Now().Unix())

	// Create a new index
	if err := api.indexer.CreateIndex(ctx, indexName, Cats); err != nil {
		return err
	}

	// Assign the in-progress alias to the new index
	if err := api.indexer.CreateAlias(ctx, indexName, "in-progress"); err != nil {
		return err
	}

	// Retrieve all categories
	cats, err := api.catAdapter.FindAll(ctx)
	if err != nil {
		return err
	}

	// Bulk index categories
	if err := api.indexer.RecordBulkItems(ctx, indexName, mapArrayToInterface(cats), 5, 5); err != nil {
		return err
	}

	// Get all indexes associated with the latest alias
	latestIndexesNames := api.indexer.IndexByAlias(ctx, "latest")

	// Find the latest index name for categories
	var catsLatestIndexName string
	for _, latestIndexName := range latestIndexesNames {
		if strings.Contains(latestIndexName, "cats") {
			catsLatestIndexName = latestIndexName
		}
	}

	// If there is no latest cats index
	if catsLatestIndexName == "" {
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

		// Find the previous index name for categories
		var catsPreviousIndexName string
		for _, iName := range previousIndexesNames {
			if strings.Contains(iName, "cats") {
				catsPreviousIndexName = iName
			}
		}

		if catsPreviousIndexName == "" {
			// Make the latest index the previous index
			if err := api.indexer.DeleteAlias(ctx, catsLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, catsLatestIndexName, "previous"); err != nil {
				return err
			}
		} else {
			// Delete the previous index
			if err := api.indexer.DeleteIndexes(ctx, []string{catsPreviousIndexName}); err != nil {
				return err
			}
			// Update aliases accordingly
			if err := api.indexer.DeleteAlias(ctx, catsLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, catsLatestIndexName, "previous"); err != nil {
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
