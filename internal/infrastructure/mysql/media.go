// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// mediaAdapter is a struct that acts as an adapter for interacting with
// the media data in the MySQL database.
type mediaAdapter struct {
	client *client
}

// NewMediaAdapter creates a new media adapter with the provided MySQL client.
// It returns an implementation of the MediaPersister interface.
func NewMediaAdapter(client *client) port.Media {
	return &mediaAdapter{
		client: client,
	}
}

// FindAll retrieves all media records from the database.
// It takes a context and returns a slice of model.Media and an error if the operation fails.
func (adapter *mediaAdapter) FindAll(ctx context.Context) ([]*model.Media, error) {
	const query = `
        SELECT * FROM media;
    `
	var mediaDB []*MediaDB
	if err := adapter.client.db.SelectContext(ctx, &mediaDB, query); err != nil {
		return nil, err
	}
	var media []*model.Media
	for _, mediaEntry := range mediaDB {
		mappedMedia := mediaEntry.ToDomainModel()
		media = append(media, &mappedMedia)
	}
	return media, nil
}

// MediaDB is a struct representing the media database model.
type MediaDB struct {
	UUID       uuid.UUID      `db:"UUID"`
	DirectLink sql.NullString `db:"direct_link"`
	Kind       sql.NullString `db:"kind"`
	EpisodeID  uuid.UUID      `db:"episodeUUID"`
}

// ToDomainModel converts a MediaDB database model to a model.Media domain model.
// It returns the corresponding model.Media.
func (db *MediaDB) ToDomainModel() model.Media {
	return model.Media{
		ID:         db.UUID.String(),
		DirectLink: db.DirectLink.String,
		Kind:       db.Kind.String,
		EpisodeID:  db.EpisodeID.String(),
	}
}

// FromDomainModel converts a model.Media domain model to a MediaDB database model.
// It sets the fields of the MediaDB based on the given model.Media.
func (db *MediaDB) FromDomainModel(domain model.Media) {
	db.UUID = uuid.MustParse(domain.ID)
	db.DirectLink = sql.NullString{String: domain.DirectLink, Valid: domain.DirectLink != ""}
	db.Kind = sql.NullString{String: domain.Kind, Valid: domain.Kind != ""}
	db.EpisodeID = uuid.Nil
	if domain.EpisodeID != "" {
		db.EpisodeID = uuid.MustParse(domain.EpisodeID)
	}
}
