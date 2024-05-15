// Package model defines the data structures for the application domain.
package model

// Wall represents a wall entity in the system.
// A Wall is a collection of blocks, each containing content and organizational metadata.
type Wall struct {
	ID          string  // Unique identifier for the wall
	Name        string  // Name of the wall
	Description string  // Description of the wall
	Blocks      []Block // List of blocks associated with the wall
}
