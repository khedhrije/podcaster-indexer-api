// Package model defines the data structures for the application domain.
package model

// Program represents a program entity in the system.
// A Program is a collection of episodes and contains metadata about the program.
type Program struct {
	ID          string    // Unique identifier for the program
	Name        string    // Name of the program
	Description string    // Description of the program
	Episodes    []Episode // List of episodes associated with the program
}
