package main

import "math"

func find_all_semi_legal_moves(board Board, is_white bool, attaacking_squares bool) [][]int {
	legal_moves := [][]int{}

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			piece := board.board[y][x]

			if piece == 0 {
				continue
			}

			if piece > 0 && is_white == false {
				continue
			}

			if piece < 0 && is_white == true {
				continue
			}

			piece_type := math.Abs(piece)
			temp_move := [][]int{}

			switch piece_type {

			case 1: // Pawn
				if attaacking_squares {
					if is_white {
						for _, dx := range []int{-1, 1} {
							nx := x + dx
							ny := y - 1

							if inbounds(nx, ny) {
								temp_move = append(temp_move, []int{x, y, nx, ny})
							}
						}
					} else {
						for _, dx := range []int{-1, 1} {
							nx := x + dx
							ny := y + 1

							if inbounds(nx, ny) {
								temp_move = append(temp_move, []int{x, y, nx, ny})
							}
						}
					}
					continue
				}
				if is_white {

					// Move forward 1
					if inbounds(x, y-1) && board.board[y-1][x] == 0 {
						temp_move = append(temp_move, []int{x, y, x, y - 1})

						// Move forward 2
						if y == 6 && inbounds(x, y-2) && board.board[y-2][x] == 0 {
							temp_move = append(temp_move, []int{x, y, x, y - 2})
						} else if y == 1 {
							for p := range 4 {
								temp_move = append(temp_move, []int{x, y, x, y - 1, p + 2})
							}
						}

					}

					// Diagonal captures
					for _, dx := range []int{-1, 1} {
						nx := x + dx
						ny := y - 1

						if inbounds(nx, ny) && board.board[ny][nx] < 0 {
							temp_move = append(temp_move, []int{x, y, nx, ny})
							if y == 1 {
								for p := range 4 {
									temp_move = append(temp_move, []int{x, y, nx, ny, p + 2})
								}
							}
						}
					}

					// En Passant
					for _, dx := range []int{-1, 1} {
						if inbounds(x+dx, 3) && y == 3 {
							if board.board[3][x+dx] == -1 && (board.b_pawns_moved_2>>uint(x+dx))&1 == 1 {
								code := 15
								if dx == 1 {
									code = 16
								}
								temp_move = append(temp_move, []int{code, x, 0, 0})
							}
						}
					}

				} else {

					// Move forward 1
					if inbounds(x, y+1) && board.board[y+1][x] == 0 {
						temp_move = append(temp_move, []int{x, y, x, y + 1})

						// Move forward 2
						if y == 1 && inbounds(x, y+2) && board.board[y+2][x] == 0 {
							temp_move = append(temp_move, []int{x, y, x, y + 2})
						} else if y == 6 {
							for p := range 4 {
								temp_move = append(temp_move, []int{x, y, x, y + 1, p + 2})
							}
						}
					}

					// Diagonal captures
					for _, dx := range []int{-1, 1} {
						nx := x + dx
						ny := y + 1

						if inbounds(nx, ny) && board.board[ny][nx] > 0 {
							temp_move = append(temp_move, []int{x, y, nx, ny})
							if y == 6 {
								for p := range 4 {
									temp_move = append(temp_move, []int{x, y, nx, ny, p + 2})
								}
							}
						}
					}

					// En Passant
					for _, dx := range []int{-1, 1} {
						if inbounds(x+dx, 4) && y == 4 {
							if board.board[4][x+dx] == 1 && (board.w_pawns_moved_2>>uint(x+dx))&1 == 1 {
								code := 17
								if dx == 1 {
									code = 18
								}
								temp_move = append(temp_move, []int{code, x, 0, 0})
							}
						}
					}
				}

			case 2: // Knight
				directions := [][2]int{
					{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
					{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
				}

				for _, d := range directions {
					mx, my := d[0], d[1]
					if is_legal(x+mx, y+my, board, is_white) {
						temp_move = append(temp_move, []int{x, y, x + mx, y + my})
					}
				}

			case 3: // Bishop
				d_moves := generate_diagonals(board, x, y, is_white)
				temp_move = append(temp_move, d_moves...)

			case 4: // Rook
				s_moves := generate_straights(board, x, y, is_white)
				temp_move = append(temp_move, s_moves...)

			case 5: // Queen
				d_moves := generate_diagonals(board, x, y, is_white)
				temp_move = append(temp_move, d_moves...)

				s_moves := generate_straights(board, x, y, is_white)
				temp_move = append(temp_move, s_moves...)

			case 6: // King
				directions := [][2]int{
					{0, 1}, {0, -1}, {1, 0}, {-1, 0},
					{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
				}

				for _, d := range directions {
					if is_legal(x+d[0], y+d[1], board, is_white) {
						temp_move = append(temp_move, []int{x, y, x + d[0], y + d[1]})
					}
				}

				if is_white && !board.white_king_moved {
					//Short Castling
					if board.white_h1_rook_moved == false && board.board[7][6] == 0 && board.board[7][5] == 0 {
						if !board.is_black_controlled(x, y) && !board.is_black_controlled(5, 7) && !board.is_black_controlled(6, 7) {
							temp_move = append(temp_move, []int{11, 0, 0, 0})
						}
					}
					if board.white_a1_rook_moved == false && board.board[7][1] == 0 && board.board[7][2] == 0 && board.board[7][3] == 0 { //Long Castling
						if !board.is_black_controlled(x, y) && !board.is_black_controlled(1, 7) && !board.is_black_controlled(2, 7) && !board.is_black_controlled(3, 7) {
							temp_move = append(temp_move, []int{12, 0, 0, 0})
						}
					}
				} else if !is_white && !board.black_king_moved {
					//Short Castling
					if board.black_h8_rook_moved == false && board.board[0][6] == 0 && board.board[0][5] == 0 {
						if !board.is_white_controlled(x, y) && !board.is_white_controlled(5, 0) && !board.is_white_controlled(6, 0) {
							temp_move = append(temp_move, []int{13, 0, 0, 0})
						}
					}
					if board.black_a8_rook_moved == false && board.board[0][1] == 0 && board.board[0][2] == 0 && board.board[0][3] == 0 { //Long Castling
						if !board.is_white_controlled(x, y) && !board.is_white_controlled(1, 0) && !board.is_white_controlled(2, 0) && !board.is_white_controlled(3, 0) {
							temp_move = append(temp_move, []int{14, 0, 0, 0})
						}
					}
				}
			}

			legal_moves = append(legal_moves, temp_move...)
		}
	}

	return legal_moves
}

func find_all_legal_moves(board Board, is_white bool) [][]int {

	moves := [][]int{}

	board.update_w_check_bitmap()
	board.update_b_check_bitmap()
	legal_moves := find_all_semi_legal_moves(board, is_white, false)

	king_x, king_y := find_king_square(board, is_white)

	if len(legal_moves) == 0 {
		return moves
	}

	for _, m := range legal_moves {

		// Guard against any empty move slice
		if len(m) == 0 {
			continue
		}

		if m[0] >= 11 {
			moves = append(moves, m)
			continue
		}

		check_x, check_y := king_x, king_y
		if m[0] == king_x && m[1] == king_y {
			check_x, check_y = m[2], m[3]
		}

		newBoard := board
		if len(m) == 5 {
			newBoard.move_piece(m[0], m[1], m[2], m[3], m[4])
		} else {
			newBoard.move_piece(m[0], m[1], m[2], m[3])
		}

		if is_white {

			newBoard.update_b_check_bitmap()
			if !newBoard.is_black_controlled(check_x, check_y) {
				moves = append(moves, m)
			}

		} else {

			newBoard.update_w_check_bitmap()
			if !newBoard.is_white_controlled(check_x, check_y) {
				moves = append(moves, m)
			}

		}
	}

	return moves
}

func is_legal(x int, y int, board Board, is_white bool) bool {
	if inbounds(x, y) {
		if is_white {
			if board.board[y][x] <= 0 {
				return true
			}
		} else {
			if board.board[y][x] >= 0 {
				return true
			}
		}
	}
	return false
}

func inbounds(x int, y int) bool {
	if x >= 0 && x <= 7 && y >= 0 && y <= 7 {
		return true
	} else {
		return false
	}
}

func generate_diagonals(board Board, x int, y int, is_white bool) [][]int {
	temp_move := [][]int{}

	directions := [][2]int{
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}

	for _, d := range directions {
		mx, my := d[0], d[1]
		counter := 1

		for {
			nx := x + mx*counter
			ny := y + my*counter

			if !inbounds(nx, ny) {
				break
			}

			target := board.board[ny][nx]

			if is_white {
				if target > 0 {
					break
				}
				temp_move = append(temp_move, []int{x, y, nx, ny})
				if target < 0 {
					break
				}
			} else {
				if target < 0 {
					break
				}
				temp_move = append(temp_move, []int{x, y, nx, ny})
				if target > 0 {
					break
				}
			}

			counter++
		}
	}
	return temp_move
}

func generate_straights(board Board, x int, y int, is_white bool) [][]int {
	temp_move := [][]int{}

	directions := [][2]int{
		{0, 1}, {0, -1}, {1, 0}, {-1, 0},
	}

	for _, d := range directions {
		mx, my := d[0], d[1]
		counter := 1

		for {
			nx := x + mx*counter
			ny := y + my*counter

			if !inbounds(nx, ny) {
				break
			}

			target := board.board[ny][nx]

			if is_white {
				if target > 0 {
					break
				}
				temp_move = append(temp_move, []int{x, y, nx, ny})
				if target < 0 {
					break
				}
			} else {
				if target < 0 {
					break
				}
				temp_move = append(temp_move, []int{x, y, nx, ny})
				if target > 0 {
					break
				}
			}

			counter++
		}
	}
	return temp_move
}

func find_king_square(board Board, is_white bool) (int, int) {
	var piece float64 = 0

	if is_white {
		piece = 6
	} else {
		piece = -6
	}

	for y := range 8 {
		for x := range 8 {
			if board.board[y][x] == piece {
				return x, y
			}
		}
	}
	return -1, -1
}

/*
func check_special_moves(move []int) (bool, int) {

	switch move[0] {

	case 11: //W Castling Short
		return true, 11

	case 12: //W Castling Long
		return true, 12

	case 13: //B Castling Short
		return true, 13

	case 14: //B Castling Long
		return true, 14

	case 15: //W En Passant Left
		return true, 15 //Use fromy as X of pawn

	case 16: //W En Passant Right
		return true, 16

	case 17: //B En Passant Left
		return true, 17

	case 18: //B En Passant Right
		return true, 18

	}

	return false, 0
}
*/
