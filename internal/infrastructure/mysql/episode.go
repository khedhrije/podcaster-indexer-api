// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// episodeAdapter is a struct that acts as an adapter for interacting with
// the episode data in the MySQL database.
type episodeAdapter struct {
	client *client
}

// NewEpisodeAdapter creates a new episode adapter with the provided MySQL client.
// It returns an implementation of the EpisodePersister interface.
func NewEpisodeAdapter(client *client) port.Episode {
	return &episodeAdapter{
		client: client,
	}
}

// FindAll retrieves all episode records from the database.
// It takes a context and returns a slice of model.Episode and an error if the operation fails.
func (adapter *episodeAdapter) FindAll(ctx context.Context) ([]*model.Episode, error) {
	const query = `
        SELECT * FROM episode;
    `
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query); err != nil {
		return nil, err
	}
	var episodes []*model.Episode
	for _, episodeDB := range episodesDB {
		mappedEpisode := episodeDB.ToDomainModel()
		episodes = append(episodes, &mappedEpisode)
	}
	return episodes, nil
}

// EpisodeDB is a struct representing the episode database model.
type EpisodeDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Position    int            `db:"position"`
	ProgramID   uuid.UUID      `db:"programUUID"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts an EpisodeDB database model to a model.Episode domain model.
// It returns the corresponding model.Episode.
func (db *EpisodeDB) ToDomainModel() model.Episode {
	return model.Episode{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
		Position:    db.Position,
		ProgramID:   db.ProgramID.String(),
	}
}

// FromDomainModel converts a model.Episode domain model to an EpisodeDB database model.
// It sets the fields of the EpisodeDB based on the given model.Episode.
func (db *EpisodeDB) FromDomainModel(domain model.Episode) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
	db.Position = domain.Position
	db.ProgramID = uuid.Nil
	if domain.ProgramID != "" {
		db.ProgramID = uuid.MustParse(domain.ProgramID)
	}
}
