package chip8

import "fmt"

func (c *Chip8) EmulateCycle() {
	highByte := c.memory[c.pc]
	lowByte := c.memory[c.pc+1]
	c.instruction.opcode = (uint16(highByte) << 8) | uint16(lowByte)
	c.pc += 2
	fmt.Printf("opcode: %d\n", c.instruction.opcode)

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
			// for i := range c.display {
			// c.display[i] = 0
			// }
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
		_, _ = x, y
		height := c.instruction.opcode & 0x000F

		c.v[0xF] = 0
		for yline := 0; yline < int(height); yline++ {
			pixel := c.memory[c.i+uint16(yline)]
			for xline := 0; xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					// if c.display[(x+uint16(xline)+((y+uint16(yline))*64))] == 1 {
					// 	c.v[0xF] = 1
					// }
					// c.display[x+uint16(xline)+((y+uint16(yline))*64)] ^= 1
				}
			}
		}
	default:
		fmt.Printf("opcode is unimplemented or invalid: %d\n", c.instruction.opcode)
	}
}
