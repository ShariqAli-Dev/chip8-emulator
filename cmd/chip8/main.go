package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/shariqali-dev/chip8-emulator/internal/sdl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <exe> <rom_path>")
	}

	chip8 := chip8.NewChip8()
	if err := chip8.LoadROM(os.Args[2]); err != nil {
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
