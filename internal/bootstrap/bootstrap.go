package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-indexer-api/internal/configuration"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/api"
	"github.com/khedhrije/podcaster-indexer-api/internal/infrastructure/elasticsearchv7"
	"github.com/khedhrije/podcaster-indexer-api/internal/infrastructure/mysql"
	"github.com/khedhrije/podcaster-indexer-api/internal/ui/gin/handlers"
	"github.com/khedhrije/podcaster-indexer-api/internal/ui/gin/router"
	"github.com/rs/zerolog/log"
)

// Bootstrap struct encapsulates the configuration settings and the HTTP router necessary for the application to run.
type Bootstrap struct {
	Config *configuration.AppConfig // Application configuration settings
	Router *gin.Engine              // HTTP router for handling web requests
}

// InitBootstrap initializes the bootstrap process and returns a Bootstrap instance.
// It serves as a public entry point for the initialization process.
func InitBootstrap() Bootstrap {
	return initBootstrap()
}

// initBootstrap sets up the Bootstrap struct by initializing configurations, database connections, middleware, and routes.
// It panics if the configuration is not set, ensuring that the application does not run with nil configurations.
func initBootstrap() Bootstrap {
	if configuration.Config == nil {
		log.Panic().Msg("configuration is nil")
	}

	app := Bootstrap{}
	app.Config = configuration.Config

	// Initialize MySQL client using application configuration
	mysqlClient := mysql.NewClient(app.Config)

	// Initialize Meilisearch client using application configuration
	esClient, err := elasticsearchv7.NewElasticSearchClient(app.Config)
	if err != nil {
		panic("could not init elasticsearch client")
	}

	// Initialize MySQL adapters for different domain models, setting up the data access layer
	catPersistenceAdapter := mysql.NewCategoryAdapter(mysqlClient)
	tagPersistenceAdapter := mysql.NewTagAdapter(mysqlClient)
	wallPersistenceAdapter := mysql.NewWallAdapter(mysqlClient)
	blockPersistenceAdapter := mysql.NewBlockAdapter(mysqlClient)
	programPersistenceAdapter := mysql.NewProgramAdapter(mysqlClient)
	episodePersistenceAdapter := mysql.NewEpisodeAdapter(mysqlClient)
	mediaPersistenceAdapter := mysql.NewMediaAdapter(mysqlClient)

	// Initialize APIs for different domain models, enabling business logic operations
	indexationApi := api.NewIndexerApi(esClient, wallPersistenceAdapter, catPersistenceAdapter, tagPersistenceAdapter, blockPersistenceAdapter, programPersistenceAdapter, episodePersistenceAdapter, mediaPersistenceAdapter)

	// Initialize handlers for different APIs, setting up the presentation layer
	indexationHandler := handlers.NewIndexationHandler(indexationApi)

	// Create the router with the initialized handlers, configuring the request handling
	r := router.CreateRouter(
		indexationHandler,
	)
	app.Router = r
	return app
}

// Run starts the application by running the HTTP server on the configured host address and port.
// It logs a fatal error if the server cannot be started, ensuring that the failure is captured and reported.
func (b Bootstrap) Run() {
	dsn := fmt.Sprintf("%s:%d", b.Config.HostAddress, b.Config.HostPort)
	if errRun := b.Router.Run(dsn); errRun != nil {
		log.Fatal().Msg("error during service instantiation")
	}
}
