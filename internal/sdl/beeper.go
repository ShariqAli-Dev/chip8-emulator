package sdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
)

type beeper struct {
	sound *mix.Chunk
}

func NewBeeper() *beeper {
	return &beeper{}
}

func (b *beeper) Init() error {
	var err error
	b.sound, err = mix.LoadWAV("./assets/sound/blip.wav")
	if err != nil {
		return fmt.Errorf("error loading go sound blip.wav: %v", err)
	}
	return err
}

func (b *beeper) PlaySound() {
	b.sound.Play(-1, 0)
}

func (b *beeper) Close() {
	b.sound.Free()
	b.sound = nil
}
