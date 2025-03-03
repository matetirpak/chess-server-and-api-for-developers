/*
API for existing games.
*/
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/matetirpak/chess-server-and-api-for-developers/internal/game_logic"

	db "github.com/matetirpak/chess-server-and-api-for-developers/internal/database"
)

// Get the state of a game
func GetGame(w http.ResponseWriter, r *http.Request) {
	/*
		Input:
			Board ID, password, color and associated token
			are required.
			Either 'Statereq' or 'Turnreq' have to be true.

			ReqdataGetBoard
			BoardID  int32  `json:"board-id,omitempty"`
			Password string `json:"password,omitempty"`
			Color    string `json:"color,omitempty"`
			Token    string `json:"token,omitempty"`
			Statereq int32  `json:"statereq,omitempty"`
			Turnreq  int32  `json:"turnreq,omitempty"`
		Return:
			Either board information or a notification when
			it's the player's turn.

			If Statereq:
				game.BoardData of type game_logic.BoardState
			If Turnreq:
				map[string]string{"message": "It's your turn!"}
		Actions:
			If Turnreq:
				Holds the call till it's the player's turn, then
				sends a response as a notification.
	*/

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var req ReqGetGame

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	if req.Statereq == req.Turnreq {
		http.Error(w, "'statereq' or 'turnreq' has to be true exclusively.", http.StatusBadRequest)
		return
	}

	success1 := verifyGameAccess(w, req.BoardID, req.Password)
	if !success1 {
		return
	}
	var game *db.Game = db.GamesMap[req.BoardID]

	client_ip, err := extractClientIP(w, r)
	if err != nil {
		return
	}

	success2 := verifyBoardAccess(w, game, client_ip, req.Color, req.Token)
	if !success2 {
		return
	}

	if req.Statereq {
		resp := game.BoardData
		json.NewEncoder(w).Encode(resp)
		return
	}
	if req.Turnreq {
		w.Header().Set("Connection", "keep-alive")

		timeout := time.After(60 * time.Second)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				http.Error(w, `{"message":"Timeout waiting for your turn"}`, http.StatusRequestTimeout)
				return
			case <-ticker.C:
				// Check the current player's turn
				game.Mu.RLock()
				currentTurn := game.PlayerTurn
				game.Mu.RUnlock()

				if currentTurn == req.Color {
					response := map[string]string{"message": "It's your turn!"}
					json.NewEncoder(w).Encode(response)
					return
				}
			}
		}
	}
}

// Endpoint to update the player turn for testing purposes
func UpdateTurn(w http.ResponseWriter, r *http.Request) {
	type Update struct {
		BoardID int32  `json:"board-id"`
		Turn    string `json:"turn"` // "w" or "b"
	}
	var req Update
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Invalid request: %v"}`, err), http.StatusBadRequest)
		return
	}
	game, exists := db.GamesMap[req.BoardID]
	if !exists {
		http.Error(w, "Board does not exist.", http.StatusBadRequest)
		return
	}
	// Update the player turn
	game.Mu.Lock()
	game.PlayerTurn = req.Turn
	game.Mu.Unlock()

	resp := fmt.Sprintf(`{"message":"Turn updated to: %s"}`, req.Turn)
	json.NewEncoder(w).Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// Apply a move
func PutGame(w http.ResponseWriter, r *http.Request) {
	/*
		Input:
			Board ID, password, color and associated token,
			and move to be made.

			ReqdataPutBoard
			BoardID  int32  `json:"board-id,omitempty"`
			Password string `json:"password,omitempty"`
			Color    string `json:"color,omitempty"`
			Token    string `json:"token,omitempty"`
			Move     string `json:"move,omitempty"`
		Return:
			---
		Actions:
			Checks validity of the move and applies
			it to the board.
	*/
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var req ReqPutGame
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	success1 := verifyGameAccess(w, req.BoardID, req.Password)
	if !success1 {
		return
	}
	var game *db.Game = db.GamesMap[req.BoardID]

	client_ip, err := extractClientIP(w, r)
	if err != nil {
		return
	}

	success2 := verifyBoardAccess(w, game, client_ip, req.Color, req.Token)
	if !success2 {
		return
	}

	if game.Winner != 'n' {
		http.Error(w, "Can't apply move. Game has ended.", http.StatusBadRequest)
		return
	}

	move, err := game_logic.StringToMoveStruct(req.Move, rune(req.Color[0]))
	if err != nil {
		http.Error(w, "Move format is invalid.", http.StatusBadRequest)
		return
	}

	// Check validity of move
	err = game_logic.ValidateMove(&move, &game.BoardData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Move is invalid with error: %v", err), http.StatusBadRequest)
		return
	}

	game_logic.MakeMove(&move, &game.BoardData)

	w.WriteHeader(http.StatusOK)
}
