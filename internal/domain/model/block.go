// Package model defines the data structures for the application domain.
package model

// Block represents a block entity in the system.
// A Block is a logical grouping that can contain multiple Programs.
type Block struct {
	ID          string    // Unique identifier for the block
	Name        string    // Name of the block
	Description string    // Description of the block
	Kind        string    // Type or category of the block
	Programs    []Program // List of programs associated with the block
}
