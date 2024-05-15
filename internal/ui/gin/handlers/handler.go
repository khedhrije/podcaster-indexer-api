// Package handlers provides HTTP request handlers for managing categories.
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/api"
)

// Indexation represents the interface for managing indexation.
type Indexation interface {
	Categories() gin.HandlerFunc
	Tags() gin.HandlerFunc
	Walls() gin.HandlerFunc
	Blocks() gin.HandlerFunc
	Programs() gin.HandlerFunc
	Episodes() gin.HandlerFunc
	Medias() gin.HandlerFunc
	All() gin.HandlerFunc
}

// indexationHandler is an implementation of the Indexation interface.
type indexationHandler struct {
	indexerApi api.Indexer
}

// NewIndexationHandler creates a new instance of Indexation interface.
func NewIndexationHandler(indexerApi api.Indexer) Indexation {
	return &indexationHandler{
		indexerApi: indexerApi,
	}
}
