package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// r.POST("/create-room", createRoomHandle)

	log.Fatal(r.Run())
}
