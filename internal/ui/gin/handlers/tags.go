package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Tags returns a Gin handler function that executes a tag indexation in background.
//
// @Summary Execute tag indexation process
// @Description Execute tag indexation process
// @Tags indexation-full
// @ID index-tag
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/tags [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Tags() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Tags")
			err := handler.indexerApi.Tags(ctx)
			if err != nil {
				return
			}
		}()
	}
}
