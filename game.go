package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const (
	EMPTY  = 0
	PLAYER = -1
	AI     = 1
	INF    = 2
)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Strategy struct{ move, payoff int } // struct

var clear map[string]func()

func evaluate(board [9]int) int {
	for i := 0; i < 3; i++ {
		if board[i] == board[i+3] && board[i] == board[i+6] {
			return board[i]
		}
	}
	for i := 0; i <= 6; i += 3 {
		if board[i] == board[i+1] && board[i] == board[i+2] {
			return board[i]
		}
	}
	if board[0] == board[4] && board[0] == board[8] {
		return board[0]
	}
	if board[2] == board[4] && board[2] == board[6] {
		return board[2]
	}
	return 0
}

func minimax(board [9]int, whoseTurn, numBlanks int) Strategy {
	evaluation := evaluate(board)
	if evaluation != 0 {
		result := Strategy{-1, evaluation}
		return result
	}
	if numBlanks == 0 {
		result := Strategy{-1, 0}
		return result
	}
	bestMove := -1
	var bestVal int
	if whoseTurn == AI {
		bestVal = -INF
	} else {
		bestVal = INF
	}
	for i := 0; i < 9; i++ {
		if board[i] == EMPTY {
			board[i] = whoseTurn
			branchVal := minimax(board, -whoseTurn, numBlanks-1).payoff
			if whoseTurn == AI {
				bestVal = max(branchVal, bestVal)
			} else {
				bestVal = min(branchVal, bestVal)
			}

			board[i] = EMPTY
			if bestVal == branchVal {
				bestMove = i
			}
		}
	}
	result := Strategy{bestMove, bestVal}
	return result
}

func getSymbol(r int) rune {
	switch r {
	case PLAYER:
		return 'X'
	case AI:
		return 'O'
	case EMPTY:
		return ' '
	}
	return '?'
}

func getPlayerMove(board *[9]int) {
	var move int
	fmt.Printf("Where would you like to move? [1 - 9] ")
	fmt.Scanf("%d", &move)
	fmt.Printf("\n")
	for board[move-1] != EMPTY || move < 1 || move > 9 {
		fmt.Printf("Invalid. ")
		fmt.Scanf("%d", &move)
		fmt.Printf("\n")
	}
	board[move-1] = PLAYER
}

func printBoard(board [9]int) {
	fmt.Printf("\n")
	fmt.Printf(" %c | %c | %c \n", getSymbol(board[0]), getSymbol(board[1]), getSymbol(board[2]))
	fmt.Printf("-----------\n")
	fmt.Printf(" %c | %c | %c \n", getSymbol(board[3]), getSymbol(board[4]), getSymbol(board[5]))
	fmt.Printf("-----------\n")
	fmt.Printf(" %c | %c | %c \n", getSymbol(board[6]), getSymbol(board[7]), getSymbol(board[8]))
	fmt.Printf("\n")
}
func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearScreen() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Unsupported platform!") // can't clear terminal
	}
}

func main() {
	board := [9]int{}
	firstMover, numBlanks := 0, 9

	clearScreen()
	fmt.Printf("Tic-tac-toe\n")
	fmt.Printf("\twith the minimax algorithm\n")
	fmt.Printf("Human: X\n")
	fmt.Printf("Machine: O\n")
	fmt.Printf("-----------------------------------\n")

	for i := true; i; i = (firstMover != 1 && firstMover != 2) {
		fmt.Printf("Would you like to go first or second? [1 / 2] ")
		fmt.Scanf("%d", &firstMover)
	}

	clearScreen()
	firstMover %= 2

	for i := 0; i < 9; i++ {
		if (firstMover+i)%2 == 1 {
			clearScreen()
			printBoard(board)
			getPlayerMove(&board)
		} else {
			aiMove := minimax(board, AI, numBlanks).move
			board[aiMove] = AI
		}

		if evaluate(board) != 0 {
			clearScreen()
			printBoard(board)
			if evaluate(board) == AI {
				fmt.Printf("AI win.\n")
			} else {
				fmt.Printf("HUMAN win.\n")
			}
			os.Exit(0)
		}
		numBlanks--
	}

	clearScreen()
	printBoard(board)
	fmt.Printf("Draw.\n")
	os.Exit(0)
}
