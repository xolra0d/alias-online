package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//type CreateRoomConfig struct {
//	Language               string `form:"language"`
//	RudeWords              string `form:"allow-rude-words,default=off"`
//	OnlyExternalVocabulary string `form:"only-external-vocabulary,default=off"`
//	AdditionalVocabulary   string `form:"additional-vocabulary"`
//	Clock                  int    `form:"clock,default=60"`
//}

//func createRoomHandle(c *gin.Context) {
//	var data CreateRoomConfig
//	if err := c.ShouldBind(&data); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	// c.Redirect(http.StatusSeeOther, "/dsadsadsda")
//}

type Handles struct {
	postgres  *Postgres
	websocket *Rooms
}

func (h *Handles) AvailableLanguages(ctx *gin.Context) {
	languages, err := h.postgres.AvailableLanguages(ctx.Request.Context())
	if err != nil {
		ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"ok": true, "languages": languages})
}

func (h *Handles) CreateUser(ctx *gin.Context) {
	name, credentials, err := h.postgres.CreateUser(ctx.Request.Context()) // TODO: consider using `runtime.secret`
	if err != nil {
		ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"ok": true, "credentials": credentials, "name": name})
}

func (h *Handles) UserAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Request.Header.Get("User-Id")
		secret := ctx.Request.Header.Get("User-Secret")
		if idString == "" || secret == "" {
			ctx.JSON(401, gin.H{"ok": false, "reason": "missing credentials"})
			ctx.Abort()
			return
		}
		idUUID, err := uuid.Parse(idString)
		if err != nil {
			ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
			ctx.Abort()
			return
		}

		credentials := UserCredentials{idUUID, secret}
		err = h.postgres.ValidateUser(ctx.Request.Context(), credentials)
		if err != nil {
			ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

type CreateRoomConfig struct {
	Language             string `form:"language"`
	RudeWords            bool   `form:"rude-words"`
	AdditionalVocabulary string `form:"additional-vocabulary"`
	Clock                int    `form:"clock"`
}

func (h *Handles) CreateRoom(ctx *gin.Context) {
	adminId := uuid.MustParse(ctx.Request.Header.Get("User-Id")) // verified at `UserAuthMiddleware`
	var data CreateRoomConfig
	err := ctx.Bind(&data)
	if err != nil {
		ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
		return
	}

	v := strings.Split(data.AdditionalVocabulary, ",")
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}

	roomId, err := h.postgres.AddRoom(ctx.Request.Context(), adminId, data.Language, data.RudeWords, v, data.Clock)
	if err != nil {
		ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"ok": true, "room_id": roomId})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 128,
}

func (h *Handles) InitWS(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	if roomId == "" {
		ctx.JSON(200, gin.H{"ok": false, "reason": "missing room id"})
	}
	if err := h.websocket.ServeHTTP(ctx.Writer, ctx.Request, roomId); err != nil {
		ctx.JSON(200, gin.H{"ok": false, "reason": err.Error()})
	}
}
