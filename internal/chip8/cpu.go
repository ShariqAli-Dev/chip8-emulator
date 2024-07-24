package chip8

import (
	"log"
)

// EmulateCycle fetches, decodes, and executes one opcode
func (c *Chip8) EmulateCycle() {
	opcode := c.fetchOpcode()
	c.decodeAndExecuteOpcode(opcode)
	c.updateTimers()
}

// fetchOpcode retrieves the next opcode from memory
func (c *Chip8) fetchOpcode() uint16 {
	if c.programCounter+1 >= uint16(len(c.memory)) {
		log.Fatal("Program counter out of bounds")
	}
	opcode := uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
	c.programCounter += 2 // Increment the program counter by 2 to point to the next instruction
	return opcode
}

// decodeAndExecuteOpcode decodes the fetched opcode and executes the corresponding instruction
func (c *Chip8) decodeAndExecuteOpcode(opcode uint16) {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			c.op00E0()
		case 0x000E:
			c.op00EE()
		default:
			log.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	case 0x1000:
		c.op1NNN(opcode)
	case 0x2000:
		c.op2NNN(opcode)
	case 0x3000:
		c.op3XNN(opcode)
	case 0x4000:
		c.op4XNN(opcode)
	case 0x5000:
		c.op5XY0(opcode)
	case 0x6000:
		c.op6XNN(opcode)
	case 0x7000:
		c.op7XNN(opcode)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			c.op8XY0(opcode)
		case 0x0001:
			c.op8XY1(opcode)
		case 0x0002:
			c.op8XY2(opcode)
		case 0x0003:
			c.op8XY3(opcode)
		case 0x0004:
			c.op8XY4(opcode)
		case 0x0005:
			c.op8XY5(opcode)
		case 0x0006:
			c.op8XY6(opcode)
		case 0x0007:
			c.op8XY7(opcode)
		case 0x000E:
			c.op8XYE(opcode)
		default:
			log.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	case 0x9000:
		c.op9XY0(opcode)
	case 0xA000:
		c.opANNN(opcode)
	case 0xB000:
		c.opBNNN(opcode)
	case 0xC000:
		c.opCXNN(opcode)
	case 0xD000:
		c.opDXYN(opcode)
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			c.opEX9E(opcode)
		case 0x00A1:
			c.opEXA1(opcode)
		default:
			log.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			c.opFX07(opcode)
		case 0x000A:
			c.opFX0A(opcode)
		case 0x0015:
			c.opFX15(opcode)
		case 0x0018:
			c.opFX18(opcode)
		case 0x001E:
			c.opFX1E(opcode)
		case 0x0029:
			c.opFX29(opcode)
		case 0x0033:
			c.opFX33(opcode)
		case 0x0055:
			c.opFX55(opcode)
		case 0x0065:
			c.opFX65(opcode)
		default:
			log.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	default:
		log.Printf("Unknown opcode: 0x%X\n", opcode)
	}
}

// updateTimers updates the delay and sound timers
func (c *Chip8) updateTimers() {
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.soundTimer--
		if c.soundTimer == 1 {
			// Add sound play code here
		}
	}
}
