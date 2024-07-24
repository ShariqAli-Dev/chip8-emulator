package chip8

import "math/rand/v2"

// Define the opcode functions
func (c *Chip8) op00E0() {
	for i := range c.display {
		c.display[i] = 0
	}
}

func (c *Chip8) op00EE() {
	// Ensure the stack pointer is not below 0
	if c.stackPointer > 0 {
		c.stackPointer--
		c.programCounter = c.stack[c.stackPointer]
	}
}

func (c *Chip8) op1NNN(opcode uint16) {
	c.programCounter = opcode & 0x0FFF
}

func (c *Chip8) op2NNN(opcode uint16) {
	// Ensure the stack pointer does not exceed its limit
	if c.stackPointer < uint16(len(c.stack)) {
		c.stack[c.stackPointer] = c.programCounter
		c.stackPointer++
		c.programCounter = opcode & 0x0FFF
	}
}

func (c *Chip8) op3XNN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := byte(opcode & 0x00FF)
	if c.V[x] == nn {
		c.programCounter += 2
	}
}

func (c *Chip8) op4XNN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := byte(opcode & 0x00FF)
	if c.V[x] != nn {
		c.programCounter += 2
	}
}

func (c *Chip8) op5XY0(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if c.V[x] == c.V[y] {
		c.programCounter += 2
	}
}

func (c *Chip8) op6XNN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.V[x] = byte(opcode & 0x00FF)
}

func (c *Chip8) op7XNN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.V[x] += byte(opcode & 0x00FF)
}

func (c *Chip8) op8XY0(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	c.V[x] = c.V[y]
}

func (c *Chip8) op8XY1(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	c.V[x] |= c.V[y]
}

func (c *Chip8) op8XY2(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	c.V[x] &= c.V[y]
}

func (c *Chip8) op8XY3(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	c.V[x] ^= c.V[y]
}

func (c *Chip8) op8XY4(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	sum := uint16(c.V[x]) + uint16(c.V[y])
	if sum > 255 {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[x] = byte(sum)
}

func (c *Chip8) op8XY5(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if c.V[x] > c.V[y] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[x] -= c.V[y]
}

func (c *Chip8) op8XY6(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.V[0xF] = c.V[x] & 0x1
	c.V[x] >>= 1
}

func (c *Chip8) op8XY7(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if c.V[y] > c.V[x] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[x] = c.V[y] - c.V[x]
}

func (c *Chip8) op8XYE(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.V[0xF] = (c.V[x] & 0x80) >> 7
	c.V[x] <<= 1
}

func (c *Chip8) op9XY0(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if c.V[x] != c.V[y] {
		c.programCounter += 2
	}
}

func (c *Chip8) opANNN(opcode uint16) {
	c.indexRegister = opcode & 0x0FFF
}

func (c *Chip8) opBNNN(opcode uint16) {
	c.programCounter = uint16(c.V[0]) + (opcode & 0x0FFF)
}

func (c *Chip8) opCXNN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := byte(opcode & 0x00FF)
	c.V[x] = byte(rand.IntN(256)) & nn
}

func (c *Chip8) opDXYN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	height := opcode & 0x000F

	c.V[0xF] = 0
	for yline := 0; yline < int(height); yline++ {
		pixel := c.memory[c.indexRegister+uint16(yline)]
		for xline := 0; xline < 8; xline++ {
			if (pixel & (0x80 >> xline)) != 0 {
				if c.display[(x+uint16(xline)+((y+uint16(yline))*64))] == 1 {
					c.V[0xF] = 1
				}
				c.display[x+uint16(xline)+((y+uint16(yline))*64)] ^= 1
			}
		}
	}
}

func (c *Chip8) opEX9E(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	if c.keypad[c.V[x]] != 0 {
		c.programCounter += 2
	}
}

func (c *Chip8) opEXA1(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	if c.keypad[c.V[x]] == 0 {
		c.programCounter += 2
	}
}

func (c *Chip8) opFX07(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.V[x] = c.delayTimer
}

func (c *Chip8) opFX0A(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	keyPress := false
	for i := 0; i < 16; i++ {
		if c.keypad[i] != 0 {
			c.V[x] = byte(i)
			keyPress = true
		}
	}
	if !keyPress {
		c.programCounter -= 2
	}
}

func (c *Chip8) opFX15(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.delayTimer = c.V[x]
}

func (c *Chip8) opFX18(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.soundTimer = c.V[x]
}

func (c *Chip8) opFX1E(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.indexRegister += uint16(c.V[x])
}

func (c *Chip8) opFX29(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.indexRegister = uint16(c.V[x]) * 0x5
}

func (c *Chip8) opFX33(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	c.memory[c.indexRegister] = c.V[x] / 100
	c.memory[c.indexRegister+1] = (c.V[x] / 10) % 10
	c.memory[c.indexRegister+2] = (c.V[x] % 100) % 10
}

func (c *Chip8) opFX55(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	for i := uint16(0); i <= x; i++ {
		c.memory[c.indexRegister+i] = c.V[i]
	}
}

func (c *Chip8) opFX65(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	for i := uint16(0); i <= x; i++ {
		c.V[i] = c.memory[c.indexRegister+i]
	}
}
