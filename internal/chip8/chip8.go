package chip8

import (
	"errors"
	"os"
)

const (
	programStart = 0x200
	maxRomSize   = 0xFFF - 0x200
)

// http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#2.0
type Chip8 struct {
	memory         [4096]byte
	opcode         uint16
	v              [16]byte
	indexRegister  uint16
	programCounter uint16
	graphics       [64 * 32]byte

	delayTimer byte
	soundTimer byte

	stack        [16]uint16
	stackPointer uint16

	keypad [16]byte
}

func NewChip8() *Chip8 {
	return &Chip8{
		programCounter: programStart,
	}
}

func (chip *Chip8) LoadROM(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(rom) > maxRomSize {
		return errors.New("error: rom too lang. Max size: 3583 bytes")
	}
	for index, byte := range rom {
		chip.memory[programStart+index] = byte
	}
	return nil
}
