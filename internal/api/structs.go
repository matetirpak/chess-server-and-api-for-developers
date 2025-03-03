/*
Structs for API communication are defined here.
*/
package api

// Create new game
type ReqPostSessions struct {
	Name string `json:"name,omitempty"`
}
type RespPostSessions struct {
	BoardID  int32  `json:"board-id,omitempty"`
	Password string `json:"password,omitempty"`
}

// Delete an ongoing game
type ReqDeleteSessions struct {
	BoardID  int32  `json:"board-id"`
	Password string `json:"password"`
}

// Entry a game
type ReqPutSessions struct {
	BoardID  int32  `json:"board-id"`
	Password string `json:"password"`
	Color    string `json:"color"`
}
type RespPutSessions struct {
	Token string `json:"token,omitempty"`
}

// Get all sessions
type RespGetSessions struct {
	Games []GameNameAndID `json:"games,omitempty"`
}
type GameNameAndID struct {
	Name    string `json:"name,omitempty"`
	BoardID int32  `json:"board-id,omitempty"`
}

// Get game data
type ReqGetGame struct {
	BoardID  int32  `json:"board-id"`
	Password string `json:"password"`
	Color    string `json:"color"`
	Token    string `json:"token"`
	Statereq bool   `json:"statereq"`
	Turnreq  bool   `json:"turnreq"`
}

// Apply move
type ReqPutGame struct {
	BoardID  int32  `json:"board-id"`
	Password string `json:"password"`
	Color    string `json:"color"`
	Token    string `json:"token"`
	Move     string `json:"move"`
}
