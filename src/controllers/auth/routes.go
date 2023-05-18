// platform/router/router.go

package auth

import (
	"encoding/gob"

	"github.com/JakubC-projects/myshare-activity-telegram/src/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func AddRoutes(router *gin.Engine) {
	auth, _ := auth.GetAuthenticator()
	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte(config.Get().Oauth.Secret))
	router.Use(sessions.Sessions("auth-session", store))

	router.GET("/login", loginHandler(auth))
	router.GET("/callback", callbackHandler(auth))
}
