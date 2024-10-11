package main

import (
    "context"
	"testing"

	pb "github.com/Mitanshk01/DS_HW4/Q1/protofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRevelioServer struct {
	mock.Mock
	pb.LabyrinthService_RevelioServer
}

type Position struct {
	PositionX uint32
	PositionY uint32
}

func (m *mockRevelioServer) Send(pos *pb.Position) error {
	args := m.Called(pos)
	return args.Error(0)
}

type mockBombardaServer struct {
	mock.Mock
	pb.LabyrinthService_BombardaServer
}

func (m *mockBombardaServer) Recv() (*pb.BombardaRequest, error) {
	args := m.Called()
	return args.Get(0).(*pb.BombardaRequest), args.Error(1)
}

func (m *mockBombardaServer) SendAndClose(*pb.EmptyMessage) error {
	args := m.Called()
	return args.Error(0)
}

func TestGetLabyrinthInfo(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "WALL"}},
			{Tile{Type: "COIN"}, Tile{Type: "EMPTY"}},
		},
	}

	response, err := server.GetLabyrinthInfo(context.Background(), &pb.EmptyMessage{})

	assert.NoError(t, err)
	assert.Equal(t, uint32(2), response.Width)
	assert.Equal(t, uint32(2), response.Height)
}

func TestRevelio(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "WALL"}, Tile{Type: "COIN"}},
			{Tile{Type: "COIN"}, Tile{Type: "EMPTY"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "COIN"}, Tile{Type: "EMPTY"}},
		},
		player: Player{Spells: 1},
	}

	mockStream := new(mockRevelioServer)
	mockStream.On("Send", mock.Anything).Return(nil)

	err := server.Revelio(&pb.RevelioRequest{
		TargetPosition: &pb.Position{PositionX: 1, PositionY: 1},
		TileType:       "COIN",
	}, mockStream)

	assert.NoError(t, err)
	assert.Equal(t, uint32(0), server.player.Spells)
	mockStream.AssertNumberOfCalls(t, "Send", 3)
}

func TestBombarda(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
		},
		player: Player{Spells: 1},
	}

	mockStream := new(mockBombardaServer)
	mockStream.On("Recv").Return(&pb.BombardaRequest{
		TargetPosition: &pb.Position{PositionX: 1, PositionY: 1},
	}, nil).Times(3)
	mockStream.On("Recv").Return((*pb.BombardaRequest)(nil), assert.AnError).Once()
	mockStream.On("SendAndClose", mock.Anything).Return(nil)

	err := server.Bombarda(mockStream)

	assert.NoError(t, err)
	assert.Equal(t, uint32(0), server.player.Spells)
	assert.Equal(t, "EMPTY", server.labyrinth[1][1].Type)
}

func TestGetPlayerStatus(t *testing.T) {
	labyrinthServer := &LabyrinthServer{
		player: Player{
			PositionX: 2,
			PositionY: 3,
			Health:    5,
			Score:     10,
		},
	}
	server := &PlayerServer{labyrinthServer: labyrinthServer}

	response, err := server.GetPlayerStatus(context.Background(), &pb.EmptyMessage{})

	assert.NoError(t, err)
	assert.Equal(t, uint32(10), response.Score)
	assert.Equal(t, uint32(5), response.HealthPoints)
	assert.Equal(t, uint32(2), response.Position.PositionX)
	assert.Equal(t, uint32(3), response.Position.PositionY)
}

func TestRegisterMove(t *testing.T) {
	testCases := []struct {
		name           string
		initialPos     Position
		direction      string
		expectedResult pb.MoveResult
		expectedPos    Position
	}{
		{"Move Up", Position{1, 1}, "U", pb.MoveResult_SUCCESS, Position{1, 0}},
		{"Move Down", Position{1, 1}, "D", pb.MoveResult_SUCCESS, Position{1, 2}},
		{"Move Left", Position{1, 1}, "L", pb.MoveResult_SUCCESS, Position{0, 1}},
		{"Move Right", Position{1, 1}, "R", pb.MoveResult_SUCCESS, Position{2, 1}},
		{"Hit Wall", Position{0, 0}, "L", pb.MoveResult_FAILURE, Position{0, 0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			labyrinthServer := &LabyrinthServer{
				labyrinth: [][]Tile{
					{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
					{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
					{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
				},
				player: Player{
					PositionX: tc.initialPos.PositionX,
					PositionY: tc.initialPos.PositionY,
					Health:    3,
				},
			}
			server := &PlayerServer{labyrinthServer: labyrinthServer}

			response, err := server.RegisterMove(context.Background(), &pb.MoveRequest{Direction: tc.direction})

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, response.Result)
			assert.Equal(t, tc.expectedPos.PositionX, server.labyrinthServer.player.PositionX)
			assert.Equal(t, tc.expectedPos.PositionY, server.labyrinthServer.player.PositionY)
		})
	}
}

func TestRegisterMoveCollectCoin(t *testing.T) {
	labyrinthServer := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "COIN"}, Tile{Type: "EMPTY"}},
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
		},
		player: Player{
			PositionX: 0,
			PositionY: 0,
			Health:    3,
			Score:     0,
		},
	}
	server := &PlayerServer{labyrinthServer: labyrinthServer}

	response, err := server.RegisterMove(context.Background(), &pb.MoveRequest{Direction: "R"})

	assert.NoError(t, err)
	assert.Equal(t, pb.MoveResult_SUCCESS, response.Result)
	assert.Equal(t, uint32(1), server.labyrinthServer.player.PositionX)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionY)
	assert.Equal(t, uint32(1), server.labyrinthServer.player.Score)
	assert.Equal(t, "EMPTY", server.labyrinthServer.labyrinth[0][1].Type)
}

