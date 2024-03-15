package main

import (
	"fmt"
	"math/rand"
	"runtime"
)

var DIGITS = [9]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

var threadCount = runtime.GOMAXPROCS(0)

type Solver struct {
	id                                  int
	digits                              []byte
	board, origin                       [][]byte
	rowChecker, colChecker, gridChecker [9][9]bool
}

func NewSolver(board [][]byte, id int) *Solver {
	solver := &Solver{
		id:     id,
		board:  make([][]byte, len(board)),
		origin: board,
		digits: make([]byte, 9),
	}

	copy(solver.digits, DIGITS[:])
	rand.Shuffle(len(solver.digits), func(i, j int) {
		solver.digits[i], solver.digits[j] = solver.digits[j], solver.digits[i]
	})

	fmt.Printf("Random shuffle: %v\n", solver.digits)
	for r, row := range board {
		solver.board[r] = make([]byte, len(board))
		for c, val := range row {
			if val != '.' {
				solver.fill(r, c, val)
			}
		}
	}
	return solver
}

func (solver *Solver) canFill(r, c int, digit byte) bool {
	return !solver.checkInRow(r, digit) && !solver.checkInCol(c, digit) && !solver.checkInGrid(r, c, digit)
}

func (solver *Solver) checkInRow(r int, digit byte) bool {
	return solver.rowChecker[r][digit-'1']
}

func (solver *Solver) checkInCol(c int, digit byte) bool {
	return solver.colChecker[c][digit-'1']
}
func (solver *Solver) checkInGrid(r, c int, digit byte) bool {
	return solver.gridChecker[solver.hashGrid(r, c)][digit-'1']
}

func (solver *Solver) hashGrid(r, c int) int {
	return r/3*3 + c/3
}

func (solver *Solver) fill(r, c int, digit byte) {
	solver.board[r][c] = digit
	solver.rowChecker[r][digit-'1'] = true
	solver.colChecker[c][digit-'1'] = true
	solver.gridChecker[solver.hashGrid(r, c)][digit-'1'] = true
}

func (solver *Solver) unfill(r, c int, digit byte) {
	solver.board[r][c] = '.'
	solver.rowChecker[r][digit-'1'] = false
	solver.colChecker[c][digit-'1'] = false
	solver.gridChecker[solver.hashGrid(r, c)][digit-'1'] = false
}

func (solver *Solver) solve(ansChan chan *Solver, done <-chan struct{}) {
	fmt.Printf("Start solver %v\n", solver.id)
	var recurse func(r, c int) bool
	recurse = func(r, c int) bool {
		select {
		case _, ok := <-done:
			if !ok {
				fmt.Printf("Shutdown solver %v\n", solver.id)
				return true
			}
		default:
		}
		if c >= len(solver.board[0]) {
			return recurse(r+1, 0)
		}

		if r >= len(solver.board) {
			ansChan <- solver
			return true
		}

		if solver.origin[r][c] != '.' {
			return recurse(r, c+1)
		}

		for _, digit := range DIGITS {
			// fmt.Println(solver.canFill(r, c, digit), r, c, digit)
			if solver.canFill(r, c, digit) {
				solver.fill(r, c, digit)
				if recurse(r, c+1) {
					return true
				}
				solver.unfill(r, c, digit)
			}
		}

		return false
	}

	recurse(0, 0)
	fmt.Printf("End solver %v\n", solver.id)
}

func solveSudoku(board [][]byte) {
	done := make(chan struct{})
	ansChan := make(chan *Solver)

	for i := 0; i < threadCount; i++ {
		solver := NewSolver(board, i)
		go solver.solve(ansChan, done)
	}
	solved := <-ansChan
	close(done)
	for r, row := range solved.board {
		copy(board[r], row)
	}
}

func main() {
	input := []string{
		"9..8.....",
		"......5..",
		".........",
		".2..1...3",
		".1.....6.",
		"...4...7.",
		"7.86.....",
		"....3.1..",
		"4.....2..",
	}
	// input := []string{
	// 	"53..7....",
	// 	"6..195...",
	// 	".98....6.",
	// 	"8...6...3",
	// 	"4..8.3..1",
	// 	"7...2...6",
	// 	".6....28.",
	// 	"...419..5",
	// 	"....8..79",
	// }
	board := make([][]byte, len(input))
	for i, row := range input {
		board[i] = []byte(row)
	}

	solveSudoku(board)
	for _, row := range board {
		fmt.Println(string(row))
	}
}
