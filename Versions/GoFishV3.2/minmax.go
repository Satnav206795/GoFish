package main

func minmax(board Board, is_white bool, depth int, alpha int, beta int) int {

	if board.check_and_update_repition() {
		return 0 //Draw by repition
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

	if depth == 0 {
		return eval(board)
	}

	if is_white {
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

			score := minmax(newBoard, !is_white, depth-1, alpha, beta)

			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				break
			}

		}
		return alpha
	} else {
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
			score := minmax(newBoard, !is_white, depth-1, alpha, beta)

			if score < beta {
				beta = score
			}
			if alpha >= beta {
				break
			}

		}
		return beta
	}
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
			score := minmax(newBoard, !is_white, depth-1, -999999999, 999999999)
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
			score := minmax(newBoard, !is_white, depth-1, -999999999, 999999999)
			if score < bestScore {
				bestScore = score
				bestMove = m
			}
		}
	}

	return bestMove
}
