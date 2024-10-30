package main

import (
	"math"

	"github.com/gonutz/prototype/draw"
)

func main() {
	draw.RunWindow("Close me", 640, 480, func(window draw.Window) {
		w, h := window.Size()
		centerX, centerY := w/2, h/2

		mouseX, mouseY := window.MousePosition()
		mouseInCircle := math.Hypot(float64(mouseX-centerX), float64(mouseY-centerY)) < 20
		color := draw.DarkRed
		if mouseInCircle {
			color = draw.Red
		}
		window.FillEllipse(centerX-20, centerY-20, 40, 40, color)
		window.DrawEllipse(centerX-20, centerY-20, 40, 40, draw.White)
		if mouseInCircle {
			window.DrawScaledText("Close!", centerX-40, centerY+25, 1.6, draw.Green)
		}

		for _, click := range window.Clicks() {
			dx, dy := click.X-centerX, click.Y-centerY
			squareDist := dx*dx + dy*dy
			if squareDist <= 20*20 {
				window.Close()
			}
		}
	})
}
