build:
	@go build -o ./bin/chip8 ./cmd/chip8
run: build
	@./bin/chip8 ./roms/games/Soccer.ch8
