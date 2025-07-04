// Package entities provides the definition and implementation of game entities
// and their interactions within the game world, including various types of monsters.
package entities

import (
	"math"
	"math/rand"
	"time"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

// snapToGrid snaps the given coordinates to the nearest 16x16 grid position
// to prevent monsters from drifting off-grid due to continuous movement.
// Only snaps if the entity is very close to a grid position (within 0.3 pixels)
func snapToGrid(x, y float64) (float64, float64) {
	gridX := math.Round(x/16.0) * 16.0
	gridY := math.Round(y/16.0) * 16.0

	// Only snap if we're very close to the grid position
	if math.Abs(x-gridX) < 0.3 && math.Abs(y-gridY) < 0.3 {
		return gridX, gridY
	}

	return x, y
}

// Monster defines the behavior of a monster entity in the game.
// It extends the Entity interface with additional methods for setting a target and updating the monster's state.
type Monster interface {
	Entity                    // The monster entity inherits the entity properties.
	SetTarget(target *Player) // SetTarget sets the target player for the monster to follow.
	Update() error            // Update updates the monster's state.
}

// NewGhost creates a new ghost monster entity at the specified position.
// It initializes the ghost with a circular collider and walk animations for different directions.
func NewGhost(collisionSpace *collider.SpatialHash, x, y float64, gameWidth, gameHeight float64) Monster {
	spritePaths := map[string]string{
		"walkLeft":  "assets/enemy/ghost/ghost-left.png",
		"walkRight": "assets/enemy/ghost/ghost-right.png",
		"walkDown":  "assets/enemy/ghost/ghost-down.png",
		"walkUp":    "assets/enemy/ghost/ghost-up.png",
	}

	m := &Ghost{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(x, y, 7),
			speed:            0.6,
			currentAnimation: "walkDown",
			animations:       LoadAnimations(1, 1, spritePaths),
		},
		gameWidth:  gameWidth,
		gameHeight: gameHeight,
	}
	m.collider.SetParent((Monster)(m))

	return m
}

// NewBallon creates a new balloon monster entity at the specified position.
// It initializes the balloon with a circular collider and walk animations for different directions.
func NewBallon(collisionSpace *collider.SpatialHash, x, y float64) Monster {
	spritePaths := map[string]string{
		"walkLeft":  "assets/enemy/balloon/balloon-left.png",
		"walkRight": "assets/enemy/balloon/balloon-right.png",
		"walkDown":  "assets/enemy/balloon/balloon-down.png",
		"walkUp":    "assets/enemy/balloon/balloon-up.png",
	}

	m := &Ballon{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(x, y, 7),
			speed:            0.6,
			currentAnimation: "walkDown",
			animations:       LoadAnimations(1, 1, spritePaths),
		},
		direction:      rand.Intn(4),
		collisionSpace: collisionSpace,
	}
	m.collider.SetParent((Monster)(m))

	return m
}

// NewSlime creates a new slime monster entity at the specified position.
// It initializes the slime with a circular collider and walk animations for different directions.
func NewSlime(collisionSpace *collider.SpatialHash, x, y float64) Monster {
	spritePaths := map[string]string{
		"walkLeft":  "assets/enemy/slime/slime-left.png",
		"walkRight": "assets/enemy/slime/slime-right.png",
		"walkDown":  "assets/enemy/slime/slime-down.png",
		"walkUp":    "assets/enemy/slime/slime-up.png",
	}

	m := &Slime{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(x, y, 7),
			speed:            0.6,
			currentAnimation: "walkDown",
			animations:       LoadAnimations(1, 1, spritePaths),
		},
		direction:      rand.Intn(4),
		collisionSpace: collisionSpace,
	}
	m.collider.SetParent((Monster)(m))

	return m
}

// NewOnion creates a new onion monster entity at the specified position.
// It initializes the onion with a circular collider and walk animations for different directions.
func NewOnion(collisionSpace *collider.SpatialHash, x, y float64) Monster {
	spritePaths := map[string]string{
		"walkLeft":  "assets/enemy/onion/onion-left.png",
		"walkRight": "assets/enemy/onion/onion-right.png",
		"walkDown":  "assets/enemy/onion/onion-down.png",
		"walkUp":    "assets/enemy/onion/onion-up.png",
	}

	m := &Onion{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(x, y, 7),
			speed:            0.6,
			currentAnimation: "walkDown",
			animations:       LoadAnimations(1, 1, spritePaths),
		},
		direction:      rand.Intn(4),
		collisionSpace: collisionSpace,
	}
	m.collider.SetParent((Monster)(m))

	return m
}

