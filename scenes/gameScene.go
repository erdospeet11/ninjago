// Package scenes provides the implementation of various game scenes and their management.
package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/entities"
)

// GameScene represents a game scene containing all necessary entities and game state.
type GameScene struct {
	screenHeight   int                   // The height of the game screen.
	screenWidth    int                   // The width of the game screen.
	collisionSpace *collider.SpatialHash // The collision space for the game entities.
	staticEntities [][]entities.Terrain  // The static entities in the game scene.
	monsters       []entities.Monster    // The monsters in the game scene.
	bombs          []*entities.Bomb      // The bombs in the game scene.
	explosions     []*entities.Explosion // The explosions in the game scene.
	boxes          []entities.Box        // The boxes in the game scene.
	statusEffects  []entities.Effect     // The status effects in the game scene.
	players        []entities.Player     // The players in the game scene.
}

// NewSinglePlayerGameScene creates a new single-player game scene by loading a level from a text file.
// It initializes the game scene with entities and sets the screen dimensions based on the level data.
//
// Parameters:
//   - filepath: The path to the text file containing the level data.
//
// Returns:
//   - Scene: The initialized single-player game scene.
func NewSinglePlayerGameScene(filepath string) Scene {
	s := LoadLevelFromTextFile(filepath)

	s.screenHeight = 16 * len(s.staticEntities)
	s.screenWidth = 16 * len(s.staticEntities[0])

	return s
}

