// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/model"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
)

// categoryAdapter is a struct that acts as an adapter for interacting with
// the category data in the MySQL database.
type categoryAdapter struct {
	client *client
}

// NewCategoryAdapter creates a new category adapter with the provided MySQL client.
// It returns an implementation of the port.Category interface.
func NewCategoryAdapter(client *client) port.Category {
	return &categoryAdapter{
		client: client,
	}
}

// FindAll retrieves all category records from the database.
// It takes a context and returns a slice of model.Category and an error if the operation fails.
func (adapter *categoryAdapter) FindAll(ctx context.Context) ([]*model.Category, error) {
	const query = `
        SELECT * FROM category;
    `
	var categoriesDB []*CategoryDB
	if err := adapter.client.db.SelectContext(ctx, &categoriesDB, query); err != nil {
		return nil, err
	}
	var categories []*model.Category
	for _, categoryDB := range categoriesDB {
		mappedCategory := categoryDB.ToDomainModel()
		categories = append(categories, &mappedCategory)
	}
	return categories, nil
}

// Find retrieves a category record from the database by its UUID.
// It takes a context and the category's UUID, and returns a model.Category and an error if the operation fails.
func (adapter *categoryAdapter) Find(ctx context.Context, categoryUUID string) (*model.Category, error) {
	const query = `
        SELECT * FROM category WHERE UUID = UUID_TO_BIN(?);
    `
	var categoryDB CategoryDB
	if err := adapter.client.db.GetContext(ctx, &categoryDB, query, categoryUUID); err != nil {
		return nil, err
	}
	result := categoryDB.ToDomainModel()
	return &result, nil
}

// CategoryDB is a struct representing the category database model.
type CategoryDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	ParentID    uuid.UUID      `db:"parentUUID"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a CategoryDB database model to a model.Category domain model.
// It returns the corresponding model.Category.
func (db *CategoryDB) ToDomainModel() model.Category {
	return model.Category{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
		Parent: &model.Category{
			ID: db.ParentID.String(),
		},
	}
}

// FromDomainModel converts a model.Category domain model to a CategoryDB database model.
// It sets the fields of the CategoryDB based on the given model.Category.
func (db *CategoryDB) FromDomainModel(domain model.Category) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
	db.ParentID = uuid.Nil
	if domain.Parent != nil && domain.Parent.ID != "" {
		db.ParentID = uuid.MustParse(domain.Parent.ID)
	}
}
