// Package entities provides the definition and implementation of game entities
// and their interactions within the game world.
package entities

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// PlayerControls represents the control state for a player, including movement and abilities.
type PlayerControls struct {
	Up       bool // Indicates if the player is moving up.
	Down     bool // Indicates if the player is moving down.
	Left     bool // Indicates if the player is moving left.
	Right    bool // Indicates if the player is moving right.
	Ability1 bool // Indicates if the player is using ability 1.
	Ability2 bool // Indicates if the player is using ability 2.
}

// Player represents a player entity in the game, with various attributes and abilities.
type Player struct {
	entity
	Control             PlayerControls     // The control state of the player.
	userData            *userinfo.UserInfo // The user data associated with the player.
	State               StatusEffect       // The current status effect applied to the player.
	BombRange           int                // The range of the player's bombs.
	NumberOfBombs       int                // The number of bombs the player can place.
	NumberOfObstacles   int                // The number of obstacles the player can place.
	canPlaceBomb        bool               // Indicates if the player can place a bomb.
	autoPlaceBomb       bool               // Indicates if the player should automatically place a bomb.
	manualDetonateBombs []*Bomb            // The list of bombs set to manual detonation.
	ColorOverLay        color.Color        // The color overlay for the player sprite.
}

// NewPlayer creates a new player entity at the specified position with the provided user data and color overlay.
// It initializes the player with a circular collider and walk animations for different directions.
func NewPlayer(colliderSpace *collider.SpatialHash, x, y float64, userData *userinfo.UserInfo, colorOverlay color.Color) *Player {
	spritePaths := map[string]string{
		"walkLeft":  "assets/player/character-left.png",
		"walkRight": "assets/player/character-right.png",
		"walkDown":  "assets/player/character-down.png",
		"walkUp":    "assets/player/character_up.png",
	}
	p := &Player{
		entity: entity{
			collider:         colliderSpace.NewCircleShape(x, y, 7), // 8 is the radius of the player (hitboxRadius
			speed:            0.6,
			currentAnimation: "walkDown",
			animations:       LoadAnimations(1, 1, spritePaths),
		},
		userData:          userData,
		BombRange:         2,
		NumberOfBombs:     1,
		NumberOfObstacles: 0,
		canPlaceBomb:      true,
		ColorOverLay:      colorOverlay,
	}
	p.collider.SetParent(p)

	return p
}

// Update updates the player's state, handling movement, abilities, and status effects.
// It returns the bomb and box placed by the player (if any), and an error if the animation cannot be set.
func (p *Player) Update() (bomb *Bomb, box *Box, err error) {
	//log.Println(p.NumberOfBombs)
	//x, y := p.collider.GetPosition().X, p.collider.GetPosition().Y
	//log.Println("at x:", math.Floor(x), "y:", math.Floor(y), math.Floor(x/16), math.Floor(y/16))

	//log.Println(p.NumberOfObstacles)

	if p.State != nil && p.State.Update(p) {
		p.State = nil
	}

	if p.Control.Ability2 && p.State != nil && p.State.GetName() == "DetonatorIncrease" && len(p.manualDetonateBombs) > 0 {
		for _, b := range p.manualDetonateBombs {
			b.Detonate()
		}
		p.manualDetonateBombs = nil
	}

	if p.Control.Ability2 && p.State != nil && p.State.GetName() == "ObstacleIncrease" && p.NumberOfObstacles > 0 {
		box = p.PlaceObstacle()
	}

	if p.autoPlaceBomb {
		bomb = p.PlaceBomb()
	} else if p.Control.Ability1 && p.canPlaceBomb && p.NumberOfBombs > 0 {
		bomb = p.PlaceBomb()
	}

	if p.Control.Up {
		p.collider.Move(0, -p.speed)
		err = p.SetCurrentAnimation("walkUp")
	}
	if p.Control.Down {
		p.collider.Move(0, p.speed)
		err = p.SetCurrentAnimation("walkDown")
	}
	if p.Control.Left {
		p.collider.Move(-p.speed, 0)
		err = p.SetCurrentAnimation("walkLeft")
	}
	if p.Control.Right {
		p.collider.Move(p.speed, 0)
		err = p.SetCurrentAnimation("walkRight")
	}

	return bomb, box, err
}

// PlaceBomb places a bomb at the player's current position if the player can place a bomb and has bombs available.
// It returns the placed bomb.
func (p *Player) PlaceBomb() (bomb *Bomb) {
	if p.canPlaceBomb && p.NumberOfBombs > 0 {
		x, y := p.TilePosition()
		bomb = NewBomb(p.collider.GetHash(), p, p.BombRange, float64(x*16), float64(y*16))
		if p.State != nil && p.State.GetName() == "DetonatorIncrease" {
			bomb.manualDetonate = true
			p.manualDetonateBombs = append(p.manualDetonateBombs, bomb)
		}
		p.NumberOfBombs--
		return bomb
	}
	return nil
}

// PlaceObstacle places an obstacle at the player's current position if the player has obstacles available.
// It returns the placed obstacle.
func (p *Player) PlaceObstacle() (obstacle *Box) {
	if p.NumberOfObstacles > 0 {
		x, y := p.TilePosition()
		obstacle = NewBox(p.collider.GetHash(), float64(x*16), float64(y*16), true)
		p.NumberOfObstacles--
		return obstacle
	}
	return nil
}

// Draw renders the player on the screen with the current animation and color overlay.
func (p *Player) Draw(screen *ebiten.Image) {
	anim, ok := p.animations[p.currentAnimation]
	if !ok {
		log.Printf("No animation found for key %s", p.currentAnimation)
		return
	}
	position := p.collider.GetPosition()
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(position.X, position.Y)
	r, g, b, a := p.ColorOverLay.RGBA()
	opts.ColorScale.SetR(float32(r))
	opts.ColorScale.SetG(float32(g))
	opts.ColorScale.SetB(float32(b))
	opts.ColorScale.SetA(float32(a))

	screen.DrawImage(anim.sprite, opts)
}
