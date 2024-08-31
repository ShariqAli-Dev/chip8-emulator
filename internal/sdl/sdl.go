package sdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func Init() error {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("error initializing sdl: %v", err)
	}
	if err = mix.Init(mix.INIT_FLAC | mix.INIT_OGG); err != nil {
		return fmt.Errorf("error initializing sdl mixer: %v", err)
	}
	if err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		return fmt.Errorf("error opening audio: %v", err)
	}

	return err
}

func Close() {
	mix.CloseAudio()
	mix.Quit()
	sdl.Quit()
}