// Update updates the game scene based on the current game state.
// This handles single-player game logic including player movement, monster updates, and collisions.
//
// Parameters:
//   - state: The current game state.
//
// Returns:
//   - error: An error if the update fails.
func (s *GameScene) Update(state *GameState) error {
	// Check if we have a player to control
	if len(s.players) == 0 {
		return nil
	}

	// Update player controls
	s.players[0].Control.Up = state.Input.StateForUp() > 0
	s.players[0].Control.Down = state.Input.StateForDown() > 0
	s.players[0].Control.Left = state.Input.StateForLeft() > 0
	s.players[0].Control.Right = state.Input.StateForRight() > 0
	s.players[0].Control.Ability1 = state.Input.IsAbilityOneJustPressed()
	s.players[0].Control.Ability2 = state.Input.IsAbilityTwoJustPressed()

	// Update monsters
	for i := range s.monsters {
		s.monsters[i].Update()
		// Set target to player
		if len(s.players) > 0 {
			s.monsters[i].SetTarget(&s.players[0])
		}
	}

	// Update player
	newBomb, newBox, err := s.players[0].Update()
	if err != nil {
		return err
	}

	// Handle new bomb placement
	if newBomb != nil {
		s.bombs = append(s.bombs, newBomb)
	}

	// Handle new box placement
	if newBox != nil {
		s.boxes = append(s.boxes, *newBox)
	}

	// Handle player collisions
	playerCollision := s.collisionSpace.CheckCollisions(s.players[0].GetCollider())
	for _, collision := range playerCollision {
		sep := collision.SeparatingVector
		switch collidingEntity := collision.Other.GetParent().(type) {
		case nil:
			break
		case *entities.Bomb:
			if !(s.players[0].State != nil && (s.players[0].State.GetName() == "GhostIncrease")) {
				s.players[0].GetCollider().Move(sep.X, sep.Y)
			}
		case *entities.Explosion:
			if !(s.players[0].State != nil && (s.players[0].State.GetName() == "InvincibilityIncrease")) {
				// Player died - could restart level or go to game over
				state.SceneManager.GoTo(NewMainMenuScene(s.screenWidth, s.screenHeight))
				return nil
			}
		case *entities.Box:
			if !(s.players[0].State != nil && (s.players[0].State.GetName() == "GhostIncrease")) {
				s.players[0].GetCollider().Move(sep.X, sep.Y)
			}
		case entities.Effect:
			s.players[0].State = collidingEntity.StatusEffect
			// Remove the collected effect
			var effectsToRemove []int
			for j, effect := range s.statusEffects {
				if effect.GetCollider() == collision.Other {
					s.players[0].State = effect.StatusEffect
					s.collisionSpace.Remove(effect.GetCollider())
					effectsToRemove = append(effectsToRemove, j)
				}
			}
			for _, index := range effectsToRemove {
				if index < len(s.statusEffects)-1 {
					s.statusEffects[index] = s.statusEffects[len(s.statusEffects)-1]
				}
				s.statusEffects = s.statusEffects[:len(s.statusEffects)-1]
			}
		case entities.Terrain:
			if collidingEntity.IsSolid() {
				s.players[0].GetCollider().Move(sep.X, sep.Y)
			}
		case entities.Monster:
			if !(s.players[0].State != nil && (s.players[0].State.GetName() == "InvincibilityIncrease")) {
				// Player died - restart level or go to game over
				state.SceneManager.GoTo(NewMainMenuScene(s.screenWidth, s.screenHeight))
				return nil
			}
		}
	}

	// Update bombs and handle explosions
	for i, bomb := range s.bombs {
		if bomb.Update() {
			x, y := bomb.TilePosition()
			s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64(x*16), float64(y*16)))

			// Create explosion pattern
			for j := 1; j <= bomb.ExplosionRange; j++ {
				if y+j < len(s.staticEntities) && (s.staticEntities[x][y+j].IsDestroyable() || !s.staticEntities[x][y+j].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64(x*16), float64((y+j)*16)))
				} else {
					break
				}
			}
			for j := 1; j <= bomb.ExplosionRange; j++ {
				if y-j >= 0 && (s.staticEntities[x][y-j].IsDestroyable() || !s.staticEntities[x][y-j].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64(x*16), float64((y-j)*16)))
				} else {
					break
				}
			}
			for j := 1; j <= bomb.ExplosionRange; j++ {
				if x+j < len(s.staticEntities) && (s.staticEntities[x+j][y].IsDestroyable() || !s.staticEntities[x+j][y].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64((x+j)*16), float64((y)*16)))
				} else {
					break
				}
			}
			for j := 1; j <= bomb.ExplosionRange; j++ {
				if x-j >= 0 && (s.staticEntities[x-j][y].IsDestroyable() || !s.staticEntities[x-j][y].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64((x-j)*16), float64(y*16)))
				} else {
					break
				}
			}

			// Restore bomb count to player
			s.bombs[i].Owner.NumberOfBombs++
			s.collisionSpace.Remove(bomb.GetCollider())

			// Remove bomb from slice
			if i < len(s.bombs)-1 {
				s.bombs[i] = s.bombs[len(s.bombs)-1]
			}
			s.bombs = s.bombs[:len(s.bombs)-1]
		}
	}

	// Update explosions and handle destruction
	for i := 0; i < len(s.explosions); i++ {
		if s.explosions[i].Update() {
			s.collisionSpace.Remove(s.explosions[i].GetCollider())
			// Remove explosion from slice
			if i < len(s.explosions)-1 {
				s.explosions[i] = s.explosions[len(s.explosions)-1]
			}
			s.explosions = s.explosions[:len(s.explosions)-1]
			i-- // Adjust index since we removed an element
		} else {
			// Check for box destruction
			exploX, exploY := s.explosions[i].TilePosition()
			for j, box := range s.boxes {
				boxX, boxY := box.TilePosition()
				if exploX == boxX && exploY == boxY {
					if !box.IsBlank {
						newEffect := box.DropRandomStatusEffect()
						if newEffect.StatusEffect != nil {
							s.statusEffects = append(s.statusEffects, newEffect)
						}
					}
					s.collisionSpace.Remove(box.GetCollider())
					// Remove box from slice
					if j < len(s.boxes)-1 {
						s.boxes[j] = s.boxes[len(s.boxes)-1]
					}
					s.boxes = s.boxes[:len(s.boxes)-1]
				}
			}

			// Check for terrain destruction
			if s.staticEntities[exploX][exploY].IsDestroyable() {
				s.collisionSpace.Remove(s.staticEntities[exploX][exploY].GetCollider())
				s.staticEntities[exploX][exploY] = entities.NewGrass(s.collisionSpace, float64(exploX*16), float64(exploY*16), "assets/map/grass_block.png")
			}
		}
	}

	// Update status effects
	for i := range s.statusEffects {
		s.statusEffects[i].Update()
	}

	// Update boxes
	for i := range s.boxes {
		s.boxes[i].Update()
	}

	return nil
}

// Draw renders the game scene onto the provided screen image.
//
// Parameters:
//   - screen: The image to which the scene is drawn.
func (s *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{148, 247, 252, 1})

	for _, rows := range s.staticEntities {
		for _, entity := range rows {
			entity.Draw(screen)
		}
	}

	for _, effect := range s.statusEffects {
		effect.Draw(screen)
	}

	for _, box := range s.boxes {
		box.Draw(screen)
	}

	for i := range s.bombs {
		s.bombs[i].Draw(screen)
	}

	for _, entity := range s.monsters {
		entity.Draw(screen)
	}

	for i := range s.explosions {
		s.explosions[i].Draw(screen)
	}
}
