package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// All returns a Gin handler function that executes all indexations in background.
//
// @Summary Execute all indexation processes
// @Description Execute all indexation processes
// @Tags indexation-full
// @ID index-all
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/all [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) All() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-All")
			err := handler.indexerApi.All(ctx)
			if err != nil {
				return
			}
		}()
	}
}
