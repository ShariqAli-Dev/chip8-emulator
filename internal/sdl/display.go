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
	window   *sdl.Window
	renderer *sdl.Renderer
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

	return err
}

func (d *display) Close() {
	d.renderer.Destroy()
	d.renderer = nil

	d.window.Destroy()
	d.window = nil
}

func (d *display) Render(display chip8.Display) {
	d.renderer.Clear()

	for j := 0; j < len(display); j++ {
		for i := 0; i < len(display[j]); i++ {
			if display[j][i] == 1 {
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
			d.renderer.SetDrawColor(45, 45, 45, 255)
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
				case sdl.SCANCODE_1:
					c.Keypad[0x1] = true
				case sdl.SCANCODE_2:
					c.Keypad[0x2] = true
				case sdl.SCANCODE_3:
					c.Keypad[0x3] = true
				case sdl.SCANCODE_4:
					c.Keypad[0xC] = true
				case sdl.SCANCODE_Q:
					c.Keypad[0x4] = true
				case sdl.SCANCODE_W:
					c.Keypad[0x5] = true
				case sdl.SCANCODE_E:
					c.Keypad[0x6] = true
				case sdl.SCANCODE_R:
					c.Keypad[0xD] = true
				case sdl.SCANCODE_A:
					c.Keypad[0x7] = true
				case sdl.SCANCODE_S:
					c.Keypad[0x8] = true
				case sdl.SCANCODE_D:
					c.Keypad[0x9] = true
				case sdl.SCANCODE_F:
					c.Keypad[0xE] = true
				case sdl.SCANCODE_Z:
					c.Keypad[0xA] = true
				case sdl.SCANCODE_X:
					c.Keypad[0x0] = true
				case sdl.SCANCODE_C:
					c.Keypad[0xB] = true
				case sdl.SCANCODE_V:
					c.Keypad[0xF] = true
				}
			} else if e.Type == sdl.KEYUP {
				switch e.Keysym.Scancode {
				case sdl.SCANCODE_1:
					c.Keypad[0x1] = false
				case sdl.SCANCODE_2:
					c.Keypad[0x2] = false
				case sdl.SCANCODE_3:
					c.Keypad[0x3] = false
				case sdl.SCANCODE_4:
					c.Keypad[0xC] = false
				case sdl.SCANCODE_Q:
					c.Keypad[0x4] = false
				case sdl.SCANCODE_W:
					c.Keypad[0x5] = false
				case sdl.SCANCODE_E:
					c.Keypad[0x6] = false
				case sdl.SCANCODE_R:
					c.Keypad[0xD] = false
				case sdl.SCANCODE_A:
					c.Keypad[0x7] = false
				case sdl.SCANCODE_S:
					c.Keypad[0x8] = false
				case sdl.SCANCODE_D:
					c.Keypad[0x9] = false
				case sdl.SCANCODE_F:
					c.Keypad[0xE] = false
				case sdl.SCANCODE_Z:
					c.Keypad[0xA] = false
				case sdl.SCANCODE_X:
					c.Keypad[0x0] = false
				case sdl.SCANCODE_C:
					c.Keypad[0xB] = false
				case sdl.SCANCODE_V:
					c.Keypad[0xF] = false
				}
			}
		}
	}

	if c.GetShouldDraw() {
		d.Render(c.GetDisplay())
	}
	sdl.Delay(uint32(1000 / fps))
	return true
}
