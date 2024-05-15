package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Walls returns a Gin handler function that executes a wall indexation in background.
//
// @Summary Execute wall indexation process
// @Description Execute wall indexation process
// @Tags indexation-full
// @ID index-wall
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/walls [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Walls() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Walls")
			err := handler.indexerApi.Walls(ctx)
			if err != nil {
				return
			}
		}()
	}
}
