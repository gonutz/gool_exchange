package main

import (
	"math"
	"math/rand"

	"github.com/gonutz/prototype/draw"
)

const (
	bulletSize      = 20
	bulletSpeed     = 12
	bulletDelay     = 3
	bulletStartLive = 3.0
	bulletDecay     = 0.03
)

const (
	enemySize       = 80
	enemySpeed      = 7
	enemySpawnDelay = 30
	enemyStartLive  = 3
	enemyDamage     = 1.1
)

const (
	initialPlayerSize = 40.0
	playerDamagePerHit = 3
	playerBonusSizePerKill = 0.25
)

func main() {
	var bullets []bullet
	lastBullet := -bulletDelay
	playerSize := initialPlayerSize

	var enemies []enemy
	lastEnemySpawn := 0

	time := 0

	draw.RunWindow("Shooter", 800, 600, func(window draw.Window) {
		time++

		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		window.SetFullscreen(true)

		windowW, windowH := window.Size()
		cx, cy := float64(windowW)/2, float64(windowH)/2

		if window.IsMouseDown(draw.LeftButton) {
			if time-lastBullet > bulletDelay {
				lastBullet = time
				cx := float64(windowW) / 2
				cy := float64(windowH) / 2
				mx, my := window.MousePosition()
				dx := float64(mx) - cx
				dy := float64(my) - cy
				length := math.Hypot(dx, dy)
				dx /= length
				dy /= length
				if length > 0 {
					bullets = append(bullets, bullet{
						x:    cx,
						y:    cy,
						dx:   bulletSpeed * dx,
						dy:   bulletSpeed * dy,
						live: bulletStartLive,
					})
				}
			}
		}

		if time-lastEnemySpawn > enemySpawnDelay {
			lastEnemySpawn = time
			dy, dx := math.Sincos(2 * math.Pi * rand.Float64())
			length := float64(max(windowW, windowH)) * math.Sqrt2
			enemies = append(enemies, enemy{
				x:    cx - dx*length,
				y:    cy - dy*length,
				dx:   enemySpeed * dx,
				dy:   enemySpeed * dy,
				live: enemyStartLive,
			})
		}

		bulletCount := 0
		for i := range bullets {
			bullets[i].x += bullets[i].dx
			bullets[i].y += bullets[i].dy
			bullets[i].live -= bulletDecay
			if bullets[i].live > 0 {
				bullets[bulletCount] = bullets[i]
				bulletCount++
			}
		}
		bullets = bullets[:bulletCount]

		enemyCount := 0
		for i := range enemies {
			enemies[i].x += enemies[i].dx
			enemies[i].y += enemies[i].dy

			hit := func() bool {
				for j := range bullets {
					if collide(bullets[j], enemies[i]) {
						bullets[j].live = 0
						return true
					}
				}
				return false
			}()
			if hit {
				enemies[i].live -= enemyDamage
			}

			if enemies[i].live > 0 {
				enemies[enemyCount] = enemies[i]
				enemyCount++
			} else {
				playerSize += playerBonusSizePerKill
			}
		}
		enemies = enemies[:enemyCount]

		playerWasHit := false
		for i, e := range enemies {
			hits := abs(e.x-cx) < 10 && abs(e.y-cy) < 10
			playerWasHit = playerWasHit || hits
			if hits {
				enemies[i].live = 0
			}
		}
		if playerWasHit {
			playerSize -= playerDamagePerHit
			if playerSize <= 0 {
				window.Close() // TODO
			}
		}

		for _, b := range bullets {
			x, y := round(b.x), round(b.y)
			window.FillEllipse(
				x-bulletSize/2,
				y-bulletSize/2,
				bulletSize,
				bulletSize,
				draw.RGBA(1, 0, 0, min(1, float32(b.live))),
			)
		}

		window.FillEllipse(
			windowW/2-round(playerSize),
			windowH/2-round(playerSize),
			round(2*playerSize),
			round(2*playerSize),
			draw.LightRed,
		)

		for _, e := range enemies {
			x, y := round(e.x), round(e.y)
			window.FillRect(
				x-enemySize/2,
				y-enemySize/2,
				enemySize,
				enemySize,
				draw.Blue,
			)
		}
	})
}

type bullet struct {
	x, y   float64
	dx, dy float64
	live   float64
}

type enemy struct {
	x, y   float64
	dx, dy float64
	live   float64
}

func round(x float64) int {
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}

func collide(b bullet, e enemy) bool {
	de := enemySize / 2.0
	x := max(e.x-de, min(e.x+de, b.x))
	y := max(e.y-de, min(e.y+de, b.y))
	dx := b.x - x
	dy := b.y - y
	d := dx*dx + dy*dy
	bSize := bulletSize / 2.0
	return d <= bSize*bSize
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
