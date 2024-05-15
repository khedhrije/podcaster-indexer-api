// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// tagAdapter is a struct that acts as an adapter for interacting with
// the tag data in the MySQL database.
type tagAdapter struct {
	client *client
}

// NewTagAdapter creates a new tag adapter with the provided MySQL client.
// It returns an instance of tagAdapter.
func NewTagAdapter(client *client) port.Tag {
	return &tagAdapter{
		client: client,
	}
}

// FindAll retrieves all tag records from the database.
// It takes a context and returns a slice of model.Tag and an error if the operation fails.
func (adapter *tagAdapter) FindAll(ctx context.Context) ([]*model.Tag, error) {
	const query = `
        SELECT * FROM tag
    `
	var tagsDB []*TagDB
	if err := adapter.client.db.SelectContext(ctx, &tagsDB, query); err != nil {
		return nil, err
	}
	var tags []*model.Tag
	for _, tagDB := range tagsDB {
		mappedTag := tagDB.ToDomainModel()
		tags = append(tags, &mappedTag)
	}
	return tags, nil
}

// TagDB is a struct representing the tag database model.
type TagDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a TagDB database model to a model.Tag domain model.
// It returns the corresponding model.Tag.
func (db *TagDB) ToDomainModel() model.Tag {
	return model.Tag{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Tag domain model to a TagDB database model.
// It sets the fields of the TagDB based on the given model.Tag.
func (db *TagDB) FromDomainModel(domain model.Tag) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
