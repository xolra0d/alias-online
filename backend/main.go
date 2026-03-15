package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	pgPool, err := InitPool()
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()
	postgres := &Postgres{pgPool}
	vocabs, err := postgres.LoadVocabs(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	handles := &Handles{postgres, &Rooms{rooms: map[string]*Room{}}, &Vocabularies{vocabulary: vocabs}}

	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return true // todo change
	}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "User-Id", "User-Secret"}

	router := gin.Default()
	router.Use(cors.New(config))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	{
		api := router.Group("/api")
		api.GET("/ok", func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})
		api.GET("/available-vocabs", handles.AvailableLanguages)
		api.POST("/create-user", handles.CreateUser)
		api.GET("/ws/:roomId", handles.InitWS)

		{
			protected := api.Group("/protected")
			protected.Use(handles.UserAuthMiddleware())
			protected.GET("/ok", func(c *gin.Context) {
				c.JSON(200, gin.H{"ok": true})
			})
			protected.POST("/create-room", handles.CreateRoom)
		}

		{
			admin := api.Group("/admin")
			admin.Use(handles.AdminAuthMiddleware())
			admin.POST("/refresh-vocabularies", handles.RefreshVocabularies)
		}
	}

	log.Fatal(router.Run())
}
