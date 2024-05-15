// Package model defines the data structures for the application domain.
package model

// Media represents a media entity in the system.
// Media is associated with an episode and contains a direct link to the media content.
type Media struct {
	ID         string // Unique identifier for the media
	DirectLink string // Direct link to the media content
	Kind       string // Type or category of the media (e.g., audio, video)
	EpisodeID  string // Unique identifier for the associated episode
}
