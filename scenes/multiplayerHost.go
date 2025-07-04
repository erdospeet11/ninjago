package scenes

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/entities"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/multiplayer"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// MultiPlayerGameSceneHost represents a scene for hosting a multiplayer game.
type MultiPlayerGameSceneHost struct {
	Server    *multiplayer.GameServer // The game server managing the multiplayer game.
	GameScene                         // The game scene containing all necessary entities and game state.
}

// NewMultiPlayerGameSceneHost initializes a new multiplayer game scene for the host.
//
// Parameters:
//   - server: The game server managing the multiplayer game.
//   - level: The path to the level file to load.
//   - host: The user information of the host player.
//
// Returns:
//   - Scene: The initialized multiplayer game scene for the host.
func NewMultiPlayerGameSceneHost(server *multiplayer.GameServer, level string, host *userinfo.UserInfo) Scene {
	s := MultiPlayerGameSceneHost{
		Server: server,
	}
	s.GameScene = *LoadLevelFromTextFile(level)

	if len(s.players) == 0 {
		s.players = append(s.players, *entities.NewPlayer(s.collisionSpace, 1, 1, host, s.Server.GameInfo.Players[0].Color))
	}
	s.Server.GameInfo.GameState = multiplayer.GameStateRunning

	s.screenHeight = 16 * len(s.staticEntities)
	s.screenWidth = 16 * len(s.staticEntities[0])

	s.Server.GameInfo.Monsters = make([]multiplayer.ProtoEntity, len(s.monsters))
	s.Server.GameInfo.Boxes = make([]multiplayer.ProtoEntity, len(s.boxes))
	s.Server.GameInfo.StatusEffects = make([]multiplayer.ProtoEntity, len(s.statusEffects))
	for i, statusEffect := range s.statusEffects {
		s.Server.GameInfo.StatusEffects[i] = multiplayer.ProtoEntity{X: statusEffect.GetCollider().GetPosition().X, Y: statusEffect.GetCollider().GetPosition().Y, Type: statusEffect.StatusEffect.GetName()}
	}

	return &s
}

