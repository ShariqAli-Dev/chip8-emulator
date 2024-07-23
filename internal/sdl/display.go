package sdl

import (
	"fmt"
	"unsafe"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	fps           = 30
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

// Render draws the current state of the CHIP-8 display
func (d *display) Render(display [64 * 32]byte) {
	// Create a buffer to hold pixel data
	pixels := make([]byte, displayWidth*displayHeight*4) // 4 bytes per pixel (RGBA)

	// Convert CHIP-8 display data to RGBA pixel data
	for i := 0; i < len(display); i++ {
		index := i * 4
		if display[i] == 1 {
			pixels[index] = 255   // R
			pixels[index+1] = 255 // G
			pixels[index+2] = 255 // B
			pixels[index+3] = 255 // A
		} else {
			pixels[index] = 0     // R
			pixels[index+1] = 0   // G
			pixels[index+2] = 0   // B
			pixels[index+3] = 255 // A
		}
	}

	// Update the texture with the new pixel data
	d.pixelTexture.Update(nil, unsafe.Pointer(&pixels[0]), displayWidth*4)

	// Clear the renderer and draw the texture
	d.renderer.SetDrawColor(0, 0, 0, 255) // black
	d.renderer.Clear()

	destRect := sdl.Rect{X: 0, Y: 0, W: displayWidth * scaleFactor, H: displayHeight * scaleFactor}
	d.renderer.Copy(d.pixelTexture, nil, &destRect)

	// Present the renderer
	d.renderer.Present()

	// Cap the frame rate
	sdl.Delay(uint32(1000 / fps))
}

func (d *display) Tick(c *chip8.Chip8) bool {
	// Handle SDL events
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
