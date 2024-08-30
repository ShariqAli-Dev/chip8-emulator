package sdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("error initializing sdl: %v", err)
	}
	return nil
}

func Close() {
	sdl.Quit()
}
