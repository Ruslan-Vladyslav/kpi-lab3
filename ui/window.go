package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event:
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonRight {
				log.Printf("Clicked: %v, %v", e.X, e.Y)
				x := int(e.X)
				y := int(e.Y)
				pw.drawDefaultUI(&x, &y)
			}
		}

	case paint.Event:
		if t == nil {
			pw.drawDefaultUI(nil, nil)
		} else {
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI(x, y *int) {
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{0, 255, 0, 255}, draw.Src) 

	if x == nil {
		defaultX := pw.sz.WidthPx / 2
		x = &defaultX
	}
	if y == nil {
		defaultY := pw.sz.HeightPx / 2
		y = &defaultY
	}

	DrawShape(pw.w.Fill, *x, *y, 1)

	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func DrawShape(Fill func(dr image.Rectangle, src color.Color, op draw.Op), x, y int, scale float64) {
	size := int(200 * scale)      
	thickness := int(40 * scale)   

	shapeColor := color.RGBA{R: 255, G: 230, B: 69, A: 255}

	vertical := image.Rect(
		x-thickness/2,
		y-size/2,
		x+thickness/2,
		y+size/2,
	)

	horizontal := image.Rect(
		x-thickness/2-size/2,
		y-thickness/2,
		x-thickness/2,
		y+thickness/2,
	)

	Fill(vertical, shapeColor, draw.Src)
	Fill(horizontal, shapeColor, draw.Src)
}


