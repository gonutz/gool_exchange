package main

import "github.com/gonutz/prototype/draw"

const empty = 0
const player1 = 1
const player2 = 2

func main() {
	var columnSet [6][7]int
	nextPlayer := player1

	hasPlayerWon := func(player int) bool {
		for y := 0; y <= 5; y++ {
			for x := 0; x <= 3; x++ {
				if columnSet[y][0+x] == player &&
					columnSet[y][1+x] == player &&
					columnSet[y][2+x] == player &&
					columnSet[y][3+x] == player {
					return true
				}
			}
		}

		for x := 0; x <= 6; x++ {
			for y := 0; y <= 2; y++ {
				if columnSet[0+y][x] == player &&
					columnSet[1+y][x] == player &&
					columnSet[2+y][x] == player &&
					columnSet[3+y][x] == player {
					return true
				}
			}
		}

		return false
	}

	nextEmptyRow := func(column int) int {
		for y := 5; y >= 0; y-- {
			if columnSet[y][column] == empty {
				return y
			}
		}
		return -1
	}

    draw.RunWindow("4 gewinnt", 700, 600, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		mouseX, _ := window.MousePosition()
		mouseColumn := mouseX / 100

		leftDown := false
		for _, click := range window.Clicks() {
			if click.Button == draw.LeftButton {
				leftDown = true
			}
		}

		if leftDown {
			y := nextEmptyRow(mouseColumn)
			if y != -1 {
				if nextPlayer == player1 {
					columnSet[y][mouseColumn] = player1
					nextPlayer = player2
				} else {
					columnSet[y][mouseColumn] = player2
					nextPlayer = player1
				}
			}

			if hasPlayerWon(player1) || hasPlayerWon(player2) {
				window.Close()
			}
		}

		// Draw the game.
		window.FillRect(0, 0, 700, 600, draw.White)
        for x := 0; x < 7; x++ {
			for y := 0; y < 6; y++ {
            	window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.Black)

				if columnSet[y][x] == player1 {
		            window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGB(0.33, 0.33, 1))
				}
				if columnSet[y][x] == player2 {
		            window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGB(1, 0.5, 0.38))
				}

				if x == mouseColumn && y == nextEmptyRow(mouseColumn) {
		            window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGB(0.5, 0.86, 0.4))
				}
        	}
        }
    })
}
