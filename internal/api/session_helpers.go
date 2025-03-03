/*
Helper functions for API communication and game setup.
*/
package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"

	db "github.com/matetirpak/chess-server-and-api-for-developers/internal/database"
)

// Verifies whether a user has access to a session.
func verifyGameAccess(w http.ResponseWriter, id int32, password string) bool {
	game, exists := db.GamesMap[id]
	if !exists {
		http.Error(w, "Board not found", http.StatusNotFound)
		return false
	}
	if password != game.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return false
	}
	return true
}

// Verifies whether a user is registered as a player and has access to a session.
func verifyBoardAccess(w http.ResponseWriter, game *db.Game, client_ip string, color string, token string) bool {
	// Both ip and token have to match with the database, session password is not required.
	switch color {
	case "w":
		if !game.HasWPlayer {
			http.Error(w, "White player does not exist.", http.StatusNotFound)
			return false
		}
		if game.W_playerIP != client_ip {
			http.Error(w, "IP is not authorized.", http.StatusUnauthorized)
			return false
		}
		if game.W_playerToken != token {
			http.Error(w, "Token is invalid.", http.StatusUnauthorized)
			return false
		}
		return true

	case "b":
		if !game.HasBPlayer {
			http.Error(w, "Black player does not exist.", http.StatusNotFound)
			return false
		}
		if game.B_playerIP != client_ip {
			http.Error(w, "IP is not authorized.", http.StatusUnauthorized)
			return false
		}
		if game.B_playerToken != token {
			http.Error(w, "Token is invalid.", http.StatusUnauthorized)
			return false
		}
		return true

	default:
		http.Error(w, "Invalid color. Enter 'w' or 'b'", http.StatusBadRequest)
		return false
	}
}

func extractClientIP(w http.ResponseWriter, r *http.Request) (string, error) {
	// Split the RemoteAddr into host (IP) and port
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when extracting IP: %v", err), http.StatusInternalServerError)
		return "", err
	}
	return ip, nil
}

func generateToken() string {
	return uuid.New().String()
}

func initializeNewGame(name string) *db.Game {
	var game db.Game
	game.Name = name
	game.ID = db.Ids
	db.Ids++
	game.Password = generateToken()
	game.HasWPlayer = false
	game.HasBPlayer = false
	game.PlayerTurn = "NotStarted"
	game.Winner = 'n'
	return &game
}
