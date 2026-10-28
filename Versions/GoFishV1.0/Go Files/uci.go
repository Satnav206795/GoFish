package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	Depth int `json:"depth"`
}

func load_settings(path string) Settings {
	settings := Settings{Depth: 4}

	data, err := os.ReadFile(path)
	if err != nil {
		return settings
	}

	json.Unmarshal(data, &settings)
	return settings
}

func fresh_board() Board {
	return Board{
		board: [8][8]float64{
			{-4, -2, -3, -5, -6, -3, -2, -4},
			{-1, -1, -1, -1, -1, -1, -1, -1},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1},
			{4, 2, 3, 5, 6, 3, 2, 4},
		},
	}
}

func uci_loop(settings Settings) {
	board := fresh_board()
	is_white := true

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {

		case "uci":
			fmt.Println("id name GoFishV1.0")
			fmt.Println("id author Satnav")
			fmt.Println("uciok")

		case "isready":
			fmt.Println("readyok")

		case "ucinewgame":
			board = fresh_board()
			is_white = true

		case "position":
			if len(parts) < 2 {
				continue
			}
			board = fresh_board()
			is_white = true

			move_idx := -1
			for i, p := range parts {
				if p == "moves" {
					move_idx = i
					break
				}
			}

			if move_idx != -1 {
				for _, mv := range parts[move_idx+1:] {
					m := uci_to_move(mv, board, is_white)
					if m != nil {
						if len(m) == 5 {
							board.move_piece(m[0], m[1], m[2], m[3], m[4])
						} else {
							board.move_piece(m[0], m[1], m[2], m[3])
						}
						is_white = !is_white
					}
				}
			}

		case "go":
			depth := settings.Depth
			for i, p := range parts {
				if p == "depth" && i+1 < len(parts) {
					if d, err := strconv.Atoi(parts[i+1]); err == nil {
						depth = d
					}
				}
			}

			best := find_best_move(board, is_white, depth)
			if best != nil {
				fmt.Printf("bestmove %s\n", move_to_uci(best))
			} else {
				fmt.Println("bestmove 0000")
			}

		case "quit":
			return
		}
	}
}

func uci_to_move(s string, board Board, is_white bool) []int {
	if len(s) < 4 {
		return nil
	}

	fromx := int(s[0] - 'a')
	fromy := 8 - int(s[1]-'0')
	tox := int(s[2] - 'a')
	toy := 8 - int(s[3]-'0')

	// Promotion — 5th char is piece e.g. q/r/b/n
	if len(s) == 5 {
		piece := 5 // default queen
		switch s[4] {
		case 'q':
			piece = 5
		case 'r':
			piece = 4
		case 'b':
			piece = 3
		case 'n':
			piece = 2
		}
		if !is_white {
			piece = -piece
		}
		return []int{fromx, fromy, tox, toy, piece}
	}

	// Castling — king moves 2 squares
	if board.board[fromy][fromx] == 6 && fromx == 4 && fromy == 7 {
		if tox == 6 {
			return []int{11, 0, 0, 0} // W kingside
		}
		if tox == 2 {
			return []int{12, 0, 0, 0} // W queenside
		}
	}
	if board.board[fromy][fromx] == -6 && fromx == 4 && fromy == 0 {
		if tox == 6 {
			return []int{13, 0, 0, 0} // B kingside
		}
		if tox == 2 {
			return []int{14, 0, 0, 0} // B queenside
		}
	}

	// En passant — pawn moves diagonally to empty square
	if board.board[fromy][fromx] == 1 && fromx != tox && board.board[toy][tox] == 0 {
		if tox == fromx-1 {
			return []int{15, fromx, 0, 0} // W en passant left
		}
		return []int{16, fromx, 0, 0} // W en passant right
	}
	if board.board[fromy][fromx] == -1 && fromx != tox && board.board[toy][tox] == 0 {
		if tox == fromx-1 {
			return []int{17, fromx, 0, 0} // B en passant left
		}
		return []int{18, fromx, 0, 0} // B en passant right
	}

	return []int{fromx, fromy, tox, toy}
}

func move_to_uci(m []int) string {
	switch m[0] {
	case 11:
		return "e1g1"
	case 12:
		return "e1c1"
	case 13:
		return "e8g8"
	case 14:
		return "e8c8"
	case 15:
		return string(rune('a'+m[1])) + "5" + string(rune('a'+m[1]-1)) + "6"
	case 16:
		return string(rune('a'+m[1])) + "5" + string(rune('a'+m[1]+1)) + "6"
	case 17:
		return string(rune('a'+m[1])) + "4" + string(rune('a'+m[1]-1)) + "3"
	case 18:
		return string(rune('a'+m[1])) + "4" + string(rune('a'+m[1]+1)) + "3"
	}

	// Promotion
	if len(m) == 5 {
		from := string(rune('a'+m[0])) + strconv.Itoa(8-m[1])
		to := string(rune('a'+m[2])) + strconv.Itoa(8-m[3])
		pieces := map[int]string{5: "q", -5: "q", 4: "r", -4: "r", 3: "b", -3: "b", 2: "n", -2: "n"}
		return from + to + pieces[m[4]]
	}

	from := string(rune('a'+m[0])) + strconv.Itoa(8-m[1])
	to := string(rune('a'+m[2])) + strconv.Itoa(8-m[3])
	return from + to
}
