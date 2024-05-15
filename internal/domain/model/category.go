// Package model defines the data structures for the application domain.
package model

// Category represents a category entity in the system.
// A Category can have a hierarchical relationship with other categories, allowing for nested structures.
type Category struct {
	ID          string      // Unique identifier for the category
	Name        string      // Name of the category
	Description string      // Description of the category
	Parent      *Category   // Reference to the parent category, if any
	Children    []*Category // List of child categories
}
