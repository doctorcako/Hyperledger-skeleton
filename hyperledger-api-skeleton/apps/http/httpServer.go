package http

import (
	"context"
	httpController "hyperledger-api-skeleton/hyperledger-business-logic/internal/controllers/http"
	"hyperledger-api-skeleton/hyperledger-business-logic/internal/core/ports"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/errors"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/properties"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/sdk/fabricGateway"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customError"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customHttp"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customLogger"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

type Binder struct {
	Cxt           context.Context
	Properties    *properties.Properties
	LogCustom     customLogger.Log
	ServerHttp    ServerHttp
	FabricGateway fabricGateway.GatewayFabric
	StatusService ports.StatusService
	KafkaInput    ports.InputEvents
}

type ServerHttp struct {
	Router *gin.Engine
}

// @title           Swagger API SOAR Business Logic
// @version         1.0
// @description     Hyperledger API Business Logic
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func StartApp(b Binder) customError.Error {
	headerLog := "Http Server - StartApp - "

	b.ServerHttp.Router = CORS(b.ServerHttp.Router)
	b.ServerHttp.Router = Routes(b)
	b.LogCustom.Info(headerLog, b.Properties.ApplicationProperties.HttpServer.Url+":"+b.Properties.ApplicationProperties.HttpServer.Port)

	b.ServerHttp.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	b.ServerHttp.Router.GET("/debug/pprof/*any", gin.WrapF(http.DefaultServeMux.ServeHTTP))

	err := b.ServerHttp.Router.Run(b.Properties.ApplicationProperties.HttpServer.Url + ":" + b.Properties.ApplicationProperties.HttpServer.Port)
	return errors.NewCustomError(customError.InternalError, err.Error())
}

// Routes configures the server routers
func Routes(binder Binder) *gin.Engine {
	router := binder.ServerHttp.Router

	router.Use(
		logger.SetLogger(),                 // Log API request calls
		gzip.Gzip(gzip.DefaultCompression), // Compress results, mostly gzipping assets and json
		gin.Recovery(),                     // Recover from panics without crashing server
		//middleware.RequestID(),             // Assigns a unique request ID
		customHttp.NoSniffGin,
	)

	publicAPI := router.Group("/status")
	{
		publicAPI.Use(apiVersionCtx("v1"))
		statusHandler := publicAPI.Group("/")
		{
			httpController.NewStatusHandler(binder.LogCustom, statusHandler, binder.StatusService)
		}
	}

	privateAPI := router.Group("/v1")
	{
		privateAPI.Use(apiVersionCtx("v1"))
		networkManagerHandler := privateAPI.Group("/soar")
		{
			httpController.NewSoarWebhookHandler(binder.LogCustom, networkManagerHandler, binder.SoarService)
		}
	}

	return router
}

func CORS(router *gin.Engine) *gin.Engine {
	corsOptions := cors.New(cors.Config{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of the major browsers
	})
	router.Use(corsOptions)
	return router
}

func apiVersionCtx(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("api.version", version)
		c.Next()
	}
}
