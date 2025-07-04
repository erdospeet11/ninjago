// Package scenes provides the implementation of various game scenes,
// including multiplayer game scenes for joining and managing multiplayer games.
package scenes

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/entities"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/multiplayer"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// MultiPlayerGameSceneJoin represents a scene for joining a multiplayer game.
// It includes the game client, game scene, and state for terrain changes.
type MultiPlayerGameSceneJoin struct {
	Client *multiplayer.GameClient
	GameScene
	terrainChange int
}

// NewMultiPlayerGameSceneJoin creates a new scene for joining a multiplayer game.
// It initializes the game scene based on the provided game client.
//
// Parameters:
//   - client: The game client used to connect to the multiplayer game.
//
// Returns:
//   - Scene: The initialized multiplayer game join scene.
func NewMultiPlayerGameSceneJoin(client *multiplayer.GameClient) Scene {
	s := MultiPlayerGameSceneJoin{
		Client: client,
	}
	s.GameScene = *LoadLevelFromTextFile(client.GameInfo.Level)

	s.screenHeight = 16 * len(s.staticEntities)
	s.screenWidth = 16 * len(s.staticEntities[0])

	return &s
}

// Update updates the multiplayer game scene based on the game state.
// It handles player controls, monster updates, effects, explosions, bombs, and terrain changes.
//
// Parameters:
//   - state: The current game state.
//
// Returns:
//   - error: An error if the update fails.
func (s *MultiPlayerGameSceneJoin) Update(state *GameState) error {
	if s.Client.GameInfo.GameState == multiplayer.GameStateEnd && state.Input.IsAbilityOneJustPressed() {
		state.SceneManager.GoTo(NewMainMenuScene(s.screenWidth, s.screenHeight))
	}

	// Update player controls
	s.Client.Player.Control.Up = state.Input.StateForUp() > 0
	s.Client.Player.Control.Down = state.Input.StateForDown() > 0
	s.Client.Player.Control.Left = state.Input.StateForLeft() > 0
	s.Client.Player.Control.Right = state.Input.StateForRight() > 0
	s.Client.Player.Control.Ability1 = state.Input.IsAbilityOneJustPressed()
	s.Client.Player.Control.Ability2 = state.Input.IsAbilityTwoJustPressed()

	// update monsters
	for i, entity := range s.monsters {
		if i < len(s.Client.GameInfo.Monsters) {
			entity.GetCollider().MoveTo(s.Client.GameInfo.Monsters[i].X, s.Client.GameInfo.Monsters[i].Y)
		}
	}

	//update effects - only rebuild if there are actual changes
	if len(s.Client.GameInfo.StatusEffects) != len(s.statusEffects) {
		// Clear existing effects from collision space
		for _, effect := range s.statusEffects {
			s.collisionSpace.Remove(effect.GetCollider())
		}

		// Rebuild effects list
		s.statusEffects = []entities.Effect{}
		for _, effect := range s.Client.GameInfo.StatusEffects {
			switch effect.Type {
			case "SkullDebuff":
				s.statusEffects = append(s.statusEffects, entities.NewSkullDebuff(s.collisionSpace, effect.X, effect.Y))
			case "BombCountIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewBombEffect(s.collisionSpace, effect.X, effect.Y))
			case "RadiusIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewRadiusEffect(s.collisionSpace, effect.X, effect.Y))
			case "RollerIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewRollerEffect(s.collisionSpace, effect.X, effect.Y))
			case "ObstacleIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewObstacleEffect(s.collisionSpace, effect.X, effect.Y))
			case "DetonatorIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewDetonatorEffect(s.collisionSpace, effect.X, effect.Y))
			case "GhostIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewGhostEffect(s.collisionSpace, effect.X, effect.Y))
			case "InvincibilityIncrease":
				s.statusEffects = append(s.statusEffects, entities.NewInvincibilityEffect(s.collisionSpace, effect.X, effect.Y))
			}
		}
	} else {
		// If count is the same, just update positions to keep effects stable
		for i, effect := range s.statusEffects {
			if i < len(s.Client.GameInfo.StatusEffects) {
				serverEffect := s.Client.GameInfo.StatusEffects[i]
				currentPos := effect.GetCollider().GetPosition()
				// Only update position if there's a significant change to avoid visual jitter
				if math.Abs(currentPos.X-serverEffect.X) > 0.5 || math.Abs(currentPos.Y-serverEffect.Y) > 0.5 {
					effect.GetCollider().MoveTo(serverEffect.X, serverEffect.Y)
				}
			}
		}
	}

	// Update status effects for animation
	for i := range s.statusEffects {
		s.statusEffects[i].Update()
	}

	for i, v := range s.Client.GameInfo.Players {
		if !v.IsDead && i < len(s.players) {
			s.players[i].ColorOverLay = v.Color
		}
	}

	for i := len(s.explosions); i < len(s.Client.GameInfo.Explosions); i++ {
		s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, s.Client.GameInfo.Explosions[i].X, s.Client.GameInfo.Explosions[i].Y))
	}

	// Iterate backwards to safely remove elements during iteration
	for i := len(s.explosions) - 1; i >= 0; i-- {
		if s.explosions[i].Update() {
			s.collisionSpace.Remove(s.explosions[i].GetCollider())
			if len(s.explosions) > 1 {
				s.explosions[i] = s.explosions[len(s.explosions)-1]
				s.explosions = s.explosions[:len(s.explosions)-1]
			} else {
				s.explosions = make([]*entities.Explosion, 0)
			}
		}
	}

	for i := len(s.bombs); i < len(s.Client.GameInfo.Bombs); i++ {
		s.bombs = append(s.bombs, entities.NewBomb(s.collisionSpace, nil, 1, s.Client.GameInfo.Bombs[i].X, s.Client.GameInfo.Bombs[i].Y))
	}

	// Iterate backwards to safely remove elements during iteration
	for i := len(s.bombs) - 1; i >= 0; i-- {
		if s.bombs[i].Update() {
			s.collisionSpace.Remove(s.bombs[i].GetCollider())
			if len(s.bombs) > 1 {
				s.bombs[i] = s.bombs[len(s.bombs)-1]
				s.bombs = s.bombs[:len(s.bombs)-1]
			} else {
				s.bombs = make([]*entities.Bomb, 0)
			}
		}
	}

	if len(s.players) != len(s.Client.GameInfo.Players) {
		s.players = make([]entities.Player, len(s.Client.GameInfo.Players))
		for i, v := range s.Client.GameInfo.Players {
			s.players[i] = *entities.NewPlayer(s.collisionSpace, v.X, v.Y, &userinfo.UserInfo{Username: v.Username}, v.Color)
		}
	}

	for s.terrainChange < len(s.Client.GameInfo.TerrainChanges) {
		terrainChange := s.Client.GameInfo.TerrainChanges[s.terrainChange]

		// Check bounds to prevent index out of range
		if terrainChange.X >= 0 && terrainChange.X < len(s.staticEntities) &&
			terrainChange.Y >= 0 && len(s.staticEntities) > 0 && terrainChange.Y < len(s.staticEntities[0]) {
			s.staticEntities[terrainChange.X][terrainChange.Y] = entities.NewGrass(s.collisionSpace, float64(terrainChange.X*16), float64(terrainChange.Y*16), "assets/map/grass_block.png")
		}
		s.terrainChange++
	}

	return nil
}

// Draw renders the multiplayer game scene onto the screen.
//
// Parameters:
//   - screen: The image to which the scene is drawn.
func (s *MultiPlayerGameSceneJoin) Draw(screen *ebiten.Image) {
	s.GameScene.Draw(screen)

	for i, v := range s.Client.GameInfo.Players {
		if v.IsDead {
			continue
		}

		if i < len(s.players) {
			s.players[i].GetCollider().MoveTo(v.X, v.Y)
			s.players[i].Draw(screen)
		}
	}

	if s.Client.GameInfo.GameState == multiplayer.GameStateEnd {
		drawLogo(screen, 400, 30, "Game Over")

		OrientationY := 50
		for _, player := range s.Client.GameInfo.Players {
			drawLogo(screen, 400, OrientationY, fmt.Sprintln(player.Username, player.Score))
			OrientationY += 30
		}
	}
}
