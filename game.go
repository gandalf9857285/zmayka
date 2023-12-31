package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const (
	GameW = 720.0
	GameH = 720.0
)

type Game struct {
	cv       *canvas.Canvas
	wnd      *sdlcanvas.Window
	worldS   float64
	snake    *Snake
	needMove bool
	food     []Point
	speed    int
	gameOver bool
}

func NewGame() *Game {
	wnd, cv, err := sdlcanvas.CreateWindow(1080, GameH+30, "Hello, Snake!")
	if err != nil {
		panic(err)
	}
	g := &Game{
		cv:       cv,
		wnd:      wnd,
		speed:    300,
		gameOver: false,
	}

	return g
}

func (g *Game) SetSnake(s *Snake) {
	g.snake = s
}

func (g *Game) CreateWorld(s float64) {
	g.worldS = s
}

func (g *Game) Run() {
	go g.snakeMovement()
	go g.foodGeneration()
	g.renderLoop()
}

func (g *Game) Exit() {
	defer g.wnd.Destroy()
}

func (g *Game) snakeMovement() {
	var snakeTimer *time.Timer
	var snakeDir Dir = Right
	var snakeLock sync.Mutex

	resetTimer := func() {
		snakeTimer = time.NewTimer(time.Duration(g.speed) * time.Millisecond)
	}
	resetTimer()

	g.wnd.KeyUp = func(code int, rn rune, name string) {
		if code < 79 && code > 82 || g.needMove {
			return
		}
		snakeLock.Lock()
		newDir := snakeDir
		switch code {
		case 80:
			newDir = Left
		case 82:
			newDir = Bottom
		case 79:
			newDir = Right
		case 81:
			newDir = Top
		}

		if !snakeDir.CheckParallel(newDir) {
			snakeDir = newDir
			g.needMove = true
		}

		snakeLock.Unlock()
	}

	for {
		<-snakeTimer.C
		snakeLock.Lock()

		if !g.gameOver {
			newPos := snakeDir.Exec(g.snake.Parts[0])
			if newPos.X <= 0 || newPos.X >= g.worldS-1 ||
				newPos.Y <= 0 || newPos.Y >= g.worldS-1 {
				g.gameOver = true
			}

			g.snake.CutIfSnake(newPos)

			isFood := false
			for i := range g.food {
				if newPos.X == g.food[i].X && newPos.Y == g.food[i].Y {
					g.food = append(g.food[:i], g.food[i+1:]...)
					g.snake.Add(newPos)
					g.speed -= 5
					isFood = true
					break
				}
			}
			if !isFood {
				g.snake.Move(snakeDir)
				g.needMove = false
			}
		}
		snakeLock.Unlock()
		resetTimer()
	}
}

func (g *Game) foodGeneration() {
	var foodTimer *time.Timer
	resetTimer := func() {
		foodTimer = time.NewTimer(1 * time.Second)
	}
	resetTimer()

	for {
		<-foodTimer.C
		if !g.gameOver {
			min := 1
			max := int(g.worldS) - 1
			randX := rand.Intn(max-min) + min
			randY := rand.Intn(max-min) + min
			newPoint := Point{float64(randX), float64(randY)}
			check := true
			if g.snake.IsSnake(newPoint) {
				check = false
			}
			for _, p := range g.food {
				if p.X == newPoint.X && p.Y == newPoint.Y {
					check = false
					break
				}
			}
			if check {
				g.food = append(g.food, newPoint)
			}
		}
		resetTimer()
	}
}

func (g *Game) renderLoop() {

	gameAreaSP := Point{X: 15, Y: 15}
	gameAreaEP := Point{GameW + gameAreaSP.X, GameH + gameAreaSP.Y}

	cellW := GameW / g.worldS
	cellH := GameH / g.worldS

	g.wnd.MainLoop(func() {
		g.cv.ClearRect(0, 0, 1080, 750)
		g.cv.BeginPath()
		g.cv.SetFillStyle("#06C258")
		g.cv.FillRect(gameAreaSP.X, gameAreaSP.Y, gameAreaEP.X-15, gameAreaEP.Y-15)
		g.cv.Stroke()
		g.cv.BeginPath()
		g.cv.SetStrokeStyle("#FFF001")
		g.cv.SetLineWidth(1)
		for i := 0; i < int(g.worldS)+1; i++ {
			g.cv.MoveTo(gameAreaSP.X+float64(i)*cellW, gameAreaSP.Y)
			g.cv.LineTo(gameAreaSP.X+float64(i)*cellW, gameAreaEP.Y)
		}
		for i := 0; i < int(g.worldS)+1; i++ {
			g.cv.MoveTo(gameAreaSP.X, gameAreaSP.Y+float64(i)*cellH)
			g.cv.LineTo(gameAreaEP.X, gameAreaSP.Y+float64(i)*cellH)
		}
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetFillStyle("#FFF001")

		for i := 0; i < int(g.worldS); i++ {
			g.cv.FillRect(
				gameAreaSP.X+float64(i)*cellW+1,
				gameAreaSP.Y,
				cellW-1*2,
				cellH-1*2)
		}

		for i := 0; i < int(g.worldS); i++ {
			g.cv.FillRect(
				gameAreaSP.X+float64(i)*cellW+1,
				gameAreaSP.Y+cellH*(g.worldS-1),
				cellW,
				cellH-1*2)
		}

		for i := 1; i < int(g.worldS)-1; i++ {
			g.cv.FillRect(
				gameAreaSP.X,
				gameAreaSP.Y+float64(i)*cellH+1,
				cellW,
				cellH-1*2)
		}

		for i := 1; i < int(g.worldS)-1; i++ {
			g.cv.FillRect(
				gameAreaSP.X+cellW*(g.worldS-1),
				gameAreaSP.Y+float64(i)*cellH+1,
				cellW-1*2,
				cellH-1*2)
		}
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetFillStyle("#FFF")
		for _, p := range g.snake.Parts {
			g.cv.FillRect(
				gameAreaSP.X+p.X*cellW+1,
				gameAreaSP.Y+p.Y*cellH+1,
				cellW-1*2,
				cellH-1*2)
		}
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetFillStyle("#F15555")
		for _, p := range g.food {
			g.cv.FillRect(
				gameAreaSP.X+p.X*cellW+1,
				gameAreaSP.Y+p.Y*cellH+1,
				cellW-1*2,
				cellH-1*2)
		}
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text := fmt.Sprintf("Score: %d", g.snake.Len())
		g.cv.FillText(text, GameW+50, 50)

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text = fmt.Sprintf("Food: %d", len(g.food))
		g.cv.FillText(text, GameW+50, 85)

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text = fmt.Sprintf("Speed: %d", 350-g.speed)
		g.cv.FillText(text, GameW+50, 120)

		if g.gameOver {
			g.cv.BeginPath()
			g.cv.SetFont(font, 30)
			text = fmt.Sprintf("Game over")
			g.cv.FillText(text, GameW+100, 175)
		}
	})
}
