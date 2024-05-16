package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	spec "github.com/khedhrije/podcaster-indexer-api/deployments/swagger"
	"github.com/khedhrije/podcaster-indexer-api/internal/configuration"
	"github.com/khedhrije/podcaster-indexer-api/internal/ui/gin/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
)

// CreateRouter sets up and returns a new Gin router with the defined routes.
func CreateRouter(handler handlers.Indexation) *gin.Engine {
	// Initialize a new Gin router without any middleware by default.
	r := gin.New()
	r.Use(cors.Default())
	r.GET("/health", health())

	// Configure Swagger documentation URL based on the environment.
	if configuration.Config.Env == "dev" {
		spec.SwaggerInfo.Host = configuration.Config.DocsAddress
	} else {
		spec.SwaggerInfo.Host = removeHTTPS(configuration.Config.DocsAddress)
	}

	// Set up the route for Swagger documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Define private routes that require authentication.
	private := r.Group("/private")
	//private.Use(TokenValidatorMiddleware())
	{
		// Routes for managing walls.
		walls := private.Group("/indexation")
		{
			walls.POST("/categories", handler.Categories())
			walls.POST("/tags", handler.Tags())
			walls.POST("/walls", handler.Walls())
			walls.POST("/blocks", handler.Blocks())
			walls.POST("/programs", handler.Programs())
			walls.POST("/episodes", handler.Episodes())
			walls.POST("/medias", handler.Medias())
			walls.POST("/all", handler.All())
		}

	}

	// Return the configured router.
	return r
}

// removeHTTPS removes the "https://" prefix from a URL string.
func removeHTTPS(url string) string {
	// Check if the URL starts with "https://"
	if strings.HasPrefix(url, "https://") {
		// Remove "https://" from the URL
		return strings.TrimPrefix(url, "https://")
	}
	// If the URL doesn't start with "https://", return it as is
	return url
}
