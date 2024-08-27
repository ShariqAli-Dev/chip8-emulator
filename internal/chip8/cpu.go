package chip8

import "fmt"

func (c *Chip8) EmulateCycle() {
	// // fetch and set the current opcode
	highByte := c.memory[c.pc]
	lowByte := c.memory[c.pc+1]
	c.instruction.opcode = (uint16(highByte) << 8) | uint16(lowByte)

	// pre increment program counter for the next opcode
	c.pc += 2

	// DXYN masking
	c.instruction.nnn = c.instruction.opcode & 0x0FFF
	c.instruction.nn = uint8(c.instruction.opcode & 0x0FF)
	c.instruction.n = uint8(c.instruction.opcode & 0x0F)
	c.instruction.x = uint8(c.instruction.opcode>>8) & 0x0F
	c.instruction.y = uint8((c.instruction.opcode >> 4) & 0x000F)

	// CXNN --> 000C,
	// 000C & 000F --> 0101 & 1111 --> 0101 -> C
	opcodeCategory := (c.instruction.opcode >> 12) & 0x0F

	switch opcodeCategory {
	case 0x0:
		switch c.instruction.nn {
		case 0xE0:
			c.debugPrintOpcode("0x00E0: clears the screen")
			for i := 0; i < len(c.display); i++ {
				for j := 0; j < len(c.display[i]); j++ {
					c.display[i][j] = 0
				}
			}
			c.shouldDraw = true
		case 0xEE:
			c.debugPrintOpcode("0x00EE: returns from a subroutine")
			c.sp-- // "pop" the last value from the stack
			c.pc = c.stack[c.sp]
		default:
			c.debugPrintOpcode("0x0? INVALID OPCODE")
		}
	case 0x1:
		c.debugPrintOpcode("0x1NNN: jumps to address NNN")
		c.pc = c.instruction.nnn
	case 0x2:
		c.debugPrintOpcode("0x2NNN: call subroutine at NNN")
		// store current address to return to on subroutine stack ("push" it on the stack)
		// set the program ocutner to subroitune address so that the next opcode is gotten from there
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.instruction.nnn
	case 0x3:
		c.debugPrintOpcode(fmt.Sprintf("0x3NNN: skips the next instruction if V%X == NN (0x3%X)", c.v[c.instruction.x], c.instruction.nn))
		if c.v[c.instruction.x] == c.instruction.nn {
			c.pc += 2
		}
	case 0x4:
		c.debugPrintOpcode(fmt.Sprintf("0x4NNN: skips the next instruction if V%X != NN (0x3%X)", c.v[c.instruction.x], c.instruction.nn))
		if c.v[c.instruction.x] != c.instruction.nn {
			c.pc += 2
		}
	case 0x5:
		c.debugPrintOpcode(fmt.Sprintf("0x5XYN: skips the next instruction if V%X == V%X", c.v[c.instruction.x], c.v[c.instruction.y]))
		if c.v[c.instruction.x] == c.v[c.instruction.y] {
			c.pc += 2
		}
	case 0x6:
		c.debugPrintOpcode(fmt.Sprintf("0x6XNN: sets register V%X to NN (0x6%X)", c.instruction.x, c.instruction.nn))
		c.v[c.instruction.x] = c.instruction.nn
	case 0x7:
		c.debugPrintOpcode(fmt.Sprintf("0x7XNN: set register V%X += NN (0x7%X)  Result: 0x7%X", c.instruction.x, c.instruction.nn, c.v[c.instruction.x]+c.instruction.nn))
		c.v[c.instruction.x] += c.instruction.nn
	case 0x9:
		c.debugPrintOpcode(fmt.Sprintf("0x9XY0: skips the next instruction if V%X != V%X", c.v[c.instruction.x], c.v[c.instruction.y]))
		if c.v[c.instruction.x] != c.v[c.instruction.y] {
			c.pc += 2
		}
	case 0xA:
		c.debugPrintOpcode(fmt.Sprintf("0xANNN: set I to NNN (0xA%X)", c.instruction.nnn))
		c.i = c.instruction.nnn
	case 0xD:
		c.debugPrintOpcode(fmt.Sprintf("0xDXYN: draws a sprite at coordinate (V%X, V%X) with a height of N (0xD%X)", c.instruction.x, c.instruction.y, c.instruction.n))
		c.v[0xF] = 0 // init carry flag to 0
		// loop over the n rows of sprite
		for j := uint8(0); j < c.instruction.n; j++ {
			// get the next byte/row of sprite data
			pixel := c.memory[c.i+uint16(j)]
			for i := uint16(0); i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if c.display[c.v[c.instruction.y]+uint8(j)][c.v[c.instruction.x]+uint8(i)] == 1 {
						c.v[0xF] = 1 // Set carry flag if pixel is already set
					}
					c.display[c.v[c.instruction.y]+uint8(j)][c.v[c.instruction.x]+uint8(i)] ^= 1
				}
			}
		}
		c.shouldDraw = true
	default:
		c.debugPrintOpcode("UNIMPLEMENTED OPCODE")
	}

	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}
	if c.soundTimer > 0 {
		if c.soundTimer == 1 {
			// play the beepor sound effect
		}
		c.soundTimer -= 1
	}
}

func (c *Chip8) debugPrintOpcode(opcodeDescription string) {
	fmt.Printf("address: 0x0%X  opcode 0x%X: %s\n", c.pc-2, c.instruction.opcode, opcodeDescription)
}
