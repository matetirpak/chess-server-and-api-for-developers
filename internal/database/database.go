/*
Manages the data storage of the server.
*/
package database

import (
	"sync"

	"github.com/matetirpak/chess-server-and-api-for-developers/internal/game_logic"
)

type Game struct {
	Name          string
	ID            int32
	Password      string
	HasWPlayer    bool
	W_playerIP    string
	W_playerToken string
	HasBPlayer    bool
	B_playerIP    string
	B_playerToken string
	PlayerTurn    string
	Winner        rune
	Mu            sync.RWMutex
	BoardData     game_logic.BoardState
}

var GamesMap = make(map[int32]*Game)

var Ids int32 = 1
