package main

import "github.com/gonutz/prototype/draw"

const empty = 0
const player1 = 1
const player2 = 2

func main() {
	var columnSet [6][7]int
	var blinking [6][7]bool
	nextPlayer := player1
	moves := 0
	gameIsOver := false

	hasPlayerWon := func(player int) bool {
		won := false

		for y := 0; y <= 5; y++ {
			for x := 0; x <= 3; x++ {
				if columnSet[y][0+x] == player &&
					columnSet[y][1+x] == player &&
					columnSet[y][2+x] == player &&
					columnSet[y][3+x] == player {

					blinking[y][0+x] = true
					blinking[y][1+x] = true
					blinking[y][2+x] = true
					blinking[y][3+x] = true
					won = true
				}
			}
		}

		for x := 0; x <= 6; x++ {
			for y := 0; y <= 2; y++ {
				if columnSet[0+y][x] == player &&
					columnSet[1+y][x] == player &&
					columnSet[2+y][x] == player &&
					columnSet[3+y][x] == player {

					blinking[0+y][x] = true
					blinking[1+y][x] = true
					blinking[2+y][x] = true
					blinking[3+y][x] = true
					won = true

				}
			}
		}

		for x := 0; x <= 3; x++ {
			for y := 0; y <= 2; y++ {
				if columnSet[0+y][0+x] == player &&
					columnSet[1+y][1+x] == player &&
					columnSet[2+y][2+x] == player &&
					columnSet[3+y][3+x] == player {

					blinking[0+y][0+x] = true
					blinking[1+y][1+x] = true
					blinking[2+y][2+x] = true
					blinking[3+y][3+x] = true
					won = true

				}

				if columnSet[0+y][6-x] == player &&
					columnSet[1+y][5-x] == player &&
					columnSet[2+y][4-x] == player &&
					columnSet[3+y][3-x] == player {

					blinking[0+y][6-x] = true
					blinking[1+y][5-x] = true
					blinking[2+y][4-x] = true
					blinking[3+y][3-x] = true
					won = true

				}
			}
		}

		return won
	}

	isDraw := func() bool {
		return moves == 42
	}

	nextEmptyRow := func(column int) int {
		for y := 5; y >= 0; y-- {
			if columnSet[y][column] == empty {
				return y
			}
		}
		return -1
	}

	time := 0
	blinkersVisible := true
	mouseColumn := 3
	lastMouseX := -1

	draw.RunWindow("4 gewinnt - F2 für neues Spiel", 700, 600, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyF2) {
			columnSet = [6][7]int{}
			blinking = [6][7]bool{}
			nextPlayer = player1
			moves = 0
			gameIsOver = false
			time = 0
			blinkersVisible = true
			lastMouseX = -1
		}

		if window.WasKeyPressed(draw.KeyF9) {
			a := player1
			b := player2
			i := empty
			columnSet = [6][7]int{
				{i, a, b, i, b, a, b},
				{b, b, a, i, a, b, b},
				{a, a, a, i, a, a, a},
				{b, b, a, a, a, b, b},
				{b, a, b, a, b, a, b},
				{a, b, b, a, b, b, a},
			}
		}

		time++
		if time%20 == 0 {
			blinkersVisible = !blinkersVisible
		}

		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		mouseX, _ := window.MousePosition()
		if mouseX != lastMouseX {
			col := mouseX / 100
			if 0 <= col && col <= 6 {
				mouseColumn = col
			}
		}
		lastMouseX = mouseX

		if window.WasKeyPressed(draw.KeyLeft) {
			mouseColumn = (mouseColumn + 6) % 7
		}
		if window.WasKeyPressed(draw.KeyRight) {
			mouseColumn = (mouseColumn + 1) % 7
		}
		if window.WasKeyPressed(draw.Key1) || window.WasKeyPressed(draw.KeyNum1) {
			mouseColumn = 0
		}
		if window.WasKeyPressed(draw.Key2) || window.WasKeyPressed(draw.KeyNum2) {
			mouseColumn = 1
		}
		if window.WasKeyPressed(draw.Key3) || window.WasKeyPressed(draw.KeyNum3) {
			mouseColumn = 2
		}
		if window.WasKeyPressed(draw.Key4) || window.WasKeyPressed(draw.KeyNum4) {
			mouseColumn = 3
		}
		if window.WasKeyPressed(draw.Key5) || window.WasKeyPressed(draw.KeyNum5) {
			mouseColumn = 4
		}
		if window.WasKeyPressed(draw.Key6) || window.WasKeyPressed(draw.KeyNum6) {
			mouseColumn = 5
		}
		if window.WasKeyPressed(draw.Key7) || window.WasKeyPressed(draw.KeyNum7) {
			mouseColumn = 6
		}

		leftDown := window.WasKeyPressed(draw.KeyEnter) || window.WasKeyPressed(draw.KeyNumEnter)
		for _, click := range window.Clicks() {
			if click.Button == draw.LeftButton {
				leftDown = true
			}
		}

		if leftDown && !gameIsOver {
			y := nextEmptyRow(mouseColumn)
			if y != -1 {
				if nextPlayer == player1 {
					columnSet[y][mouseColumn] = player1
					nextPlayer = player2
				} else {
					columnSet[y][mouseColumn] = player2
					nextPlayer = player1
				}

				moves++
			}

			if hasPlayerWon(player1) || hasPlayerWon(player2) || isDraw() {
				gameIsOver = true
			}
		}

		// Draw the game.
		window.FillRect(0, 0, 700, 600, draw.White)
		for x := 0; x < 7; x++ {
			for y := 0; y < 6; y++ {
				window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.Black)

				if !blinking[y][x] || blinkersVisible {
					if columnSet[y][x] == player1 {
						window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGB(0.33, 0.33, 1))
					}
					if columnSet[y][x] == player2 {
						window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGB(1, 0.5, 0.38))
					}
				}

				if !gameIsOver && x == mouseColumn && y == nextEmptyRow(mouseColumn) {
					if nextPlayer == player1 {
						window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGBA(0.33, 0.33, 1, 0.75))
					} else {
						window.FillEllipse(x*100+10, y*100+10, 80, 80, draw.RGBA(1, 0.5, 0.38, 0.75))
					}
				}
			}
		}
	})
}
