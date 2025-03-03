/*
This module implements move functionalities including
move formatting, move generation and move validation.
*/

package game_logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

const Empty = ' '

// Move represents a chess move.
type Move struct {
	from    [2]int // [row, col]
	to      [2]int // [row, col]
	color   rune   // 'w' or 'b'
	capture bool
	next    *Move
}

// Comparator for Move
func eqMove(move1, move2 *Move) bool {
	if move1 == nil || move2 == nil {
		return move1 == move2 // Both must be nil to be equal.
	}

	return move1.from == move2.from &&
		move1.to == move2.to &&
		move1.color == move2.color
}

// Checks whether a move is part of a move-list
func isMoveInMoves(move *Move, moves *Move) bool {
	for moves != nil {
		if eqMove(move, moves) {
			return true
		}
		moves = moves.next
	}
	return false
}

// AllPossibleMoves generates all moves for the given player and evaluates board value
func allPossibleMoves(color rune, boardState *BoardState, exclude []rune) *Move {
	var head *Move
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			target_color, target_piece := getColorAndPiece(row, col, boardState.Board)
			if target_piece == Empty {
				continue
			}
			if color == target_color {
				if isExcluded(target_piece, exclude) {
					continue
				}
				generateMovesForPiece(row, col, boardState, &head)
			}
		}
	}
	return head
}

func isExcluded(piece rune, exclude []rune) bool {
	for _, excluded := range exclude {
		if piece == excluded {
			return true
		}
	}
	return false
}

// generateMovesForPiece generates moves for a specific piece
func generateMovesForPiece(row, col int, boardState *BoardState, moves **Move) {
	color, piece := getColorAndPiece(row, col, boardState.Board)

	switch piece {
	case 'p', 'P': // Pawn
		pawnMoves(row, col, color, boardState, moves)
	case 'k', 'K': // Knight
		knightMoves(row, col, color, boardState, moves)
	case 'b', 'B': // Bishop
		lineMoves(row, col, color, boardState, moves, "diagonal")
	case 'r', 'R': // Rook
		lineMoves(row, col, color, boardState, moves, "straight")
	case 'q', 'Q': // Queen
		lineMoves(row, col, color, boardState, moves, "diagonal")
		lineMoves(row, col, color, boardState, moves, "straight")
	case 'x', 'X': // King
		kingMoves(row, col, color, boardState, moves)
	}
}

// pawnMoves generates moves for a pawn
func pawnMoves(row, col int, color rune, boardState *BoardState, moves **Move) {
	direction := -1
	if color == 'b' {
		direction = 1
	}
	board := boardState.Board

	// Single step forward
	if isValid(row+direction, col) && board[row+direction][col] == Empty {
		addMove(row, col, row+direction, col, color, false, moves)
	}

	// Double step on initial position
	startRow := 6
	if color == 'b' {
		startRow = 1
	}
	if row == startRow && isValid(row+2*direction, col) && board[row+direction][col] == Empty && board[row+2*direction][col] == Empty {
		addMove(row, col, row+2*direction, col, color, false, moves)
	}

	// Capture moves
	for _, offset := range []int{-1, 1} {
		if isValid(row+direction, col+offset) {
			target_color, target_piece := getColorAndPiece(row+direction, col+offset, board)
			if target_piece == Empty {
				continue
			}
			if color != target_color {
				addMove(row, col, row+direction, col+offset, color, true, moves)
			}
		}
	}
}

// knightMoves generates moves for a knight
func knightMoves(row, col int, color rune, boardState *BoardState, moves **Move) {
	knightOffsets := [][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}

	for _, offset := range knightOffsets {
		r, c := row+offset[0], col+offset[1]
		if !isValid(r, c) {
			continue
		}
		target_color, target_piece := getColorAndPiece(r, c, boardState.Board)
		if target_piece == Empty {
			addMove(row, col, r, c, color, false, moves)
		} else if color != target_color {
			addMove(row, col, r, c, color, true, moves)
		}
	}
}

func isValid(row, col int) bool {
	return row >= 0 && row < 8 && col >= 0 && col < 8
}

func lineMoves(start_row, start_col int, color rune, boardState *BoardState,
	moves **Move, direction string) {
	getDeltas := func(direction string) [4][2]int {
		if direction != "straight" && direction != "diagonal" {
			// error handling required
			return [4][2]int{}
		}
		if direction == "straight" {
			return [4][2]int{
				{-1, 0}, // N
				{1, 0},  // S
				{0, -1}, // W
				{0, 1},  // E
			}
		}
		if direction == "diagonal" {
			return [4][2]int{
				{-1, -1}, // NW
				{-1, 1},  // NE
				{1, -1},  // SW
				{1, 1},   // SE
			}
		}
		return [4][2]int{}
	}
	deltas := getDeltas(direction)

	for _, delta := range deltas {
		row, col := start_row, start_col
		for {
			row += delta[0]
			col += delta[1]
			if !isValid(row, col) {
				break
			}

			target_color, target_piece := getColorAndPiece(row, col, boardState.Board)
			if color == target_color {
				break
			}
			if target_piece == Empty {
				addMove(start_row, start_col, row, col, color, false, moves)
				continue
			}
			if color != target_color {
				addMove(start_row, start_col, row, col, color, true, moves) // Capture
			}
			break // Stop on encountering any piece
		}
	}

}

