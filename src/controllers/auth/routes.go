// platform/router/router.go

package auth

import (
	"encoding/gob"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte(config.Get().Oauth.Secret))
	router.Use(sessions.Sessions("auth-session", store))

	router.GET("/login", loginHandler)
	router.GET("/callback", callbackHandler)
	router.GET("/logout", logoutHandler)
}