// Ghost represents a ghost monster in the game, which moves randomly within the game area.
type Ghost struct {
	entity
	direction           int       // The current direction of the ghost.
	gameWidth           float64   // The width of the game area.
	gameHeight          float64   // The height of the game area.
	lastDirectionChange time.Time // The time of the last direction change.
}

// SetTarget sets the target player for the ghost.
func (g *Ghost) SetTarget(target *Player) {
	return
}

// Update updates the ghost's state, changing its direction randomly and moving it accordingly.
// It handles collisions with other entities and the boundaries of the game area.
func (g *Ghost) Update() error {
	if time.Since(g.lastDirectionChange) > time.Second/2 {
		g.direction = rand.Intn(4)
		g.lastDirectionChange = time.Now()
	}

	switch g.direction {
	case 0:
		g.collider.Move(0, -g.entity.speed)
	case 1:
		g.collider.Move(0, g.entity.speed)
	case 2:
		g.collider.Move(-g.entity.speed, 0)
	case 3:
		g.collider.Move(g.entity.speed, 0)
	}

	// Snap to grid to prevent drift from continuous movement
	currentX, currentY := g.GetCollider().GetPosition().X, g.GetCollider().GetPosition().Y
	gridX, gridY := snapToGrid(currentX, currentY)
	g.GetCollider().MoveTo(gridX, gridY)

	ghostCollision := g.collider.GetHash().CheckCollisions(g.collider)

	for _, collision := range ghostCollision {
		sep := collision.SeparatingVector

		switch collidingEntity := collision.Other.GetParent().(type) {
		case nil:
			break
		case *Bomb:
			g.GetCollider().Move(sep.X, sep.Y)
			g.direction = rand.Intn(4)
		case *terrain:
			if collidingEntity.IsSolid() {
				continue
			}
		}
	}

	// Check if the ghost is out of bounds
	x, y := g.GetCollider().GetPosition().X, g.GetCollider().GetPosition().Y
	if x <= 0 {
		g.GetCollider().Move(1, 0)
	} else if x >= g.gameWidth {
		g.GetCollider().Move(-1, 0)
	}
	if y <= 0 {
		g.GetCollider().Move(0, 1)
	} else if y >= g.gameHeight {
		g.GetCollider().Move(0, -1)
	}

	return nil
}

// Ballon represents a balloon monster in the game, which moves randomly within the game area.
type Ballon struct {
	entity                               // The balloon entity inherits the entity properties.
	direction      int                   // The current direction of the balloon.
	collisionSpace *collider.SpatialHash // The spatial hash for collision detection.
	target         *Player               // The target player for the balloon to follow.
}

// SetTarget sets the target player for the balloon.
func (b *Ballon) SetTarget(target *Player) {
	b.target = target
}

// Update updates the balloon's state, moving it randomly and handling collisions with other entities.
func (b *Ballon) Update() error {
	b.entity.Update()

	switch b.direction {
	case 0: // Up
		b.GetCollider().Move(0, -b.speed)
	case 1: // Down
		b.GetCollider().Move(0, b.speed)
	case 2: // Left
		b.GetCollider().Move(-b.speed, 0)
	case 3: // Right
		b.GetCollider().Move(b.speed, 0)
	}

	// Snap to grid to prevent drift from continuous movement
	currentX, currentY := b.GetCollider().GetPosition().X, b.GetCollider().GetPosition().Y
	gridX, gridY := snapToGrid(currentX, currentY)
	b.GetCollider().MoveTo(gridX, gridY)

	// Check for collisions
	balloonCollision := b.collider.GetHash().CheckCollisions(b.collider)

	// Handle collisions
	for _, collision := range balloonCollision {
		sep := collision.SeparatingVector

		// Check the type of entity we collided with
		switch collidingEntity := collision.Other.GetParent().(type) {
		case nil:
			break
		case *Bomb, *Box:
			b.GetCollider().Move(sep.X, sep.Y)
			// Generate a new direction
			b.direction = rand.Intn(4)
		case *terrain:
			if collidingEntity.IsSolid() {
				b.GetCollider().Move(sep.X, sep.Y)
				monsterX, monsterY := b.GetCollider().GetPosition().X, b.GetCollider().GetPosition().Y
				targetX, targetY := b.target.GetCollider().GetPosition().X, b.target.GetCollider().GetPosition().Y
				b.direction = directionTowardsTarget(monsterX, monsterY, targetX, targetY)
			}
		}
	}

	return nil
}

