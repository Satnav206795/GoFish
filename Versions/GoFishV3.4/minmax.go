package main

func minmax(board Board, is_white bool, depth int, alpha int, beta int) int {

	if board.check_and_update_repition() {
		return 0 //Draw by repition
	}

	kx, ky := find_king_square(board, is_white)
	if king_in_check(board, is_white, kx, ky) && depth == 0 {
		depth += 1
	}

	if depth == 0 {
		return eval(board)
	}

	moves := find_all_legal_moves(board, is_white)

	if len(moves) == 0 {
		if king_in_check(board, is_white, kx, ky) {
			if is_white {
				return -999999999
			}
			return 999999999
		}
		return 0
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

func king_in_check(board Board, is_white bool, king_x int, king_y int) bool {
	if is_white {
		_, captures := generate_diagonals(board, king_x, king_y, true)
		for _, m := range captures {
			piece := board.board[m[3]][m[2]]
			if piece == -3 || piece == -5 {
				return true
			}
		}

		_, captures = generate_straights(board, king_x, king_y, true)
		for _, m := range captures {
			piece := board.board[m[3]][m[2]]
			if piece == -4 || piece == -5 {
				return true
			}
		}

		directions := [][2]int{
			{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
			{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
		}
		for _, d := range directions {
			mx, my := d[0]+king_x, d[1]+king_y
			if inbounds(mx, my) && board.board[my][mx] == -2 {
				return true
			}
		}
		return false

	} else {
		_, captures := generate_diagonals(board, king_x, king_y, false)
		for _, m := range captures {
			piece := board.board[m[3]][m[2]]
			if piece == 3 || piece == 5 {
				return true
			}
		}

		_, captures = generate_straights(board, king_x, king_y, false)
		for _, m := range captures {
			piece := board.board[m[3]][m[2]]
			if piece == 4 || piece == 5 {
				return true
			}
		}

		directions := [][2]int{
			{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
			{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
		}
		for _, d := range directions {
			mx, my := d[0]+king_x, d[1]+king_y
			if inbounds(mx, my) && board.board[my][mx] == 2 {
				return true
			}
		}
		return false
	}
}
