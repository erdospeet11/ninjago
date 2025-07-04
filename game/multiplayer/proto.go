package multiplayer

import (
	"image/color"
	"time"

	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/entities"
)

// ProtoPlayer represents the player information to be shared across the network.
type ProtoPlayer struct {
	Username string                  // The username of the player.
	Ping     time.Duration           // The ping of the player.
	X, Y     float64                 // The x and y coordinates of the player.
	Control  entities.PlayerControls // The player controls.
	Color    color.RGBA              // The color associated with the player.
	IsDead   bool                    // Whether the player is dead.
	Score    int                     // The score of the player.
}

// ProtoEntity represents a generic game entity to be shared across the network.
type ProtoEntity struct {
	X, Y             float64 // The x and y coordinates of the entity.
	IsDead           bool    // Whether the entity is dead.
	Type             string  // The type of the entity.
	CurrentAnimation string  // The current animation of the entity.
}

// ProtoTerrainChange represents a change in the terrain to be shared across the network.
type ProtoTerrainChange struct {
	X, Y int    // The x and y coordinates of the terrain change.
	To   string // The new terrain type.
}

// ProtoGameInfo represents the overall game state information to be shared across the network.
type ProtoGameInfo struct {
	GameState string // The current game state.
	Level     string // The level of the game.

	TerrainChanges []ProtoTerrainChange // The terrain changes in the game.
	Players        []ProtoPlayer        // The players in the game.
	Monsters       []ProtoEntity        // The monsters in the game.
	Bombs          []ProtoEntity        // The bombs in the game.
	Explosions     []ProtoEntity        // The explosions in the game.
	Boxes          []ProtoEntity        // The boxes in the game.
	StatusEffects  []ProtoEntity        // The status effects in the game.
}

// Constants representing the possible game states.
const GameStateLobby = "lobby"
const GameStateRunning = "running"
const GameStateEnd = "end"
