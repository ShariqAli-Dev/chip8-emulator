package main

import (
	"log"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/shariqali-dev/chip8-emulator/internal/sdl"
)

func main() {
	chip8 := chip8.NewChip8()

	defer sdl.Close()
	if err := sdl.Init(); err != nil {
		log.Fatal(error.Error(err))
	}

	display := sdl.NewDisplay()
	defer display.Close()
	if err := display.Init(); err != nil {
		log.Fatal(error.Error(err))
	}

	for display.Tick(chip8) {

	}
}
