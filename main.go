package main

import (
	"errors"
	"fmt"
)

const (
	n = iota
	o
	x
)

// Allows to play a single game of tic tac toe and then display the successive states

// Enhancements
// - Have the state contains each player move and not the turn result
// - Instead of saving multiple arrays in the state, only hold the couple { playerNumber, position }
// - Use a map ? something like  [0] -> { 1, 4 } | [1] -> { 2, 0 } | [3] -> { 1, 2 } ...
// 		- that would allow to get closer to the use of Merkle Trees ;
//		- unfortunately for these simple types, no performance gain in using the hash references (the contrary instead)

// Player corresponds to the three players : player 1 (O), player two (X),
// .... and the 'no player' (N)
// referenced by the iota constants (nox)
type Player int

func (e Player) String() string {
	return [...]string{" ", "O", "X"}[e]
}

// Game instance
type Game []Grid

func (s Game) String() string {
	str := "=== Game Summary ===\n"
	for _, grid := range s {
		if len(grid.grid) != 0 {
			str += fmt.Sprintf("\n%s\n", grid)
		}
	}
	return str + "\n===============\n"
}

// Grid is a 'state'
type Grid struct {
	turnNumber int
	grid       []Player
}

// Create initializes a game
func Create() Game {
	s := make(Game, 9)
	s[0] = Grid{
		turnNumber: 0,
		grid:       make([]Player, 9),
	}
	return s
}

// Play launches a game
func (s Game) Play() {

	winner := -1

	turnNumber := 0
	g := Grid{
		turnNumber: turnNumber,
		grid:       make([]Player, 9),
	}

	for {
		g.turnNumber++

		fmt.Println(g)

		if g.executePlayerTurn(1) {
			winner = 1
			break
		}

		// easiest dumb way to find the game has no winner
		if g.turnNumber == 5 {
			break
		}

		if g.executePlayerTurn(2) {
			winner = 2
			break
		}
		s = append(s, Grid{
			turnNumber: g.turnNumber,
			grid:       g.copyGrid(),
		})
	}

	fmt.Println(g)

	s = append(s, Grid{
		turnNumber: g.turnNumber,
		grid:       g.copyGrid(),
	})

	if winner == -1 {
		fmt.Println("\033[33mThere was no winner in that game :/\033[0m")
		fmt.Scanln()
	} else {
		fmt.Printf("\033[32mPlayer %d won !\033[0m", winner)
		fmt.Scanln()
	}

	fmt.Println(s)
}

func (g Grid) executePlayerTurn(playerNumber int) bool {
	var pos int
	for {
		fmt.Printf("Player %d, Enter a number to place a pawn (0-8): ", playerNumber)
		_, err := fmt.Scanln(&pos)
		if err != nil {
			fmt.Println("\033[31mPlease enter a valid number\033[0m")
			continue
		}
		win, err := g.placeSymbol(pos, Player(playerNumber))
		if err != nil {
			fmt.Println("\033[31m" + err.Error() + "\033[0m")
			continue
		}
		if win {
			return true
		}
		break
	}
	fmt.Println(g)
	return false
}

func (g Grid) placeSymbol(pos int, e Player) (bool, error) {
	if pos < 0 || pos > 8 {
		return false, errors.New("Enter a correct number (0-8)")
	}
	if g.grid[pos] != n {
		return false, errors.New("Place is already taken, please enter another one")
	}
	g.grid[pos] = e
	return g.checkWin(), nil
}

func (g Grid) checkWin() bool {
	if g.turnNumber < 3 {
		return false
	}

	return g.checkRow(0, 1, 2) || g.checkRow(3, 4, 5) || g.checkRow(6, 7, 8) ||
		g.checkRow(0, 3, 6) || g.checkRow(1, 4, 7) || g.checkRow(2, 5, 8) ||
		g.checkRow(0, 4, 8) || g.checkRow(2, 4, 6)
}

func (g Grid) checkRow(pos1, pos2, pos3 int) bool {
	return g.grid[pos1] == g.grid[pos2] &&
		g.grid[pos2] == g.grid[pos3] &&
		g.grid[pos1] != n
}

func (g Grid) copyGrid() []Player {
	result := make([]Player, len(g.grid))
	for i, v := range g.grid {
		result[i] = v
	}
	return result
}

func (g Grid) String() string {

	return fmt.Sprintf(`Turn number %d:
|----|----|----|
| %s  | %s  | %s  |
|----|----|----|
| %s  | %s  | %s  |
|----|----|----|
| %s  | %s  | %s  |
|----|----|----|
`,
		g.turnNumber,
		g.grid[0], g.grid[1], g.grid[2],
		g.grid[3], g.grid[4], g.grid[5],
		g.grid[6], g.grid[7], g.grid[8])
}

func main() {
	fmt.Println("Hello world")

	game := Create()
	game.Play()
}
