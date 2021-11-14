package server

import (
	"github.com/gin-gonic/gin"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
	"github.com/karthiklsarma/cedar-logging/logging"
)

func InitiateServerEntry() {
	logging.SetInfoLogLevel()
	logging.Info("Initializing server entry")
	router := gin.Default()
	setupRouting(router)
	router.Run(":8080")
}

func setupRouting(router *gin.Engine) {
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server has started",
		})
	})

	handlerFunc := graphqlHandler()
	router.POST("/graphql", handlerFunc)
	router.OPTIONS("/graphql", handlerFunc)
}

func graphqlHandler() gin.HandlerFunc {
	gqlSchema := StartGraphQlServer()
	gqHandler := gqlhandler.New(&gqlhandler.Config{
		Schema: &gqlSchema,
		Pretty: true,
	})
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Max-Age", "10000")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		gqHandler.ContextHandler(c, c.Writer, c.Request)
	}
}
