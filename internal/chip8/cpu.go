package chip8

import "fmt"

func (c *Chip8) EmulateCycle() {
	highByte := c.memory[c.pc]
	lowByte := c.memory[c.pc+1]
	c.instruction.opcode = (uint16(highByte) << 8) | uint16(lowByte)

	// DXYN
	c.instruction.nnn = c.instruction.opcode & 0x0FFF
	c.instruction.nn = uint8(c.instruction.opcode & 0x00FF)
	c.instruction.n = uint8(c.instruction.opcode & 0x000F)
	c.instruction.x = uint8((c.instruction.opcode >> 8) & 0x000F)
	c.instruction.y = uint8((c.instruction.opcode >> 4) & 0x000F)

	switch c.instruction.opcode & 0xF000 {
	case 0x0000:
		switch c.instruction.opcode & 0x000F {
		case 0x0000: // 0x00E0 clears screen
			for i := 0; i < len(c.display); i++ {
				for j := 0; j < len(c.display[i]); j++ {
					c.display[i][j] = 0
				}
			}
			c.pc += 2
		case 0x000E: // 0x00EE returns from a subroutine
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2
		default:
			fmt.Printf("INVALID OPCODE %X\n", c.instruction.opcode)
		}
	case 0x1000: // 0x1NNN jump to address NNN
		c.pc = c.instruction.opcode & 0x0FFF
	case 0x2000: // 0x2NNN calls subroutine at NNN
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.instruction.opcode & 0x0FFF
	case 0x6000: // 0x6XNN sets VX to NN
		c.v[c.instruction.x] = c.instruction.nn
		c.pc += 2
	case 0x7000: // 0x7XNN adds NN to VX (without carry)
		c.v[c.instruction.x] += c.instruction.nn
		c.pc += 2
	case 0xA000: // 0xANNN sets I to the address NNN
		c.i = c.instruction.opcode & 0x0FFF
		c.pc += 2
	case 0xD000: // 0xDXYN draws a sprite at coordinate (VX, VY) with a height of N
		h := c.instruction.opcode & 0x000F
		c.v[0xF] = 0
		for j := uint16(0); j < h; j++ {
			pixel := c.memory[c.i+j]
			for i := uint16(0); i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if c.display[c.v[c.instruction.y]+uint8(j)][c.v[c.instruction.x]+uint8(i)] == 1 {
						c.v[0xF] = 1
					}
					c.display[c.v[c.instruction.y]+uint8(j)][c.v[c.instruction.x]+uint8(i)] ^= 1
				}
			}
		}
		c.pc += 2
	case 0xF000:
		switch c.instruction.opcode & 0x00FF {
		case 0x1E: // 0xFX1E adds VX to I
			c.i += uint16(c.v[c.instruction.x])
			c.pc += 2
		case 0x55: // 0xFX55 stores V0 to VX in memory starting at address I
			for i := uint8(0); i <= c.instruction.x; i++ {
				c.memory[c.i+uint16(i)] = c.v[i]
			}
			c.pc += 2
		case 0x65: // 0xFX65 fills V0 to VX with values from memory starting at I
			for i := uint8(0); i <= c.instruction.x; i++ {
				c.v[i] = c.memory[c.i+uint16(i)]
			}
			c.pc += 2
		default:
			fmt.Printf("INVALID OPCODE %X\n", c.instruction.opcode)
		}
	default:
		fmt.Printf("INVALID OPCODE %X\n", c.instruction.opcode)
	}
}
