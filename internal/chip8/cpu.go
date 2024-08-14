package chip8

import "fmt"

func (c *Chip8) EmulateCycle() {
	highByte := c.memory[c.pc]
	lowByte := c.memory[c.pc+1]
	c.instruction.opcode = uint16(highByte)<<8 | uint16(lowByte)
	c.pc += 2

	// Fill out current instruction format
	// DXYN
	c.instruction.nnn = c.instruction.opcode & 0x0FFF
	c.instruction.nn = uint8(c.instruction.opcode & 0x00FF)
	c.instruction.n = uint8(c.instruction.opcode & 0x000F)
	c.instruction.x = uint8((c.instruction.opcode >> 8) & 0x000F)
	c.instruction.y = uint8((c.instruction.opcode >> 4) & 0x000F)

	category := (c.instruction.opcode >> 12) & 0x0F
	switch category {
	case 0x00:
		if c.instruction.nn == 0xE0 {
			// 0x00E0: Clear the screen
			for i := range c.display {
				c.display[i] = 0
			}
		} else if c.instruction.nn == 0xEE {
			// 0x0EE: retun from subroutine
			if c.sp > 0 {
				c.sp--
				c.pc = c.stack[c.sp]
			}
		}
	case 0x02:
		// 0x2NNN: call subroutine at nnn
		if c.sp < uint16(len(c.stack)) {
			c.stack[c.sp] = c.pc
			c.sp++
			c.pc = c.instruction.opcode & 0x0FFF
		}
	case 0x06:
		// 0x6XNN: set register vx to NN
		c.v[c.instruction.x] = c.instruction.nn
	case 0x0A:
		// 0xANNN: set index register i to nnnn
		c.i = c.instruction.opcode & 0x0FFF
	case 0x0D:
		// 0x0DXYN: draw n height sprite at these cooridens
		x := (c.instruction.opcode & 0x0F00) >> 8
		y := (c.instruction.opcode & 0x00F0) >> 4
		height := c.instruction.opcode & 0x000F

		c.v[0xF] = 0
		for yline := 0; yline < int(height); yline++ {
			pixel := c.memory[c.i+uint16(yline)]
			for xline := 0; xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if c.display[(x+uint16(xline)+((y+uint16(yline))*64))] == 1 {
						c.v[0xF] = 1
					}
					c.display[x+uint16(xline)+((y+uint16(yline))*64)] ^= 1
				}
			}
		}
	default:
		fmt.Printf("opcode is unimplemented or invalid: %d\n", c.instruction.opcode)
	}
}

// func (c *Chip8) decodeAndExecuteOpcode(opcode uint16) {
// 	switch opcode & 0xF000 {
// 	case 0x0000:
// 		switch opcode & 0x000F {
// 		case 0x0000:
// 			c.op00E0()
// 		case 0x00E:
// 			c.op00EE()
// 		default:
// 			log.Printf("unkwnow opcode 0x%X\n", opcode)
// 		}
// 	case 0x1000:
// 		c.op1NNN(opcode)
// 	// rest of the opcodes go here
// 	default:
// 		log.Printf("unknown opcode: 0x%X\n", opcode)
// 	}
// }

// // updates the delay and sound timers
// func (c *Chip8) updateTimers() {
// 	if c.delayTimer > 0 {
// 		c.delayTimer--
// 	}
// 	if c.soundTimer > 0 {
// 		c.soundTimer--
// 		if c.soundTimer == 1 {
// 			// add sound to play here
// 		}
// 	}
// }

// // op00E0 clears the display
// func (c *Chip8) op00E0() {
// 	for i := range c.graphics {
// 		c.graphics[i] = 0
// 	}
// }

// // op00EE returns from a subroutine
// func (c *Chip8) op00EE() {
// 	c.programCounter = c.stack[c.stackPointer]
// 	c.stackPointer--
// }

// // op1NNN jumps to address NNN
// func (c *Chip8) op1NNN(opcode uint16) {
// 	c.programCounter = opcode & 0x0FFF
// }
