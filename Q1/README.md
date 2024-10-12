# Homework-4: gRPC

```
Names : Mitansh Kayathwal, Pradeep Mishra
Roll Nos: 2021101026, 2023801013
Branch : CSE, PhD
Course : Distributed Systems, Monsoon '24
```

# **_Directory Structure_**

```
📁 Q1
├── 📁 client
│   └── 📄 main.go
├── 📁 protofiles
│   └── 📄 common.pb.go
│   └── 📄 common.proto
│   └── 📄 labyrinth_grpc.pb.go
│   └── 📄 labyrinth.pb.go
│   └── 📄 labyrinth.proto
│   └── 📄 player_grpc.pb.go
│   └── 📄 player.pb.go
│   └── 📄 player.proto
├── 📁 server
│   └── 📄 labyrinth_game_test.go
│   └── 📄 labyrinth.txt
│   └── 📄 main.go
├── 📄 go.mod
├── 📄 go.sum
├── 📄 Makefile
├── 📄 README.md
```

# Implementation Details

## Overview
This part implements a single player labyrinth game using gRPC for client-server communication. The game allows the player to navigate through a labyrinth, collect treasures, and use spells to reach the exit and win the game.

## Components

### Proto Files
- `common.proto`: Defines common message types used across the game.
- `labyrinth.proto`: Defines the Labyrinth service and related message types.
- `player.proto`: Defines the Player service and related message types.

### Server
- `main.go`: Implements the game server, handling player connections and game logic.
- `labyrinth_game_test.go`: Contains unit tests for the game logic.
- `labyrinth.txt`: Text file containing the labyrinth layout.

### Client
- `main.go`: Implements the game client, allowing players to interact with the game.

## Game Features
- Spell Usage
- Labyrinth navigation
- Treasure collection
- Win condition (reaching the last cell of the labyrinth)

## Running the Code

To run the game, follow these steps:

1. Start the server:
   ```
   $ make server
   ```

2. In a separate terminal, start the client:
   ```
   $ make client
   ```

3. Follow the prompts in the client to play the game.

## Labyrinth File Format

The `labyrinth.txt` file (located in the server directory) should follow this format:

- Each line represents the labyrinth layout:
  - 'W' represents a wall
  - 'C' represents a coin
  - 'E' represents an empty cell
  -  Every cell must be represented by one of the above characters, separated by spaces.

Example (labyrinth.txt):

W W W

W C W

W W W 

This represents a 3x3 labyrinth with a coin at the center and walls surrounding it.

## Makefile Usage

The Makefile provides several commands to simplify development and testing:

- `make clean`: Removes generated files (proto files).
- `make proto`: Generates Go code from the proto files.
- `make server`: Builds and runs the server.
- `make client`: Builds and runs the client.
- `make tests`: Runs the unit tests.

To use these commands, simply run `make <command>` in the terminal (from the root directory).

## Unit Tests

Unit tests are implemented in `server/labyrinth_test.go`. These tests cover various aspects of the game logic, including:

- Player movement
- Treasure collection
- Win condition checking
- Collision detection
- Concurrent player actions
- Edge case scenarios (e.g., moving into walls, invalid moves)

To run the unit tests, use the following command (from the root directory, mentioned above as well):

```
$ make tests
```