// directionTowardsTarget calculates the direction from the monster to the target.
// It returns an integer representing the direction (0: up, 1: down, 2: left, 3: right).
func directionTowardsTarget(monsterX, monsterY, targetX, targetY float64) int {
	dx := targetX - monsterX
	dy := targetY - monsterY

	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			return 3 // Right
		} else {
			return 2 // Left
		}
	} else {
		if dy > 0 {
			return 1 // Down
		} else {
			return 0 // Up
		}
	}
}

// Slime represents a slime monster in the game, which moves randomly within the game area.
type Slime struct {
	entity                                    // The slime entity inherits the entity properties.
	target              *Player               // The target player for the slime to follow.
	direction           int                   // The current direction of the slime.
	lastDirectionChange time.Time             // The time of the last direction change.
	collisionSpace      *collider.SpatialHash // The spatial hash for collision detection.
}

// SetTarget sets the target player for the slime.
func (s *Slime) SetTarget(target *Player) {
	s.target = target
}

// Update updates the slime's state, moving it randomly and handling collisions with other entities.
func (s *Slime) Update() error {
	s.entity.Update()

	switch s.direction {
	case 0: // Up
		s.GetCollider().Move(0, -s.speed)
	case 1: // Down
		s.GetCollider().Move(0, s.speed)
	case 2: // Left
		s.GetCollider().Move(-s.speed, 0)
	case 3: // Right
		s.GetCollider().Move(s.speed, 0)
	}

	// Snap to grid to prevent drift from continuous movement
	currentX, currentY := s.GetCollider().GetPosition().X, s.GetCollider().GetPosition().Y
	gridX, gridY := snapToGrid(currentX, currentY)
	s.GetCollider().MoveTo(gridX, gridY)

	// Check for collisions
	onionCollision := s.collider.GetHash().CheckCollisions(s.collider)

	// Handle collisions
	for _, collision := range onionCollision {
		sep := collision.SeparatingVector

		switch collidingEntity := collision.Other.GetParent().(type) {
		case nil:
			break
		case *Bomb, *Box:
			s.GetCollider().Move(sep.X, sep.Y)
			// Generate a new direction
			s.direction = rand.Intn(4)
			s.lastDirectionChange = time.Now()
		case *terrain:
			if collidingEntity.IsSolid() {
				s.GetCollider().Move(sep.X, sep.Y)
				monsterX, monsterY := s.GetCollider().GetPosition().X, s.GetCollider().GetPosition().Y
				targetX, targetY := s.target.GetCollider().GetPosition().X, s.target.GetCollider().GetPosition().Y
				if rand.Float64() < 0.4 { // 40% chance of making a wrong decision
					s.direction = rand.Intn(4)
				} else {
					s.direction = directionTowardsTarget(monsterX, monsterY, targetX, targetY)
				}
			}
		}
	}
	return nil
}

// Onion represents an onion monster in the game, which moves randomly within the game area.
type Onion struct {
	entity                                    // The onion entity inherits the entity properties.
	direction           int                   // The current direction of the onion.
	lastDirectionChange time.Time             // The time of the last direction change.
	collisionSpace      *collider.SpatialHash // The spatial hash for collision detection.
}

// SetTarget sets the target player for the onion.
func (o *Onion) SetTarget(target *Player) {
	return
}

// Update updates the onion's state, moving it randomly and handling collisions with other entities.
func (o *Onion) Update() error {
	o.entity.Update()

	switch o.direction {
	case 0: // Up
		o.GetCollider().Move(0, -o.speed)
	case 1: // Down
		o.GetCollider().Move(0, o.speed)
	case 2: // Left
		o.GetCollider().Move(-o.speed, 0)
	case 3: // Right
		o.GetCollider().Move(o.speed, 0)
	}

	// Snap to grid to prevent drift from continuous movement
	currentX, currentY := o.GetCollider().GetPosition().X, o.GetCollider().GetPosition().Y
	gridX, gridY := snapToGrid(currentX, currentY)
	o.GetCollider().MoveTo(gridX, gridY)

	// Check for collisions
	onionCollision := o.collider.GetHash().CheckCollisions(o.collider)

	// Handle collisions
	for _, collision := range onionCollision {
		sep := collision.SeparatingVector

		// Check the type of entity we collided with
		switch collidingEntity := collision.Other.GetParent().(type) {
		case nil:
			break
		case *Bomb, *Box:
			o.GetCollider().Move(sep.X, sep.Y)
			// Generate a new direction
			o.direction = rand.Intn(4)
			o.lastDirectionChange = time.Now()
		case *terrain:
			if collidingEntity.IsSolid() {
				o.GetCollider().Move(sep.X, sep.Y)
				// Generate a new direction
				o.direction = rand.Intn(4)
				o.lastDirectionChange = time.Now()
			}
		}
	}

	return nil
}
