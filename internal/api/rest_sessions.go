/*
API for session management.
*/
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/matetirpak/chess-server-and-api-for-developers/internal/database"
)

// Deletes an ongoing game
func DeleteSessions(w http.ResponseWriter, r *http.Request) {
	/*
		Input:
			Board ID and password.

			ReqDeleteSessions
			BoardID  int32  `json:"board-id,omitempty"`
			Password string `json:"password,omitempty"`
		Return:
			---
		Actions:
			Deletes the game session if it exists
			and the password is correct.
	*/
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var req ReqDeleteSessions
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Decoding failed.\n")
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	success := verifyGameAccess(w, req.BoardID, req.Password)
	if !success {
		return
	}

	delete(db.GamesMap, req.BoardID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}")) // Empty response
}

// Displays all existing games
func GetSessions(w http.ResponseWriter, r *http.Request) {
	/*
		Input:
			---
		Return:
			List of all existing {name, game ID} pairs.

		Actions:
			---
	*/
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var resp RespGetSessions

	// Iterate through the GamesMap and populate the response
	for _, game := range db.GamesMap {
		extracted_game := GameNameAndID{
			Name:    game.Name,
			BoardID: game.ID,
		}
		// Append the response to the slice
		resp.Games = append(resp.Games, extracted_game)
	}
	json.NewEncoder(w).Encode(resp)
}

// Creates a new game
func PostSessions(w http.ResponseWriter, r *http.Request) {
	/*
		Input:
			Name of the game to be created.

			ReqPostSessions
			Name string `json:"name,omitempty"`
		Return:
			Unique ID and password to access the game.

			RespPostSessions
			BoardID  int32  `json:"board-id,omitempty"`
			Password string `json:"password,omitempty"`
		Actions:
			Create a game entry.
	*/

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var req ReqPostSessions
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	NewGame := initializeNewGame(req.Name)

	db.GamesMap[NewGame.ID] = NewGame

	resp := RespPostSessions{BoardID: NewGame.ID, Password: NewGame.Password}
	json.NewEncoder(w).Encode(resp)
}

// Entries a games
func PutSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	/*
		Input:
			Game ID, password and desired color.

			ReqPutSessions
			BoardID  int32  `json:"board-id,omitempty"`
			Password string `json:"password,omitempty"`
			Color    string `json:"color,omitempty"``
		Return:
			Token to access the game as the player of
			the specified color.

			RespPutSessions
			Token string `json:"token,omitempty"`
		Actions:
			The color in the specified game is reserved for
			the caller and protected by his token.
	*/
	var req ReqPutSessions
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	success := verifyGameAccess(w, req.BoardID, req.Password)
	if !success {
		return
	}

	var game *db.Game = db.GamesMap[req.BoardID]

	if game.HasWPlayer && game.HasBPlayer {
		http.Error(w, "Game is already full.", http.StatusForbidden)
		return
	}

	token := generateToken()
	resp := RespPutSessions{Token: token}

	switch req.Color {
	case "w":
		if game.HasWPlayer {
			http.Error(w, "White is already taken.", http.StatusForbidden)
			return
		}
		game.HasWPlayer = true
		game.W_playerToken = token
		if game.HasBPlayer {
			game.Started = true
		}
		json.NewEncoder(w).Encode(resp)

	case "b":
		if game.HasBPlayer {
			http.Error(w, "Black is already taken.", http.StatusForbidden)
			return
		}
		game.HasBPlayer = true
		game.B_playerToken = token
		if game.HasWPlayer {
			game.Started = true
		}
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Invalid color. Enter 'w' or 'b'", http.StatusBadRequest)
		return
	}
}
