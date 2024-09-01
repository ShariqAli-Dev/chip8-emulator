# Chip-8 Emulator

A CHIP-8 emulator is a program that simulates the behavior of the CHIP-8, an interpreted programming language used in simple 8-bit computers, allowing the execution of classic games and applications designed for the CHIP-8 platform.

## Screenshot

![Hello Chip8](https://user-images.githubusercontent.com/1005086/76698226-abb83c80-66a0-11ea-93bd-32776fee27e4.png)

## Installing

#### Get dependencies

```
go get -u github.com/veandco/go-sdl2/sdl
```

Note: Read the installation instructions for sdl2 library from [`go-sdl2`](github.com/veandco/go-sdl2) for your os platform.

#### Get code

```
go get -u github.com/shariqali-dev/chip8-emulator
```

## Running

Sample [roms](https://github.com/skatiyar/go-chip8/tree/master/roms) can be used to test the installation.

```
go mod tidy
go run ./cmd/chip8 <path/to/rom>
```

or

```
make build
./bin/chip8 ./roms/filter.ch8
```

## Key Bindings

```
Chip8 keypad         Keyboard mapping
1 | 2 | 3 | C        1 | 2 | 3 | 4
4 | 5 | 6 | D   =>   Q | W | E | R
7 | 8 | 9 | E   =>   A | S | D | F
A | 0 | B | F        Z | X | C | V
```

## Sources

- [Abstracted Guide](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)
- [Golang SDL2 Bindings](https://github.com/veandco/go-sdl2)
- [Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#0.1)
- [Chip-8 opcode table](https://en.wikipedia.org/wiki/CHIP-8)
