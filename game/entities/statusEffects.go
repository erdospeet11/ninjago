// Package entities provides the definition and implementation of game entities,
// including various status effects that can be applied to players.
package entities

import (
	"math/rand"
)

// statusEffect represents a basic status effect with a name and duration.
type statusEffect struct {
	name     string // The name of the status effect.
	duration int    // The duration of the status effect.
}

// GetName returns the name of the status effect.
func (s *statusEffect) GetName() string {
	return s.name
}

// StatusEffect defines the behavior of a status effect that can be applied to a player.
type StatusEffect interface {
	Update(p *Player) bool // Update updates the status effect's state for the given player.
	GetName() string       // GetName returns the name of the status effect.
}

// skullDeb represents a debuff that can apply one of several negative effects to a player.
type skullDeb struct {
	statusEffect     // The skull debuff inherits the status effect properties.
	debuffType   int // The type of debuff applied to the player.
}

// NewSkullDeb creates a new SkullDebuff status effect with a random debuff type.
func NewSkullDeb() StatusEffect {
	randomType := rand.Intn(4) + 1

	return &skullDeb{
		statusEffect: statusEffect{
			name:     "SkullDebuff",
			duration: 1800,
		},
		debuffType: randomType,
	}
}

// Update applies the SkullDebuff effect to the player, reducing their attributes based on the debuff type.
// It returns true if the debuff has expired, false otherwise.
func (s *skullDeb) Update(p *Player) bool {
	s.duration--
	//log.Println("Debuff type: ", s.debuffType)
	if s.duration <= 0 {
		switch s.debuffType {
		case 1:
			p.speed = 0.6
		case 2:
			p.BombRange = 2
		case 3:
			p.canPlaceBomb = true
		case 4:
			p.autoPlaceBomb = false
		}
		return true
	} else {
		switch s.debuffType {
		case 1:
			p.speed = 0.3
		case 2:
			p.BombRange = 1
		case 3:
			p.canPlaceBomb = false
		case 4:
			p.autoPlaceBomb = true
		}
	}

	return false
}

// skate represents a temporary speed increase status effect for the player.
type skate struct {
	statusEffect // The skate effect inherits the status effect properties.
}

// NewSkate creates a new Skate status effect.
func NewSkate() StatusEffect {
	return &skate{
		statusEffect: statusEffect{
			name:     "Skate",
			duration: 60,
		},
	}
}

// Update applies the Skate effect to the player, increasing their speed temporarily.
// It returns true if the effect has expired, false otherwise.
func (s *skate) Update(p *Player) bool {
	s.duration--
	if s.duration <= 0 {
		p.speed = 0.6
	} else {
		p.speed = 1
	}

	return s.duration <= 0
}

// radiusInc represents a status effect that increases the player's bomb explosion radius.
type radiusInc struct {
	statusEffect // The radius increase effect inherits the status effect properties.
}

// NewRadiusInc creates a new RadiusIncrease status effect.
func NewRadiusInc() StatusEffect {
	return &radiusInc{
		statusEffect: statusEffect{
			name:     "RadiusIncrease",
			duration: 1800,
		},
	}
}

// Update applies the RadiusIncrease effect to the player, increasing their bomb range.
// It returns true if the effect has expired, false otherwise.
func (s *radiusInc) Update(p *Player) bool {
	s.duration--
	if s.duration <= 0 {
		p.BombRange = 2
	} else {
		p.BombRange = 5
	}

	return s.duration <= 0
}

// bombInc represents a status effect that increases the player's bomb count.
type bombInc struct {
	statusEffect // The bomb count increase effect inherits the status effect properties.
}

// NewBombInc creates a new BombCountIncrease status effect.
func NewBombInc() StatusEffect {
	return &bombInc{
		statusEffect: statusEffect{
			name:     "BombCountIncrease",
			duration: 1800,
		},
	}
}

// Update applies the BombCountIncrease effect to the player, increasing their bomb count.
// It returns true if the effect has expired, false otherwise.
func (s *bombInc) Update(p *Player) bool {
	s.duration--
	if s.duration == 1799 {
		p.NumberOfBombs += 1
	}
	return s.duration <= 0
}

// rollerInc represents a status effect that increases the player's speed significantly.
type rollerInc struct {
	statusEffect // The roller increase effect inherits the status effect properties.
}

// NewRollerInc creates a new RollerIncrease status effect.
func NewRollerInc() StatusEffect {
	return &rollerInc{
		statusEffect: statusEffect{
			name:     "RollerIncrease",
			duration: 1800,
		},
	}
}

// Update applies the RollerIncrease effect to the player, significantly increasing their speed.
// It returns true if the effect has expired, false otherwise.
func (r *rollerInc) Update(p *Player) bool {
	r.duration--
	if r.duration <= 0 {
		p.speed = 0.6
	} else {
		p.speed = 1.2
	}

	return r.duration <= 0
}

// obstacleInc represents a status effect that increases the player's ability to place obstacles.
type obstacleInc struct {
	statusEffect      // The obstacle increase effect inherits the status effect properties.
	applied      bool // Indicates if the effect has been applied to the player.
}

// NewObstacleInc creates a new ObstacleIncrease status effect.
func NewObstacleInc() StatusEffect {
	return &obstacleInc{
		statusEffect: statusEffect{
			name:     "ObstacleIncrease",
			duration: 1800,
		},
		applied: false,
	}
}

// Update applies the ObstacleIncrease effect to the player, increasing their number of obstacles.
// It returns true if the effect has expired, false otherwise.
func (o *obstacleInc) Update(p *Player) bool {
	o.duration--
	if o.duration <= 0 {
		return true
	}

	if !o.applied {
		p.NumberOfObstacles += 3
		o.applied = true
	}

	return false
}

// detonatorInc represents a status effect that increases the player's ability to manually detonate bombs.
type detonatorInc struct {
	statusEffect
}

// NewDetonatorInc creates a new DetonatorIncrease status effect.
func NewDetonatorInc() StatusEffect {
	return &detonatorInc{
		statusEffect: statusEffect{
			name:     "DetonatorIncrease",
			duration: 1800,
		},
	}
}

// Update applies the DetonatorIncrease effect to the player.
// It returns true if the effect has expired, false otherwise.
func (d *detonatorInc) Update(p *Player) bool {
	d.duration--
	return d.duration <= 0
}

// ghostInc represents a status effect that increases the player's ability to move through obstacles.
type ghostInc struct {
	statusEffect // The ghost increase effect inherits the status effect properties.
}

// NewGhostInc creates a new GhostIncrease status effect.
func NewGhostInc() StatusEffect {
	return &ghostInc{
		statusEffect: statusEffect{
			name:     "GhostIncrease",
			duration: 1800,
		},
	}
}

// Update applies the GhostIncrease effect to the player.
// It returns true if the effect has expired, false otherwise.
func (g *ghostInc) Update(p *Player) bool {
	g.duration--
	return g.duration <= 0
}

// invincibilityInc represents a status effect that makes the player invincible for a duration.
type invincibilityInc struct {
	statusEffect
}

// NewInvincibilityInc creates a new InvincibilityIncrease status effect.
func NewInvincibilityInc() StatusEffect {
	return &invincibilityInc{
		statusEffect: statusEffect{
			name:     "InvincibilityIncrease",
			duration: 1800,
		},
	}
}

// Update applies the InvincibilityIncrease effect to the player.
// It returns true if the effect has expired, false otherwise.
func (i *invincibilityInc) Update(p *Player) bool {
	i.duration--
	return i.duration <= 0
}
