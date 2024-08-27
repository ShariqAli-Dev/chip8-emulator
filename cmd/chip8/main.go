package main

import (
	"log"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/shariqali-dev/chip8-emulator/internal/sdl"
)

// https://tobiasvl.github.io/blog/write-a-chip-8-emulator/
// https://github.dev/skatiyar/go-chip8

var roms = []string{"IBM Logo.ch8", "test_opcode.ch8", "Airplane.ch8"}

func main() {
	chip8 := chip8.NewChip8()
	if err := chip8.LoadROM("./roms/" + roms[0]); err != nil {
		log.Fatal(error.Error(err))
	}

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
		chip8.EmulateCycle()
	}
}
