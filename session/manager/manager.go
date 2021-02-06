package manager

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/saitofun/qor/middlewares"
	"github.com/saitofun/qor/session"
	"github.com/saitofun/qor/session/gorilla"
)

// SessionManager default session manager
var SessionManager session.ManagerInterface = gorilla.New("_session", sessions.NewCookieStore([]byte("secret")))

func init() {
	middlewares.Use(middlewares.Middleware{
		Name: "session",
		Handler: func(handler http.Handler) http.Handler {
			return SessionManager.Middleware(handler)
		},
	})
}
