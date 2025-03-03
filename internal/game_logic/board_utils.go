/*
Board functionalities and helper functions used to extract
board information are implemented here.
*/

package game_logic

import (
	"strings"
)

// BoardState represents the chessboard and its metadata
type BoardState struct {
	/*
		Winner: 'n': none, 'r': remis, 'w': white, 'b': black.
		EnPassant: If a pawn double moves, coordinates of
				   the target field. Otherwise {-1,-1}.
		TurnColor: Current player. 'w': white, 'b': black.
	*/
	Board          [8][8]rune `json:"board"`
	WhiteKingPos   [2]int     `json:"whitekingpos"`
	BlackKingPos   [2]int     `json:"blackkingpos"`
	WhiteKingMoved bool       `json:"whitekingmoved"`
	BlackKingMoved bool       `json:"blackkingmoved"`
	Winner         rune       `json:"winner"`
	EnPassant      [2]int     `json:"enpassant"`
	TurnColor      rune       `json:"turncolor"`
}

// Constructs the standard starting board.
func InitializeBoard(boardState *BoardState) {
	boardState.WhiteKingPos = [2]int{7, 4}
	boardState.BlackKingPos = [2]int{0, 4}
	boardState.WhiteKingMoved = false
	boardState.BlackKingMoved = false
	boardState.Winner = 'n'
	boardState.EnPassant = [2]int{-1, -1}
	boardState.TurnColor = 'w'

	board := &boardState.Board

	// Initialize black pieces
	board[0][0] = 'R'
	board[0][1] = 'K'
	board[0][2] = 'B'
	board[0][3] = 'Q'
	board[0][4] = 'X'
	board[0][5] = 'B'
	board[0][6] = 'K'
	board[0][7] = 'R'
	for i := 0; i < 8; i++ {
		board[1][i] = 'P'
	}

	// Initialize empty fields
	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			board[i][j] = Empty
		}
	}

	// Initialize white pieces
	for i := 0; i < 8; i++ {
		board[6][i] = 'p'
	}
	board[7][0] = 'r'
	board[7][1] = 'k'
	board[7][2] = 'b'
	board[7][3] = 'q'
	board[7][4] = 'x'
	board[7][5] = 'b'
	board[7][6] = 'k'
	board[7][7] = 'r'
}

// Returns the color and piece given a position on the board.
func getColorAndPiece(row int, col int, board [8][8]rune) (color rune, piece rune) {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		// Field out of bounds.
		// Error is not raised for simplicity.
		return 'n', Empty
	}

	// Extract piece
	p := board[row][col]
	if p == Empty {
		return 'n', Empty
	}

	// Return color and piece in lowercase
	if p >= 'a' && p <= 'z' {
		return 'w', rune(strings.ToLower(string(p))[0])
	}
	return 'b', rune(strings.ToLower(string(p))[0])
}

// Checks whether a position is within the board's bounds.
func isInBounds(pos [2]int) bool {
	return pos[0] >= 0 && pos[0] < 8 && pos[1] >= 0 && pos[1] < 8
}