// Update handles the game state updates for the multiplayer game hosted by this player.
//
// Parameters:
//   - state: The current game state containing input and scene manager.
//
// Returns:
//   - error: An error if the update process fails.
func (s *MultiPlayerGameSceneHost) Update(state *GameState) error {
	if s.Server.GameInfo.GameState == multiplayer.GameStateEnd && state.Input.IsAbilityOneJustPressed() {
		state.SceneManager.GoTo(NewMainMenuScene(s.screenWidth, s.screenHeight))
	}

	// Update player controls
	s.Server.GameInfo.Players[0].Control.Up = state.Input.StateForUp() > 0
	s.Server.GameInfo.Players[0].Control.Down = state.Input.StateForDown() > 0
	s.Server.GameInfo.Players[0].Control.Left = state.Input.StateForLeft() > 0
	s.Server.GameInfo.Players[0].Control.Right = state.Input.StateForRight() > 0
	s.Server.GameInfo.Players[0].Control.Ability1 = state.Input.IsAbilityOneJustPressed()
	s.Server.GameInfo.Players[0].Control.Ability2 = state.Input.IsAbilityTwoJustPressed()

	// update monsters
	for i, entity := range s.monsters {
		s.monsters[i].Update()
		target := rand.Intn(len(s.players))
		s.monsters[i].SetTarget(&s.players[target])
		s.Server.GameInfo.Monsters[i].X = entity.GetCollider().GetPosition().X
		s.Server.GameInfo.Monsters[i].Y = entity.GetCollider().GetPosition().Y
	}

	// update boxes
	for i, box := range s.boxes {
		s.boxes[i].Update()
		if i >= len(s.Server.GameInfo.Boxes) {
			s.Server.GameInfo.Boxes = append(s.Server.GameInfo.Boxes, multiplayer.ProtoEntity{X: box.GetCollider().GetPosition().X, Y: box.GetCollider().GetPosition().Y})
		}
		s.Server.GameInfo.Boxes[i].X = box.GetCollider().GetPosition().X
		s.Server.GameInfo.Boxes[i].Y = box.GetCollider().GetPosition().Y
	}

	// update status effects
	for i, effect := range s.statusEffects {
		s.statusEffects[i].Update()
		if i >= len(s.Server.GameInfo.StatusEffects) {
			s.Server.GameInfo.StatusEffects = append(s.Server.GameInfo.StatusEffects, multiplayer.ProtoEntity{X: effect.GetCollider().GetPosition().X, Y: effect.GetCollider().GetPosition().Y, Type: effect.StatusEffect.GetName()})
		}
		s.Server.GameInfo.StatusEffects[i].X = effect.GetCollider().GetPosition().X
		s.Server.GameInfo.StatusEffects[i].Y = effect.GetCollider().GetPosition().Y
	}

	// update players
	for i, v := range s.Server.GameInfo.Players {
		// remove dead players
		if v.IsDead {
			continue
		}

		s.players[i].Control = v.Control

		s.players[i].ColorOverLay = v.Color

		// add new players
		if len(s.Server.GameInfo.Players) > len(s.players) {
			s.players = append(s.players, *entities.NewPlayer(s.collisionSpace, v.X, v.Y, &userinfo.UserInfo{Username: v.Username}, v.Color))
		}

		newBomb, newBox, err := s.players[i].Update()
		if err != nil {
			log.Println(err)
		}

		if newBomb != nil {
			s.bombs = append(s.bombs, newBomb)
			s.Server.GameInfo.Bombs = append(s.Server.GameInfo.Bombs, multiplayer.ProtoEntity{X: newBomb.GetCollider().GetPosition().X, Y: newBomb.GetCollider().GetPosition().Y})
		}

		if newBox != nil {
			s.boxes = append(s.boxes, *newBox)
			s.Server.GameInfo.Boxes = append(s.Server.GameInfo.Boxes, multiplayer.ProtoEntity{X: newBox.GetCollider().GetPosition().X, Y: newBox.GetCollider().GetPosition().Y})
		}

		playerCollision := s.collisionSpace.CheckCollisions(s.players[i].GetCollider())
		//log.Println(s.players[i].State)
		for _, collision := range playerCollision {
			sep := collision.SeparatingVector
			switch collidingEntity := collision.Other.GetParent().(type) {
			case nil:
				break
			case *entities.Bomb:
				if !(s.players[i].State != nil && (s.players[i].State.GetName() == "GhostIncrease")) {
					s.players[i].GetCollider().Move(sep.X, sep.Y)
				}
			case *entities.Explosion:
				if !(s.players[i].State != nil && (s.players[i].State.GetName() == "InvincibilityIncrease")) {
					s.Server.GameInfo.Players[i].IsDead = true
				}
			case *entities.Box:
				if !(s.players[i].State != nil && (s.players[i].State.GetName() == "GhostIncrease")) {
					s.players[i].GetCollider().Move(sep.X, sep.Y)
				}
			case entities.Effect:
				s.players[i].State = collidingEntity.StatusEffect
				var effectsToRemove []int
				for j, effect := range s.statusEffects {
					if effect.GetCollider() == collision.Other {
						s.players[i].State = effect.StatusEffect
						s.collisionSpace.Remove(effect.GetCollider())
						effectsToRemove = append(effectsToRemove, j)
					}
				}

				// Sort indices in reverse order to prevent index invalidation during removal
				sort.Sort(sort.Reverse(sort.IntSlice(effectsToRemove)))
				for _, index := range effectsToRemove {
					if index < len(s.statusEffects) && index < len(s.Server.GameInfo.StatusEffects) {
						// Safe swap-remove operation
						if index < len(s.statusEffects)-1 {
							s.statusEffects[index] = s.statusEffects[len(s.statusEffects)-1]
							s.Server.GameInfo.StatusEffects[index] = s.Server.GameInfo.StatusEffects[len(s.Server.GameInfo.StatusEffects)-1]
						}
						s.statusEffects = s.statusEffects[:len(s.statusEffects)-1]
						s.Server.GameInfo.StatusEffects = s.Server.GameInfo.StatusEffects[:len(s.Server.GameInfo.StatusEffects)-1]
					}
				}

			case entities.Terrain:
				if collidingEntity.IsSolid() {
					s.players[i].GetCollider().Move(sep.X, sep.Y)
				}
			case entities.Monster:
				if !(s.players[i].State != nil && (s.players[i].State.GetName() == "InvincibilityIncrease")) {
					s.Server.GameInfo.Players[i].IsDead = true
				}
			}
		}
		s.Server.GameInfo.Players[i].X = s.players[i].GetCollider().GetPosition().X
		s.Server.GameInfo.Players[i].Y = s.players[i].GetCollider().GetPosition().Y
	}

	for i := 0; i < len(s.bombs); i++ {
		if s.bombs[i].Update() {
			x, y := s.bombs[i].TilePosition()
			s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64(x*16), float64(y*16)))
			for j := 1; j <= s.bombs[i].ExplosionRange; j++ {
				if y+j < len(s.staticEntities) && x >= 0 && x < len(s.staticEntities) &&
					len(s.staticEntities) > 0 && y+j < len(s.staticEntities[0]) &&
					(s.staticEntities[x][y+j].IsDestroyable() || !s.staticEntities[x][y+j].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64(x*16), float64((y+j)*16)))
				} else {
					break
				}
			}
			for j := 1; j <= s.bombs[i].ExplosionRange; j++ {
				if y-j >= 0 && x >= 0 && x < len(s.staticEntities) &&
					len(s.staticEntities) > 0 && y-j < len(s.staticEntities[0]) &&
					(s.staticEntities[x][y-j].IsDestroyable() || !s.staticEntities[x][y-j].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64(x*16), float64((y-j)*16)))
				} else {
					break
				}
			}
			for j := 1; j <= s.bombs[i].ExplosionRange; j++ {
				if x+j < len(s.staticEntities) && y >= 0 && y < len(s.staticEntities[0]) &&
					len(s.staticEntities) > 0 && x+j < len(s.staticEntities) &&
					(s.staticEntities[x+j][y].IsDestroyable() || !s.staticEntities[x+j][y].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64((x+j)*16), float64((y)*16)))
				} else {
					break
				}
			}
			for j := 1; j <= s.bombs[i].ExplosionRange; j++ {
				if x-j >= 0 && y >= 0 && y < len(s.staticEntities[0]) &&
					len(s.staticEntities) > 0 && x-j < len(s.staticEntities) &&
					(s.staticEntities[x-j][y].IsDestroyable() || !s.staticEntities[x-j][y].IsSolid()) {
					s.explosions = append(s.explosions, entities.NewExplosion(s.collisionSpace, float64((x-j)*16), float64(y*16)))
				} else {
					break
				}
			}
			s.bombs[i].Owner.NumberOfBombs++
			s.collisionSpace.Remove(s.bombs[i].GetCollider())
			if i < len(s.bombs)-1 {
				s.bombs[i] = s.bombs[len(s.bombs)-1]
				s.Server.GameInfo.Bombs[i] = s.Server.GameInfo.Bombs[len(s.Server.GameInfo.Bombs)-1]
			}
			s.bombs = s.bombs[:len(s.bombs)-1]
			s.Server.GameInfo.Bombs = s.Server.GameInfo.Bombs[:len(s.Server.GameInfo.Bombs)-1]
			i-- // Adjust index since we removed an element
		}
	}

	for i := 0; i < len(s.explosions); i++ {
		if s.explosions[i].Update() {
			s.collisionSpace.Remove(s.explosions[i].GetCollider())
			if i < len(s.explosions)-1 {
				s.explosions[i] = s.explosions[len(s.explosions)-1]
				s.Server.GameInfo.Explosions[i] = s.Server.GameInfo.Explosions[len(s.Server.GameInfo.Explosions)-1]
			}
			s.explosions = s.explosions[:len(s.explosions)-1]
			s.Server.GameInfo.Explosions = s.Server.GameInfo.Explosions[:len(s.Server.GameInfo.Explosions)-1]
			i-- // Adjust index since we removed an element
		} else {
			if i >= len(s.Server.GameInfo.Explosions) {
				s.Server.GameInfo.Explosions = append(s.Server.GameInfo.Explosions, multiplayer.ProtoEntity{X: s.explosions[i].GetCollider().GetPosition().X, Y: s.explosions[i].GetCollider().GetPosition().Y})
			}
			s.Server.GameInfo.Explosions[i].X = s.explosions[i].GetCollider().GetPosition().X
			s.Server.GameInfo.Explosions[i].Y = s.explosions[i].GetCollider().GetPosition().Y

			exploX, ExploY := s.explosions[i].TilePosition()
			for j, box := range s.boxes {
				boxX, boxY := box.TilePosition()
				if exploX == boxX && ExploY == boxY {
					if !box.IsBlank {
						newEffect := box.DropRandomStatusEffect()
						if newEffect.StatusEffect != nil {
							s.statusEffects = append(s.statusEffects, newEffect)
							s.Server.GameInfo.StatusEffects = append(s.Server.GameInfo.StatusEffects, multiplayer.ProtoEntity{X: box.GetCollider().GetPosition().X, Y: box.GetCollider().GetPosition().Y, Type: newEffect.StatusEffect.GetName()})
						}
					}
					s.collisionSpace.Remove(box.GetCollider())
					s.Server.GameInfo.Boxes[j] = s.Server.GameInfo.Boxes[len(s.Server.GameInfo.Boxes)-1]
					s.Server.GameInfo.Boxes = s.Server.GameInfo.Boxes[:len(s.Server.GameInfo.Boxes)-1]
					s.boxes[j] = s.boxes[len(s.boxes)-1]
					s.boxes = s.boxes[:len(s.boxes)-1]
				}
			}

			// Check bounds before accessing staticEntities
			if exploX >= 0 && exploX < len(s.staticEntities) &&
				ExploY >= 0 && len(s.staticEntities) > 0 && ExploY < len(s.staticEntities[0]) {
				if s.staticEntities[exploX][ExploY].IsDestroyable() {
					log.Println("Destroyable at x:", exploX, "y:", ExploY, "type:", s.staticEntities[exploX][ExploY])
					s.collisionSpace.Remove(s.staticEntities[exploX][ExploY].GetCollider())
					s.staticEntities[exploX][ExploY] = entities.NewGrass(s.collisionSpace, float64(exploX*16), float64(ExploY*16), "assets/map/grass_block.png")
					s.Server.GameInfo.TerrainChanges = append(s.Server.GameInfo.TerrainChanges, multiplayer.ProtoTerrainChange{X: exploX, Y: ExploY, To: "GRASS"})
				}
			}
		}
	}

	s.Server.GameInfo.GameState = multiplayer.GameStateEnd
	for _, player := range s.Server.GameInfo.Players {
		if !player.IsDead {
			s.Server.GameInfo.GameState = multiplayer.GameStateRunning
			break
		}
	}

	return nil
}

// Draw renders the multiplayer game scene onto the provided screen image.
//
// Parameters:
//   - screen: The image to which the scene is drawn.
func (s *MultiPlayerGameSceneHost) Draw(screen *ebiten.Image) {
	s.GameScene.Draw(screen)

	for i, player := range s.Server.GameInfo.Players {
		if !player.IsDead && i < len(s.players) {
			s.players[i].Draw(screen)
		}
	}

	if s.Server.GameInfo.GameState == multiplayer.GameStateEnd {
		drawLogo(screen, 400, 30, "Game Over")

		OrientationY := 50
		for _, player := range s.Server.GameInfo.Players {
			drawLogo(screen, 400, OrientationY, fmt.Sprintln(player.Username, "Score:", player.Score))
			OrientationY += 30
		}
	}
}
