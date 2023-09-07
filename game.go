package main

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const (
	GameW = 720.0
	GameH = 720.0
)

type Game struct {
	cv  *canvas.Canvas
	wnd *sdlcanvas.Window
}

func NewGame() *Game {
	wnd, cv, err := sdlcanvas.CreateWindow(1080, 750, "Hello, Snake!")
	if err != nil {
		panic(err)
	}
	g := &Game{
		cv:  cv,
		wnd: wnd,
	}
	return g
}

func (g *Game) Run() {
	g.renderLoop()
}

func (g *Game) renderLoop() {

	gameAreaSP := Point{X: 15, Y: 15}
	gameAreaEP := Point{X: 15 + GameW, Y: 15 + GameH}

	//cellW := GameW / 20
	//cellH := GameH / 20

	g.wnd.MainLoop(func() {
		// clear
		g.cv.ClearRect(0, 0, 1080, 750)
		//render
		g.cv.BeginPath()
		g.cv.SetFillStyle("#06C258")
		g.cv.FillRect(gameAreaSP.X, gameAreaSP.Y, gameAreaEP.X-15, gameAreaEP.Y-15)
		g.cv.Stroke()
	})
}
