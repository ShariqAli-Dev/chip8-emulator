package chip8

import (
	"errors"
	"os"
)

const (
	programStart        = 0x200
	maxRomSize          = 0xFFF - 0x200
	fontSetStartAddress = 0x050
	DisplayWidth        = 64
	DisplayHeight       = 32
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

// http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#2.0
type Chip8 struct {
	memory      [4096]byte
	display     [DisplayWidth * DisplayHeight]byte // resolution pixel
	pc          uint16                             // program counter
	v           [16]byte                           // data register v0-vf
	i           uint16                             // index register
	delayTimer  byte
	soundTimer  byte
	stack       [12]uint16 // subroutine stack
	sp          uint16     // stack pointer
	keypad      [16]bool
	instruction struct {
		opcode uint16
		nnn    uint16 // 12-bit address/constant
		nn     uint8  // 8-bit constant
		n      uint8  // 4-bit constant
		x      uint8  // 4-bit register identifier
		y      uint8  // 4-bit register identifier
	}
}

func NewChip8() *Chip8 {
	chip8 := Chip8{
		pc: programStart,
	}
	for i := 0; i < len(fontSet); i++ {
		chip8.memory[fontSetStartAddress+i] = fontSet[i]
	}
	return &chip8
}

func (c *Chip8) GetDisplay() [DisplayWidth * DisplayHeight]byte {
	return c.display
}

func (c *Chip8) LoadROM(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(rom) > maxRomSize {
		return errors.New("error: rom too lang. Max size: 3583 bytes")
	}
	for index, bit := range rom {
		c.memory[programStart+index] = bit
	}
	return err
}
