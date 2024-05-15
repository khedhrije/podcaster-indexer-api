package api

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Episodes handles the indexing of episodes
func (api indexerApi) Episodes(ctx context.Context) error {
	// Initialize indexName with current Unix timestamp
	indexName := fmt.Sprintf("episodes-%v", time.Now().Unix())

	// Create a new index
	if err := api.indexer.CreateIndex(ctx, indexName, Episodes); err != nil {
		return err
	}

	// Assign the in-progress alias to the new index
	if err := api.indexer.CreateAlias(ctx, indexName, "in-progress"); err != nil {
		return err
	}

	// Retrieve all episodes
	episodes, err := api.episodeAdapter.FindAll(ctx)
	if err != nil {
		return err
	}

	// Bulk index episodes
	if err := api.indexer.RecordBulkItems(ctx, indexName, mapArrayToInterface(episodes), 5, 5); err != nil {
		return err
	}

	// Get all indexes associated with the latest alias
	latestIndexesNames := api.indexer.IndexByAlias(ctx, "latest")

	// Find the latest index name for episodes
	var episodesLatestIndexName string
	for _, latestIndexName := range latestIndexesNames {
		if strings.Contains(latestIndexName, "episodes") {
			episodesLatestIndexName = latestIndexName
		}
	}

	// If there is no latest episodes index
	if episodesLatestIndexName == "" {
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

		// Find the previous index name for episodes
		var episodesPreviousIndexName string
		for _, iName := range previousIndexesNames {
			if strings.Contains(iName, "episodes") {
				episodesPreviousIndexName = iName
			}
		}

		if episodesPreviousIndexName == "" {
			// Make the latest index the previous index
			if err := api.indexer.DeleteAlias(ctx, episodesLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, episodesLatestIndexName, "previous"); err != nil {
				return err
			}
		} else {
			// Delete the previous index
			if err := api.indexer.DeleteIndexes(ctx, []string{episodesPreviousIndexName}); err != nil {
				return err
			}
			// Update aliases accordingly
			if err := api.indexer.DeleteAlias(ctx, episodesLatestIndexName, "latest"); err != nil {
				return err
			}
			if err := api.indexer.CreateAlias(ctx, episodesLatestIndexName, "previous"); err != nil {
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
