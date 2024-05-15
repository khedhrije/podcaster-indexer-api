// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// wallAdapter is a struct that acts as an adapter for interacting with
// the wall data in the MySQL database.
type wallAdapter struct {
	client *client
}

// NewWallAdapter creates a new wall adapter with the provided MySQL client.
// It returns an instance of wallAdapter.
func NewWallAdapter(client *client) port.Wall {
	return &wallAdapter{
		client: client,
	}
}

// FindAll retrieves all wall records from the database.
// It takes a context and returns a slice of model.Wall and an error if the operation fails.
func (adapter *wallAdapter) FindAll(ctx context.Context) ([]*model.Wall, error) {
	const query = `
        SELECT * FROM wall
    `
	var wallsDB []*WallDB
	if err := adapter.client.db.SelectContext(ctx, &wallsDB, query); err != nil {
		return nil, err
	}
	var walls []*model.Wall
	for _, wallDB := range wallsDB {
		mappedWall := wallDB.ToDomainModel()
		walls = append(walls, &mappedWall)
	}
	return walls, nil
}

// WallDB is a struct representing the wall database model.
type WallDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a WallDB database model to a model.Wall domain model.
// It returns the corresponding model.Wall.
func (db *WallDB) ToDomainModel() model.Wall {
	return model.Wall{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Wall domain model to a WallDB database model.
// It sets the fields of the WallDB based on the given model.Wall.
func (db *WallDB) FromDomainModel(domain model.Wall) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
