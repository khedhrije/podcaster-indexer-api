package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Programs returns a Gin handler function that executes a program indexation in background.
//
// @Summary Execute program indexation process
// @Description Execute program indexation process
// @Tags indexation-full
// @ID index-program
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/indexation/programs [post]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler indexationHandler) Programs() gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			ctx := context.WithValue(c, "goroutine", "Indexation-Programs")
			err := handler.indexerApi.Programs(ctx)
			if err != nil {
				return
			}
		}()
	}
}
