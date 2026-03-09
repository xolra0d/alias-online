package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRoomData struct {
	Language               string `form:"language"`
	RudeWords              string `form:"allow-rude-words,default=off"`
	OnlyExternalVocabulary string `form:"only-external-vocabulary,default=off"`
	AdditionalVocabulary   string `form:"additional-vocabulary"`
	Clock                  int    `form:"clock,default=60"`
}

func availableLanguages() []string {
	return []string{"en", "ru"}
}


func createRoomHandle(c *gin.Context) {
	var data CreateRoomData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// c.Redirect(http.StatusSeeOther, "/dsadsadsda")
}
