package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Episodes returns a Gin handler function that executes a episode indexation in background.
//
// @Summary Execute episode indexation process
// @Description Execute episode indexation process
// @Tags indexation-full
// @ID index-episode
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/episodes [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Episodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Episodes")
			err := handler.indexerApi.Episodes(ctx)
			if err != nil {
				return
			}
		}()
	}
}
