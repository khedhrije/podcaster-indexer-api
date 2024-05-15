package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Blocks returns a Gin handler function that executes a block indexation in background.
//
// @Summary Execute block indexation process
// @Description Execute block indexation process
// @Tags indexation-full
// @ID index-block
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/blocks [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Blocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Blocks")
			err := handler.indexerApi.Blocks(ctx)
			if err != nil {
				return
			}
		}()
	}
}
