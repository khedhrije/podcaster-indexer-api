// Package model defines the data structures for the application domain.
package model

// Tag represents a tag entity in the system.
// A Tag is used to categorize or label programs.
type Tag struct {
	ID          string // Unique identifier for the tag
	Name        string // Name of the tag
	Description string // Description of the tag
}
