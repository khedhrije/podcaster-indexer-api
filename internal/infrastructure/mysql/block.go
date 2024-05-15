// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// blockAdapter is a struct that acts as an adapter for interacting with
// the block data in the MySQL database.
type blockAdapter struct {
	client *client
}

// NewBlockAdapter creates a new block adapter with the provided MySQL client.
// It returns an implementation of the BlockPersister interface.
func NewBlockAdapter(client *client) port.Block {
	return &blockAdapter{
		client: client,
	}
}

// FindAll retrieves all block records from the database.
// It takes a context and returns a slice of model.Block and an error if the operation fails.
func (adapter *blockAdapter) FindAll(ctx context.Context) ([]*model.Block, error) {
	const query = `
        SELECT * FROM block;
    `
	var blocksDB []*BlockDB
	if err := adapter.client.db.SelectContext(ctx, &blocksDB, query); err != nil {
		return nil, err
	}
	var blocks []*model.Block
	for _, blockDB := range blocksDB {
		mappedBlock := blockDB.ToDomainModel()
		blocks = append(blocks, &mappedBlock)
	}
	return blocks, nil
}

// BlockDB is a struct representing the block database model.
type BlockDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Kind        sql.NullString `db:"kind"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a BlockDB database model to a model.Block domain model.
// It returns the corresponding model.Block.
func (db *BlockDB) ToDomainModel() model.Block {
	return model.Block{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
		Kind:        db.Kind.String,
	}
}

// FromDomainModel converts a model.Block domain model to a BlockDB database model.
// It sets the fields of the BlockDB based on the given model.Block.
func (db *BlockDB) FromDomainModel(domain model.Block) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
	db.Kind = sql.NullString{String: domain.Kind, Valid: domain.Kind != ""}
}
