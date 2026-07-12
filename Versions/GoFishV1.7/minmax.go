package main

func minmax(board Board, is_white bool, depth int) int {
	if depth == 0 {
		return eval(board)
	}

	moves := find_all_legal_moves(board, is_white)

	if len(moves) == 0 {
		board.update_w_check_bitmap()
		board.update_b_check_bitmap()
		kx, ky := find_king_square(board, is_white)
		if is_white && board.is_black_controlled(kx, ky) {
			return -999999999
		} else if !is_white && board.is_white_controlled(kx, ky) {
			return 999999999
		}
		return 0
	}

	best := 0

	if is_white {
		best = -999999999
		for _, m := range moves {
			newBoard := board

			if len(m) == 0 {
				continue
			}

			if len(m) >= 4 {
				if len(m) == 5 {
					newBoard.move_piece(m[0], m[1], m[2], m[3], m[4])
				} else {
					newBoard.move_piece(m[0], m[1], m[2], m[3])
				}
			} else {
				newBoard.move_piece(m[0], m[1], 0, 0)
			}
			score := minmax(newBoard, !is_white, depth-1)
			if score > best {
				best = score
			}
		}
	} else {
		best = 999999999
		for _, m := range moves {
			newBoard := board

			if len(m) == 0 {
				continue
			}

			if len(m) >= 4 {
				if len(m) == 5 {
					newBoard.move_piece(m[0], m[1], m[2], m[3], m[4])
				} else {
					newBoard.move_piece(m[0], m[1], m[2], m[3])
				}
			} else {
				newBoard.move_piece(m[0], m[1], 0, 0)
			}
			score := minmax(newBoard, !is_white, depth-1)
			if score < best {
				best = score
			}
		}
	}

	return best
}

func find_best_move(board Board, is_white bool, depth int) []int {
	moves := find_all_legal_moves(board, is_white)

	if len(moves) == 0 {
		return nil
	}

	bestMove := moves[0]

	if is_white {
		bestScore := -999999999
		for _, m := range moves {
			newBoard := board

			if len(m) == 0 {
				continue
			}

			if len(m) >= 4 {
				if len(m) == 5 {
					newBoard.move_piece(m[0], m[1], m[2], m[3], m[4])
				} else {
					newBoard.move_piece(m[0], m[1], m[2], m[3])
				}
			} else {
				newBoard.move_piece(m[0], m[1], 0, 0)
			}
			score := minmax(newBoard, !is_white, depth-1)
			if score > bestScore {
				bestScore = score
				bestMove = m
			}
		}
	} else {
		bestScore := 999999999
		for _, m := range moves {
			newBoard := board

			if len(m) == 0 {
				continue
			}

			if len(m) >= 4 {
				if len(m) == 5 {
					newBoard.move_piece(m[0], m[1], m[2], m[3], m[4])
				} else {
					newBoard.move_piece(m[0], m[1], m[2], m[3])
				}
			} else {
				newBoard.move_piece(m[0], m[1], 0, 0)
			}
			score := minmax(newBoard, !is_white, depth-1)
			if score < bestScore {
				bestScore = score
				bestMove = m
			}
		}
	}

	return bestMove
}
