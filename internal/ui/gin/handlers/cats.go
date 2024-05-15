package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Categories returns a Gin handler function that executes a category indexation in background.
//
// @Summary Execute category indexation process
// @Description Execute category indexation process
// @Tags indexation-full
// @ID index-category
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/categories [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Categories() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Categories")
			err := handler.indexerApi.Categories(ctx)
			if err != nil {
				return
			}
		}()
	}
}