func TestRegisterMoveVictory(t *testing.T) {
	labyrinthServer := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
		},
		player: Player{
			PositionX: 0,
			PositionY: 1,
			Health:    3,
		},
	}
	server := &PlayerServer{labyrinthServer: labyrinthServer}

	response, err := server.RegisterMove(context.Background(), &pb.MoveRequest{Direction: "R"})

	assert.NoError(t, err)
	assert.Equal(t, pb.MoveResult_VICTORY, response.Result)
	assert.Equal(t, uint32(1), server.labyrinthServer.player.PositionX)
	assert.Equal(t, uint32(1), server.labyrinthServer.player.PositionY)
}

func TestRegisterMovePlayerDead(t *testing.T) {
	labyrinthServer := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
		},
		player: Player{
			PositionX: 0,
			PositionY: 0,
			Health:    0,
		},
	}
	server := &PlayerServer{labyrinthServer: labyrinthServer}

	response, err := server.RegisterMove(context.Background(), &pb.MoveRequest{Direction: "R"})

	assert.NoError(t, err) 
	assert.Equal(t, pb.MoveResult_PLAYER_DEAD, response.Result)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionX)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionY)
}

func TestRevelioNoSpellsLeft(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "WALL"}, Tile{Type: "COIN"}},
			{Tile{Type: "COIN"}, Tile{Type: "EMPTY"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "COIN"}, Tile{Type: "EMPTY"}},
		},
		player: Player{Spells: 0},
	}

	mockStream := new(mockRevelioServer)
	err := server.Revelio(&pb.RevelioRequest{
		TargetPosition: &pb.Position{PositionX: 1, PositionY: 1},
		TileType:       "COIN",
	}, mockStream)

	assert.Error(t, err)
	assert.Equal(t, "No spells remaining", err.Error())
}

func TestRevelioInvalidTargetPosition(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "WALL"}, Tile{Type: "COIN"}},
			{Tile{Type: "COIN"}, Tile{Type: "EMPTY"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "COIN"}, Tile{Type: "EMPTY"}},
		},
		player: Player{Spells: 1},
	}

	mockStream := new(mockRevelioServer)
	err := server.Revelio(&pb.RevelioRequest{
		TargetPosition: &pb.Position{PositionX: 5, PositionY: 5},
		TileType:       "COIN",
	}, mockStream)

	assert.Error(t, err)
	assert.Equal(t, "Invalid target position", err.Error())
}

func TestBombardaInvalidNumberOfPoints(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
		},
		player: Player{Spells: 1},
	}

	mockStream := new(mockBombardaServer)
	mockStream.On("Recv").Return(&pb.BombardaRequest{
		TargetPosition: &pb.Position{PositionX: 1, PositionY: 1},
	}, nil).Times(2)
	mockStream.On("Recv").Return((*pb.BombardaRequest)(nil), assert.AnError).Once()

	err := server.Bombarda(mockStream)

	assert.Error(t, err)
	assert.Equal(t, "Invalid number of points received, expected exactly 3", err.Error())
}

func TestBombardaNoSpellsLeft(t *testing.T) {
	server := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
			{Tile{Type: "WALL"}, Tile{Type: "WALL"}, Tile{Type: "WALL"}},
		},
		player: Player{Spells: 0},
	}

	mockStream := new(mockBombardaServer)
	err := server.Bombarda(mockStream)

	assert.Error(t, err)
	assert.Equal(t, "No spells remaining", err.Error())
}

func TestRegisterMoveOutOfBounds(t *testing.T) {
	labyrinthServer := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
		},
		player: Player{
			PositionX: 0,
			PositionY: 0,
			Health:    3,
		},
	}
	server := &PlayerServer{labyrinthServer: labyrinthServer}

	response, err := server.RegisterMove(context.Background(), &pb.MoveRequest{Direction: "U"})

	assert.NoError(t, err)
	assert.Equal(t, pb.MoveResult_FAILURE, response.Result)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionX)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionY)
	assert.Equal(t, uint32(2), server.labyrinthServer.player.Health)
}

func TestRegisterMoveHealthDepletion(t *testing.T) {
	labyrinthServer := &LabyrinthServer{
		labyrinth: [][]Tile{
			{Tile{Type: "EMPTY"}, Tile{Type: "WALL"}},
			{Tile{Type: "EMPTY"}, Tile{Type: "EMPTY"}},
		},
		player: Player{
			PositionX: 0,
			PositionY: 0,
			Health:    1,
		},
	}
	server := &PlayerServer{labyrinthServer: labyrinthServer}

	response, err := server.RegisterMove(context.Background(), &pb.MoveRequest{Direction: "R"})

	assert.NoError(t, err)
	assert.Equal(t, pb.MoveResult_PLAYER_DEAD, response.Result)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionX)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.PositionY)
	assert.Equal(t, uint32(0), server.labyrinthServer.player.Health)
}
