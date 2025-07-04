// Package scenes provides the implementation of various game scenes,
// including loading levels from text files.
package scenes

import (
	"bufio"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/entities"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// LoadLevelFromTextFile loads a game level from a specified text file.
// It initializes the game scene with static entities, players, status effects, and monsters.
//
// Parameters:
//   - filepath: The path to the text file containing the level data.
//
// Returns:
//   - *GameScene: The initialized game scene.
func LoadLevelFromTextFile(filepath string) *GameScene {
	s := &GameScene{collisionSpace: collider.NewSpatialHash(16)}
	s.staticEntities = make([][]entities.Terrain, 17)

	textLevel, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer textLevel.Close()

	scanner := bufio.NewScanner(textLevel)
	rowIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if rowIndex < 17 {
			s.staticEntities[rowIndex] = make([]entities.Terrain, len(parts)) // Initialize the column size for current row

			for colIndex, char := range parts {
				xPos := float64(colIndex * 16)
				yPos := float64(rowIndex * 16)

				if char == "SOLID" {
					s.staticEntities[rowIndex][colIndex] = entities.NewWall(s.collisionSpace, xPos, yPos, "assets/map/solid_block.png")
				} else {
					s.staticEntities[rowIndex][colIndex] = entities.NewGrass(s.collisionSpace, xPos, yPos, "assets/map/grass_block.png")
				}
			}
		} else {
			if len(parts) < 2 {
				log.Printf("Invalid entity format: %s", line)
				continue
			}

			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				log.Printf("Invalid coordinates in line: %s", line)
				continue
			}

			xPos := float64(x * 16)
			yPos := float64(y * 16)

			switch parts[0] {
			case "PLAYER":
				userData := &userinfo.UserInfo{
					Username: "Player",
				}
				s.players = []entities.Player{*entities.NewPlayer(s.collisionSpace, xPos, xPos, userData, color.Opaque)}
			case "RADIUSINC":
				s.statusEffects = append(s.statusEffects, entities.NewRadiusEffect(s.collisionSpace, xPos, yPos))
			case "SKULLDEB":
				s.statusEffects = append(s.statusEffects, entities.NewSkullDebuff(s.collisionSpace, xPos, yPos))
			case "BOMBINC":
				s.statusEffects = append(s.statusEffects, entities.NewBombEffect(s.collisionSpace, xPos, yPos))
			case "INVINC":
				s.statusEffects = append(s.statusEffects, entities.NewInvincibilityEffect(s.collisionSpace, xPos, yPos))
			case "DETONATOR":
				s.statusEffects = append(s.statusEffects, entities.NewDetonatorEffect(s.collisionSpace, xPos, yPos))
			case "ROLLER":
				s.statusEffects = append(s.statusEffects, entities.NewRollerEffect(s.collisionSpace, xPos, yPos))
			case "GHOSTINC":
				s.statusEffects = append(s.statusEffects, entities.NewGhostEffect(s.collisionSpace, xPos, yPos))
			case "OBSTACLE":
				s.statusEffects = append(s.statusEffects, entities.NewObstacleEffect(s.collisionSpace, xPos, yPos))
			case "GHOST":
				ghost := entities.NewGhost(s.collisionSpace, xPos, yPos, 16*17, 16*17)
				s.monsters = append(s.monsters, ghost)
			case "SLIME":
				slime := entities.NewSlime(s.collisionSpace, xPos, yPos)
				s.monsters = append(s.monsters, slime)
			case "BALLOON":
				ballon := entities.NewBallon(s.collisionSpace, xPos, yPos)
				s.monsters = append(s.monsters, ballon)
			case "ONION":
				onion := entities.NewOnion(s.collisionSpace, xPos, yPos)
				s.monsters = append(s.monsters, onion)
			case "BOX":
				s.boxes = append(s.boxes, *entities.NewBox(s.collisionSpace, xPos, yPos, false))
			default:
				log.Printf("Unknown entity type: %s", parts[0])
			}
		}
		rowIndex++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file: ", err)
	}

	return s
}
