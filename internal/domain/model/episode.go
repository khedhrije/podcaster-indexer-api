// Package model defines the data structures for the application domain.
package model

// Episode represents an episode entity in the system.
// An Episode is a part of a Program and contains media content.
type Episode struct {
	ID          string // Unique identifier for the episode
	Name        string // Name of the episode
	Description string // Description of the episode
	Position    int    // Position of the episode within its program
	Media       Media  // Media content associated with the episode
	ProgramID   string // Unique identifier for the associated program
}
