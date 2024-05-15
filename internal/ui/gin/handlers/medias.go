package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Medias returns a Gin handler function that executes a media indexation in background.
//
// @Summary Execute media indexation process
// @Description Execute media indexation process
// @Tags indexation-full
// @ID index-media
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/medias [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Medias() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Medias")
			err := handler.indexerApi.Medias(ctx)
			if err != nil {
				return
			}
		}()
	}
}
