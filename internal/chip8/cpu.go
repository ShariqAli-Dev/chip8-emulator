package chip8

import "log"

func (c *Chip8) EmulateCycle() {
	opcode := c.fetchOpcode()
}

func (c *Chip8) fetchOpcode() uint16 {
	highByte := c.memory[c.programCounter]
	lowByte := c.memory[c.programCounter+1]
	// opcodes consint of 2 bytes, bitwise operation combining high and low byte
	opcode := uint16(highByte)>>8 | uint16(lowByte)
	c.programCounter += 2 // points to the next opcode
	return opcode
}

func (c *Chip8) decodeAndExecuteOpcode(opcode uint16) {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			c.op00E0()
		case 0x00E:
			c.op00EE()
		default:
			log.Printf("unkwnow opcode 0x%X\n", opcode)
		}
	case 0x1000:
		c.op1NNN(opcode)
	// rest of the opcodes go here
	default:
		log.Printf("unknown opcode: 0x%X\n", opcode)
	}
}

// updates the delay and sound timers
func (c *Chip8) updateTimers() {
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.soundTimer--
		if c.soundTimer == 1 {
			// add sound to play here
		}
	}
}

// op00E0 clears the display
func (c *Chip8) op00E0() {
	for i := range c.graphics {
		c.graphics[i] = 0
	}
}

// op00EE returns from a subroutine
func (c *Chip8) op00EE() {
	c.programCounter = c.stack[c.stackPointer]
	c.stackPointer--
}

// op1NNN jumps to address NNN
func (c *Chip8) op1NNN(opcode uint16) {
	c.programCounter = opcode & 0x0FFF
}
