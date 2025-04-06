/*
API for existing games.
*/
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/schema"
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

			ReqGetGame
			Moveidx  int32  `json:"moveidx,omitempty"`
			BoardID  int32  `json:"board-id"`
			Password string `json:"password"`
			Color    string `json:"color"`
			Token    string `json:"token"`
			Statereq int32  `json:"statereq"`
			Turnreq  int32  `json:"turnreq"`
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
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		http.Error(w, "Failed to parse query params: "+err.Error(), http.StatusBadRequest)
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

	success2 := verifyBoardAccess(w, game, req.Color, req.Token)
	if !success2 {
		return
	}

	if req.Statereq {
		var idx int = int(req.Moveidx)
		if idx == -1 {
			idx = len(game.BoardData) - 1
		}
		resp := game.BoardData[idx]
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
				currentTurn := game.BoardData[len(game.BoardData)-1].TurnColor
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
	game.BoardData[len(game.BoardData)-1].TurnColor = req.Turn
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

			ReqPutGame
			BoardID  int32  `json:"board-id"`
			Password string `json:"password"`
			Color    string `json:"color"`
			Token    string `json:"token"`
			Move     string `json:"move,omitempty"`
			Forfeit  bool   `json:"forfeit,omitempty"`
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

	success2 := verifyBoardAccess(w, game, req.Color, req.Token)
	if !success2 {
		return
	}

	if req.Forfeit {
		game.BoardData[len(game.BoardData)-1].TurnColor = "n"
		if req.Color == "w" {
			game.Winner = "b"
			game.BoardData[len(game.BoardData)-1].Winner = "b"
		} else {
			game.Winner = "w"
			game.BoardData[len(game.BoardData)-1].Winner = "w"
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if !game.Started {
		http.Error(w, "Can't apply move. Game has not started.", http.StatusBadRequest)
		return
	}

	if game.Winner != "n" {
		http.Error(w, "Can't apply move. Game has ended.", http.StatusBadRequest)
		return
	}

	move, err := game_logic.StringToMoveStruct(req.Move, rune(req.Color[0]))
	if err != nil {
		http.Error(w, "Move format is invalid.", http.StatusBadRequest)
		return
	}

	// Check validity of move
	err = game_logic.ValidateMove(&move, &game.BoardData[len(game.BoardData)-1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Move is invalid with error: %v", err), http.StatusBadRequest)
		return
	}

	newBstate := game_logic.MakeMove(&move, game.BoardData[len(game.BoardData)-1])
	game.BoardData = append(game.BoardData, newBstate)

	w.WriteHeader(http.StatusOK)
}
