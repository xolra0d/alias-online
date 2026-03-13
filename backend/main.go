package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	pgPool, err := InitPool()
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()
	postgres := &Postgres{pgPool}
	handles := &Handles{postgres}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	{
		api := router.Group("/api")
		api.GET("/ok", func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})
		api.GET("/available-vocabs", handles.AvailableLanguages)
		api.POST("/create-user", handles.CreateUser)

		{
			protected := api.Group("/protected")
			protected.Use(handles.UserAuthMiddleware())
			protected.GET("/ok", func(c *gin.Context) {
				c.JSON(200, gin.H{"ok": true})
			})
			protected.POST("/create-room", handles.CreateRoom)

			protected.POST("/ws", handles.InitWS)

			{
				//ws := protected.Group("/ws")

			}
		}
	}

	log.Fatal(router.Run())
}
