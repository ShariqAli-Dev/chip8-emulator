package chip8

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func (c *Chip8) EmulateCycle() {
	// fetch and set the current opcode
	highByte := c.memory[c.pc]
	lowByte := c.memory[c.pc+1]
	c.instruction.opcode = (uint16(highByte) << 8) | uint16(lowByte)

	// pre increment program counter for the next opcode
	c.pc += 2

	// DXYN masking
	c.instruction.nnn = c.instruction.opcode & 0x0FFF
	c.instruction.nn = uint8(c.instruction.opcode & 0x00FF)
	c.instruction.n = uint8(c.instruction.opcode & 0x000F)
	c.instruction.x = uint8((c.instruction.opcode >> 8) & 0x0F)
	c.instruction.y = uint8((c.instruction.opcode >> 4) & 0x0F)

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
			c.sp--
			c.pc = c.stack[c.sp]
		default:
			c.debugPrintOpcode("0x0? INVALID OPCODE")
		}
	case 0x1:
		c.debugPrintOpcode(fmt.Sprintf("0x1NNN: jumps to address NNN (0x%X)", c.instruction.nnn))
		c.pc = c.instruction.nnn
	case 0x2:
		c.debugPrintOpcode(fmt.Sprintf("0x2NNN: call subroutine at NNN (0x%X)", c.instruction.nnn))
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.instruction.nnn
	case 0x3:
		c.debugPrintOpcode(fmt.Sprintf("0x3XNN: skips the next instruction if V%X == NN (0x%X)", c.instruction.x, c.instruction.nn))
		if c.v[c.instruction.x] == c.instruction.nn {
			c.pc += 2
		}
	case 0x4:
		c.debugPrintOpcode(fmt.Sprintf("0x4XNN: skips the next instruction if V%X != NN (0x%X)", c.instruction.x, c.instruction.nn))
		if c.v[c.instruction.x] != c.instruction.nn {
			c.pc += 2
		}
	case 0x5:
		c.debugPrintOpcode(fmt.Sprintf("0x5XY0: skips the next instruction if V%X == V%X", c.instruction.x, c.instruction.y))
		if c.v[c.instruction.x] == c.v[c.instruction.y] {
			c.pc += 2
		}
	case 0x6:
		c.debugPrintOpcode(fmt.Sprintf("0x6XNN: sets register V%X to NN (0x%X)", c.instruction.x, c.instruction.nn))
		c.v[c.instruction.x] = c.instruction.nn
	case 0x7:
		c.debugPrintOpcode(fmt.Sprintf("0x7XNN: set register V%X += NN (0x%X)  Result: 0x%X", c.instruction.x, c.instruction.nn, c.v[c.instruction.x]+c.instruction.nn))
		c.v[c.instruction.x] += c.instruction.nn
	case 0x8:
		switch c.instruction.n {
		case 0x0:
			c.debugPrintOpcode("0x8XY0: sets the VX to the value of VY")
			c.v[c.instruction.x] = c.v[c.instruction.y]
		case 0x1:
			c.debugPrintOpcode("0x8XY1: sets VX = VX | VY")
			c.v[c.instruction.x] |= c.v[c.instruction.y]
		case 0x2:
			c.debugPrintOpcode("0x8XY2: sets VX = VX & VY")
			c.v[c.instruction.x] &= c.v[c.instruction.y]
		case 0x3:
			c.debugPrintOpcode("0x8XY3: sets VX = VX ^ VY (xor)")
			c.v[c.instruction.x] ^= c.v[c.instruction.y]
		case 0x4:
			c.debugPrintOpcode("0x8XY4: Add VY to VX. if overflow, VF set to 1, else 0.")
			sum := uint16(c.v[c.instruction.x]) + uint16(c.v[c.instruction.y])
			if sum > 255 {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[c.instruction.x] = uint8(sum)
		case 0x5:
			c.debugPrintOpcode("0x8XY5: VX -= VY. Set VF to 1 if there is not a borrow (result is positive).")
			if c.v[c.instruction.x] >= c.v[c.instruction.y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[c.instruction.x] -= c.v[c.instruction.y]
		case 0x6:
			c.debugPrintOpcode("0x8XY6: Shift VX right by one. Store LSB in VF.")
			c.v[0xF] = c.v[c.instruction.x] & 1
			c.v[c.instruction.x] >>= 1
		case 0x7:
			c.debugPrintOpcode("0x8XY7: VX = VY - VX. VF = 1 if no borrow, else 0.")
			if c.v[c.instruction.y] >= c.v[c.instruction.x] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[c.instruction.x] = c.v[c.instruction.y] - c.v[c.instruction.x]
		case 0xE:
			c.debugPrintOpcode("0x8XYE: Shift VX left by one. Store MSB in VF.")
			c.v[0xF] = (c.v[c.instruction.x] & 0x80) >> 7
			c.v[c.instruction.x] <<= 1
		default:
			c.debugPrintOpcode("0x8XY? INVALID OPCODE")
		}
	case 0x9:
		c.debugPrintOpcode(fmt.Sprintf("0x9XY0: skips the next instruction if V%X != V%X", c.instruction.x, c.instruction.y))
		if c.v[c.instruction.x] != c.v[c.instruction.y] {
			c.pc += 2
		}
	case 0xA:
		c.debugPrintOpcode(fmt.Sprintf("0xANNN: set I to NNN (0x%X)", c.instruction.nnn))
		c.i = c.instruction.nnn
	case 0xB:
		c.debugPrintOpcode(fmt.Sprintf("0xBNNN: jumps to the address NNN + V0 (0x%X + %X)", c.instruction.nnn, c.v[0]))
		c.pc = c.instruction.nnn + uint16(c.v[0])
	case 0xC:
		c.debugPrintOpcode("0xCXNN: Sets VX to the result of a bitwise AND operation on a random number (0 to 255) and NN")
		n, err := rand.Int(rand.Reader, big.NewInt(256))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		c.v[c.instruction.x] = uint8(n.Int64()) & c.instruction.nn
	case 0xD:
		c.debugPrintOpcode(fmt.Sprintf("0xDXYN: draws a sprite at coordinate (V%X, V%X) with a height of N (0xD%X)", c.instruction.x, c.instruction.y, c.instruction.n))
		c.v[0xF] = 0
		for j := uint8(0); j < c.instruction.n; j++ {
			pixel := c.memory[c.i+uint16(j)]
			for i := uint8(0); i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					x := (c.v[c.instruction.x] + i) % 64
					y := (c.v[c.instruction.y] + j) % 32
					if c.display[y][x] == 1 {
						c.v[0xF] = 1
					}
					c.display[y][x] ^= 1
				}
			}
		}
		c.shouldDraw = true
	case 0xE:
		switch c.instruction.nn {
		case 0x9E:
			c.debugPrintOpcode("0xEX9E: skips the next instruction if the key stored in VX is pressed")
			if c.Keypad[c.v[c.instruction.x]] {
				c.pc += 2
			}
		case 0xA1:
			c.debugPrintOpcode("0xEXA1: skips the next instruction if the key stored in VX is not pressed")
			if !c.Keypad[c.v[c.instruction.x]] {
				c.pc += 2
			}
		default:
			c.debugPrintOpcode("0xEX? INVALID OPCODE")
		}
	case 0xF:
		switch c.instruction.nn {
		case 0x07:
			c.debugPrintOpcode("0xFX07: sets VX to the value of the delay timer")
			c.v[c.instruction.x] = c.delayTimer
		case 0x0A:
			c.debugPrintOpcode("0xFX0A: waits for a key press and stores the result in VX")
			keyPressed := false
			for i := 0; i < 16; i++ {
				if c.Keypad[i] {
					c.v[c.instruction.x] = uint8(i)
					keyPressed = true
				}
			}
			if !keyPressed {
				c.pc -= 2
			}
		case 0x15:
			c.debugPrintOpcode("0xFX15: sets the delay timer to VX")
			c.delayTimer = c.v[c.instruction.x]
		case 0x18:
			c.debugPrintOpcode("0xFX18: sets the sound timer to VX")
			c.soundTimer = c.v[c.instruction.x]
		case 0x1E:
			c.debugPrintOpcode("0xFX1E: adds VX to I. Sets VF to 1 if there is a carry, else 0")
			c.i += uint16(c.v[c.instruction.x])
			if c.i > 0xFFF {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
		case 0x29:
			c.debugPrintOpcode("0xFX29: sets I to the location of the sprite for the character in VX")
			c.i = uint16(c.v[c.instruction.x]) * 5
		case 0x33:
			c.debugPrintOpcode("0xFX33: stores the binary-coded decimal representation of VX at addresses I, I+1, and I+2")
			c.memory[c.i] = c.v[c.instruction.x] / 100
			c.memory[c.i+1] = (c.v[c.instruction.x] / 10) % 10
			c.memory[c.i+2] = c.v[c.instruction.x] % 10
		case 0x55:
			c.debugPrintOpcode("0xFX55: stores registers V0 through VX in memory starting at address I")
			for i := 0; i <= int(c.instruction.x); i++ {
				c.memory[c.i+uint16(i)] = c.v[i]
			}
		case 0x65:
			c.debugPrintOpcode("0xFX65: reads registers V0 through VX from memory starting at address I")
			for i := 0; i <= int(c.instruction.x); i++ {
				c.v[i] = c.memory[c.i+uint16(i)]
			}
		default:
			c.debugPrintOpcode("0xF? INVALID OPCODE")
		}
	default:
		c.debugPrintOpcode("INVALID OPCODE")
	}
}

func (c *Chip8) debugPrintOpcode(opcodeDescription string) {
	fmt.Printf("address: 0x0%X  opcode 0x%X: %s\n", c.pc-2, c.instruction.opcode, opcodeDescription)
}
