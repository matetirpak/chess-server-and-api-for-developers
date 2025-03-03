package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/matetirpak/chess-server-and-api-for-developers/internal/api"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	// Sessions
	Route{
		"DeleteSessions",
		strings.ToUpper("Delete"),
		"/ChessServer/0.1.0/sessions",
		api.DeleteSessions,
	},

	Route{
		"GetSessions",
		strings.ToUpper("Get"),
		"/ChessServer/1.0.0/sessions",
		api.GetSessions,
	},

	Route{
		"PostSessions",
		strings.ToUpper("Post"),
		"/ChessServer/0.1.0/sessions",
		api.PostSessions,
	},

	Route{
		"PutSessions",
		strings.ToUpper("Put"),
		"/ChessServer/0.1.0/sessions",
		api.PutSessions,
	},

	// Game
	Route{
		"GetGame",
		strings.ToUpper("Get"),
		"/ChessServer/0.1.0/game",
		api.GetGame,
	},

	Route{
		"PutGame",
		strings.ToUpper("Put"),
		"/ChessServer/0.1.0/game",
		api.PutGame,
	},

	// Debugging
	Route{
		"UpdateTurn",
		strings.ToUpper("Put"),
		"/ChessServer/0.1.0/board/updateturn",
		api.UpdateTurn,
	},
}
