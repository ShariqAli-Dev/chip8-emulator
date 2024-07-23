package sdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func Init() error {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("error initializing sdl: %v", err)
	}
	return err
}

func Close() {
	sdl.Quit()
}
