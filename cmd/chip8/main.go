package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/shariqali-dev/chip8-emulator/internal/chip8"
	"github.com/shariqali-dev/chip8-emulator/internal/sdl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <exe> <rom_path>")
		return
	}

	romAbsolutePath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(fmt.Errorf("error resolving rom path: %s", err))
	}

	chip8 := chip8.NewChip8()
	if err := chip8.LoadROM(romAbsolutePath); err != nil {
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
	}
}
