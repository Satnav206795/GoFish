package main

import (
	"math"
)

func eval(board Board) int {
	var score int = 0
	temp_score := 0

	w_bishop_counter := 0
	b_bishop_counter := 0

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			piece := board.board[y][x]

			if piece == 0 {
				continue
			}

			is_white := false
			if piece > 0 {
				is_white = true
			}

			abs_piece := int(math.Abs(piece))
			temp_score = 0

			switch abs_piece {
			case 1: // Pawn
				temp_score += 100
				if is_white {
					temp_score += 10 * (6 - y)
				} else {
					temp_score += 10 * (y - 1)
				}
				if x > 2 && x < 5 && y > 2 && y < 5 {
					temp_score += 30
				}

			case 2: // Knight
				temp_score += 300
				if x >= 2 && x <= 5 && y >= 2 && y <= 5 {
					temp_score += 50
				} else if x == 0 || x == 7 {
					temp_score -= 30
				}

			case 3: // Bishop
				temp_score += 300
				if is_white {
					w_bishop_counter++
				} else {
					b_bishop_counter++
				}

			case 4: // Rook
				temp_score += 500

			case 5: // Queen
				temp_score += 900

			case 6: // King
				temp_score += 100000
			}

			if is_white {
				score += temp_score
			} else {
				score -= temp_score
			}
		}
	}

	if w_bishop_counter >= 2 {
		score += 50
	}
	if b_bishop_counter >= 2 {
		score -= 50
	}

	return score
}
