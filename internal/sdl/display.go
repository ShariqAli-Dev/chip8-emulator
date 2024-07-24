package sdl

import (
	"fmt"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	fps           = 30
	scaleFactor   = 20
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

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("error initializing SDL: %v", err)
	}

	d.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, displayWidth, displayHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("error creating window: %v", err)
	}

	d.renderer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		d.window.Destroy()
		return fmt.Errorf("error creating renderer: %v", err)
	}

	d.pixelTexture, err = d.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, displayWidth, displayHeight)
	if err != nil {
		d.renderer.Destroy()
		d.window.Destroy()
		return fmt.Errorf("error creating pixel texture: %v", err)
	}

	return nil
}

func (d *display) Close() {
	if d.pixelTexture != nil {
		d.pixelTexture.Destroy()
		d.pixelTexture = nil
	}

	if d.renderer != nil {
		d.renderer.Destroy()
		d.renderer = nil
	}

	if d.window != nil {
		d.window.Destroy()
		d.window = nil
	}

	sdl.Quit()
}

func (d *display) Tick(c *chip8.Chip8) {
	for {
		c.EmulateCycle()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Scancode {
					case sdl.SCANCODE_ESCAPE:
						return
					}
				}
			}
		}

		d.renderer.SetDrawColor(0, 0, 0, 255) // black
		d.renderer.Clear()

		for xPixel := 0; xPixel < chip8.DisplayWidth; xPixel++ {
			for yPixel := 0; yPixel < chip8.DisplayHeight; yPixel++ {
				if c.GetDisplay()[yPixel*chip8.DisplayWidth+xPixel] == 1 {
					d.renderer.SetDrawColor(255, 255, 255, 255) // white
				} else {
					d.renderer.SetDrawColor(0, 0, 0, 255) // black
				}
				pixel := sdl.Rect{
					X: int32(xPixel) * scaleFactor,
					Y: int32(yPixel) * scaleFactor,
					W: scaleFactor,
					H: scaleFactor,
				}
				d.renderer.FillRect(&pixel)
			}
		}

		d.renderer.Present()
		sdl.Delay(uint32(1000 / fps))
	}
}
