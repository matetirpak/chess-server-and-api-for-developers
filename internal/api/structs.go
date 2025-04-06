/*
Structs for API communication are defined here.
*/
package api

// Create new game
type ReqPostSessions struct {
	Name string `json:"name"`
}
type RespPostSessions struct {
	BoardID  int32  `json:"boardid"`
	Password string `json:"password"`
}

// Delete an ongoing game
type ReqDeleteSessions struct {
	BoardID  int32  `json:"boardid"`
	Password string `json:"password"`
}

// Entry a game
type ReqPutSessions struct {
	BoardID  int32  `json:"boardid"`
	Password string `json:"password"`
	Color    string `json:"color"`
}
type RespPutSessions struct {
	Token string `json:"token"`
}

// Get all sessions
type RespGetSessions struct {
	Games []GameNameAndID `json:"games"`
}
type GameNameAndID struct {
	Name    string `json:"name"`
	BoardID int32  `json:"boardid"`
}

// Get game data
type ReqGetGame struct {
	Moveidx  int32  `schema:"moveidx"`
	BoardID  int32  `schema:"boardid"`
	Password string `schema:"password"`
	Color    string `schema:"color"`
	Token    string `schema:"token"`
	Statereq bool   `schema:"statereq"`
	Turnreq  bool   `schema:"turnreq"`
}

// Apply move
type ReqPutGame struct {
	BoardID  int32  `json:"boardid"`
	Password string `json:"password"`
	Color    string `json:"color"`
	Token    string `json:"token"`
	Move     string `json:"move,omitempty"`
	Forfeit  bool   `json:"forfeit,omitempty"`
}
