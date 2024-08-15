package sdl

import (
	"fmt"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	fps           = 60
	scaleFactor   = 15
	displayWidth  = chip8.DisplayWidth * scaleFactor
	displayHeight = chip8.DisplayHeight * scaleFactor
	windowTitle   = "Chip8 Emulator"
)

type display struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	pixelTexture *sdl.Texture
}

func NewDisplay() *display {
	return &display{}
}

func (d *display) Init() error {
	var err error

	d.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, displayWidth, displayHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("error creating window: %v", err)
	}

	d.renderer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("error creating renderer: %v", err)
	}

	d.pixelTexture, err = d.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, displayWidth, displayHeight)
	if err != nil {
		return fmt.Errorf("error creating pixel texturer: %v", err)
	}

	return err
}

func (d *display) Close() {
	d.pixelTexture.Destroy()
	d.pixelTexture = nil

	d.renderer.Destroy()
	d.renderer = nil

	d.window.Destroy()
	d.window = nil
}

func (d *display) Render(display chip8.Display) {
	d.renderer.SetDrawColor(255, 0, 0, 255)
	d.renderer.Clear()

	for j := 0; j < len(display); j++ {
		for i := 0; i < len(display[j]); i++ {
			if display[j][i] != 0 {
				d.renderer.SetDrawColor(255, 255, 255, 255)
			} else {

				d.renderer.SetDrawColor(0, 0, 0, 255)
			}
			d.renderer.FillRect(&sdl.Rect{
				Y: int32(j) * scaleFactor,
				X: int32(i) * scaleFactor,
				W: scaleFactor,
				H: scaleFactor,
			})

			// rectangle border (smaller rectangle)
			d.renderer.SetDrawColor(0, 0, 0, 255)
			borderThickness := int32(1)
			d.renderer.DrawRect(&sdl.Rect{
				Y: int32(j)*scaleFactor - borderThickness,
				X: int32(i)*scaleFactor - borderThickness,
				W: scaleFactor + 2*borderThickness,
				H: scaleFactor + 2*borderThickness,
			})
		}
	}

	d.renderer.Present()
	sdl.Delay(uint32(1000 / fps))
}

func (d *display) Tick(c *chip8.Chip8) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Scancode {
				case sdl.SCANCODE_ESCAPE:
					return false
				}
			}
		}
	}

	d.Render(c.GetDisplay())

	return true
}
