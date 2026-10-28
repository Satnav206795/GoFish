package main

import (
	"hash/fnv"
	"slices"
	"unsafe"
)

type Board struct {
	board                  [8][8]float64
	white_checking_squares int64 //Bitmap
	black_checking_squares int64 // Bitmap
	white_king_moved       bool
	black_king_moved       bool
	white_a1_rook_moved    bool
	white_h1_rook_moved    bool
	black_a8_rook_moved    bool
	black_h8_rook_moved    bool
	w_pawns_moved_2        int8
	b_pawns_moved_2        int8
	black_turn             bool
	repition               [500]uint64
	historyCount           int
}

func (b *Board) inverse() {

	for s := 0; s < 8; s++ {
		slices.Reverse(b.board[s][:])
	}

	slices.Reverse(b.board[:])
}

func (b *Board) move_piece(fromx int, fromy int, tox int, toy int, promoting_piece ...int) {

	if b.black_turn {
		b.b_pawns_moved_2 = 0
	} else {
		b.w_pawns_moved_2 = 0
	}

	if fromx >= 11 { //Special Moves
		switch fromx {

		case 11:
			b.board[7][7] = 0 //Rook Gone
			b.board[7][4] = 0 //King Gone
			b.board[7][6] = 6 //King Back
			b.board[7][5] = 4 //Rook Back
			b.white_king_moved = true

		case 12:
			b.board[7][0] = 0 //Rook Gone
			b.board[7][4] = 0 //King Gone
			b.board[7][2] = 6 //King Back
			b.board[7][3] = 4 //Rook Back
			b.white_king_moved = true

		case 13:
			b.board[0][7] = 0 //Rook Gone
			b.board[0][4] = 0 //King Gone
			b.board[0][6] = 6 //King Back
			b.board[0][5] = 4 //Rook Back
			b.black_king_moved = true

		case 14:
			b.board[0][0] = 0 //Rook Gone
			b.board[0][4] = 0 //King Gone
			b.board[0][2] = 6 //King Back
			b.board[0][3] = 4 //Rook Back
			b.black_king_moved = true

		case 15:
			b.board[3][fromy] = 0 //Fromy = the x of pawn
			b.board[2][fromy-1] = 1
			b.board[3][fromy-1] = 0

		case 16:
			b.board[3][fromy] = 0 //Fromy = the x of pawn
			b.board[2][fromy+1] = 1
			b.board[3][fromy+1] = 0

		case 17:
			b.board[4][fromy] = 0 //Fromy = the x of pawn
			b.board[5][fromy-1] = -1
			b.board[4][fromy-1] = 0

		case 18:
			b.board[4][fromy] = 0 //Fromy = the x of pawn
			b.board[5][fromy+1] = -1
			b.board[4][fromy+1] = 0

		}

		return
	}

	piece := b.board[fromy][fromx]

	// Default to queen promotion; use provided piece if given.
	// For black the promo value must be negative.
	promo := 5
	if piece == -1 {
		promo = -5
	}
	if len(promoting_piece) > 0 {
		promo = promoting_piece[0]
	}

	if piece == 1 && toy == 0 {
		b.board[fromy][fromx] = 0
		b.board[toy][tox] = float64(promo)

	} else if piece == -1 && toy == 7 {
		b.board[fromy][fromx] = 0
		// Ensure the promoted piece is negative for black
		if promo > 0 {
			promo = -promo
		}
		b.board[toy][tox] = float64(promo)

	} else {
		b.board[fromy][fromx] = 0
		b.board[toy][tox] = piece
	}

	if piece == 6 {
		b.white_king_moved = true
	}
	if piece == -6 {
		b.black_king_moved = true
	}
	if fromx == 0 && fromy == 7 {
		b.white_a1_rook_moved = true
	}
	if fromx == 7 && fromy == 7 {
		b.white_h1_rook_moved = true
	}
	if fromx == 0 && fromy == 0 {
		b.black_a8_rook_moved = true
	}
	if fromx == 7 && fromy == 0 {
		b.black_h8_rook_moved = true
	}
	if piece == 1 && fromy-toy == 2 && tox >= 0 && tox <= 7 {
		b.w_pawns_moved_2 |= int8(1) << uint(tox)
	}
	if piece == -1 && toy-fromy == 2 && tox >= 0 && tox <= 7 {
		b.b_pawns_moved_2 |= int8(1) << uint(tox)
	}
}

func (b *Board) update_w_check_bitmap() {
	b.white_checking_squares = 0
	moves := find_all_semi_legal_moves(*b, true, true)

	x := 0
	y := 0
	bit := 0
	for _, m := range moves {
		x = m[2] //Move to x
		y = m[3] //Move to y

		bit = y*8 + x
		b.white_checking_squares |= int64(1) << uint(bit)
	}
}

func (b *Board) update_b_check_bitmap() {
	b.black_checking_squares = 0
	moves := find_all_semi_legal_moves(*b, false, true)

	x := 0
	y := 0
	bit := 0
	for _, m := range moves {
		x = m[2] //Move to x
		y = m[3] //Move to y

		bit = y*8 + x
		b.black_checking_squares |= int64(1) << uint(bit)
	}
}

func (b *Board) is_white_controlled(x int, y int) bool {
	bit := y*8 + x
	return (b.white_checking_squares>>uint(bit))&1 == 1
}

func (b *Board) is_black_controlled(x int, y int) bool {
	bit := y*8 + x
	return (b.black_checking_squares>>uint(bit))&1 == 1
}

func (b *Board) Hash() uint64 {
	// 1. Find exactly where the map starts in memory.
	// This tells us the exact size (in bytes) of all preceding fields.
	stateBytesSize := unsafe.Offsetof(b.repition)

	// 2. Cast the struct pointer into a raw byte slice of that exact size.
	// This reads the board array, bitmasks, and flags as raw data.
	stateBytes := unsafe.Slice((*byte)(unsafe.Pointer(b)), stateBytesSize)

	// 3. Compute a fast, non-cryptographic FNV-1a hash of the slice.
	hasher := fnv.New64a()
	hasher.Write(stateBytes)
	return hasher.Sum64()
}

func (b *Board) check_and_update_repition() bool {
	BHash := b.Hash()
	occurrences := 0

	// 1. Scan backward through your repition array
	for i := b.historyCount - 1; i >= 0; i-- {
		if b.repition[i] == BHash {
			occurrences++

			// If it has been seen 2 times before, seeing it right now
			// makes it the 3rd time (Threefold Repetition Draw)
			if occurrences >= 2 {
				return true
			}
		}
	}

	// 2. Out-of-bounds safety catch to protect the [500] fixed-size array boundary
	if b.historyCount >= 500 {
		return true
	}

	// 3. Update the array with the current position hash and advance the counter
	b.repition[b.historyCount] = BHash
	b.historyCount++

	return false
}