// kingMoves generates moves for the king
func kingMoves(row int, col int, color rune, boardState *BoardState, moves **Move) {
	kingOffsets := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // Top row
		{0, -1}, {0, 1}, // Middle row
		{1, -1}, {1, 0}, {1, 1}, // Bottom row
	}

	for _, offset := range kingOffsets {
		toRow, toCol := row+offset[0], col+offset[1]
		if !isValid(toRow, toCol) {
			continue
		}
		target_color, target_piece := getColorAndPiece(toRow, toCol, boardState.Board)
		if target_piece == Empty {
			addMove(row, col, toRow, toCol, color, false, moves)
		} else if color != target_color {
			addMove(row, col, toRow, toCol, color, true, moves)
		}
	}
}

func addMove(fromRow, fromCol, toRow, toCol int, color rune, capture bool, moves **Move) {
	newMove := &Move{
		from:    [2]int{fromRow, fromCol},
		to:      [2]int{toRow, toCol},
		color:   color,
		capture: capture,
	}
	if *moves == nil {
		*moves = newMove
	} else {
		current := *moves
		for current.next != nil {
			current = current.next
		}
		current.next = newMove
	}
}

// Converts a move string to a 'Move' struct.
func StringToMoveStruct(moveStr string, color rune) (Move, error) {
	// Ensure the input is valid
	if len(moveStr) != 5 || moveStr[2] != ' ' {
		return Move{}, errors.New("invalid move string format")
	}

	// Helper function to convert chess notation to indices
	chessToIndex := func(pos string) ([2]int, error) {
		if len(pos) != 2 {
			return [2]int{}, errors.New("invalid position format")
		}
		col := int(pos[0] - 'a')   // Convert column ('a' -> 0, ..., 'h' -> 7)
		row := 8 - int(pos[1]-'0') // Convert row ('1' -> 7, ..., '8' -> 0)
		if col < 0 || col > 7 || row < 0 || row > 7 {
			return [2]int{}, errors.New("position out of bounds")
		}
		return [2]int{row, col}, nil
	}

	// Parse positions
	from, err := chessToIndex(moveStr[:2])
	if err != nil {
		return Move{}, fmt.Errorf("invalid 'From' position: %w", err)
	}

	to, err := chessToIndex(moveStr[3:5])
	if err != nil {
		return Move{}, fmt.Errorf("invalid 'To' position: %w", err)
	}

	// Create the Move struct (Color and Capture need additional context to fill correctly)
	move := Move{
		from:    from,
		to:      to,
		color:   color,
		capture: false,
	}

	return move, nil
}

// validateMove checks whether a move is valid.
func ValidateMove(move *Move, boardState *BoardState) error {
	board := boardState.Board
	from_row, from_col := move.from[0], move.from[1]
	color := move.color
	if !isInBounds(move.from) || !isInBounds(move.to) {
		return errors.New("move out of bounds")
	}

	fromColor, _ := getColorAndPiece(move.from[0], move.from[1], board)
	if fromColor != color {
		return errors.New("the piece to be moved is not owned")
	}

	toColor, _ := getColorAndPiece(move.to[0], move.to[1], board)
	if toColor == color {
		return errors.New("the target position contains an owned piece")
	}

	pinned, err := isPinned(from_row, from_col, boardState)
	if err != nil {
		return err
	}
	if pinned {
		return errors.New("pinned piece can't be moved")
	}

	var moves *Move
	generateMovesForPiece(from_row, from_col, boardState, &moves)
	if !isMoveInMoves(move, moves) {
		return errors.New("move doesn't exist")
	}

	var tmpBoardState BoardState
	data, _ := json.Marshal(boardState)
	json.Unmarshal(data, &tmpBoardState)

	MakeMove(move, &tmpBoardState)

	attacked, err := kingAttacked(color, &tmpBoardState)
	if err != nil {
		return err
	}
	if attacked {
		return errors.New("king is under attack")
	}

	return nil
}

// Applies a specified move to the board.
func MakeMove(move *Move, boardState *BoardState) {
	fromRow, fromCol := move.from[0], move.from[1]
	toRow, toCol := move.to[0], move.to[1]
	fromColor, fromPiece := getColorAndPiece(fromRow, fromCol, boardState.Board)

	// Update king positions
	if fromPiece == 'x' {
		if fromColor == 'w' {
			boardState.WhiteKingMoved = true
			boardState.WhiteKingPos = [2]int{toRow, toCol}
		}
		if fromColor == 'b' {
			boardState.BlackKingMoved = true
			boardState.BlackKingPos = [2]int{toRow, toCol}
		}
	}

	// Update en passant
	boardState.EnPassant = [2]int{-1, -1}
	if fromPiece == 'p' {
		if math.Abs(float64(fromRow-toRow)) == 2 {
			boardState.EnPassant = [2]int{toRow, toCol}
		}
	}

	// Update board
	boardState.Board[toRow][toCol] = boardState.Board[fromRow][fromCol]
	boardState.Board[fromRow][fromCol] = Empty

	// Update player
	if fromColor == 'w' {
		boardState.TurnColor = 'b'
	} else {
		boardState.TurnColor = 'w'
	}
}
