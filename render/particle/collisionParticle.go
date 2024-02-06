package particle

import (
	"image/draw"

	"github.com/diakovliev/oak/v4/collision"
)

// A CollisionParticle is a wrapper around other particles that also
// has a collision space and can functionally react with the environment
// on collision
type CollisionParticle struct {
	Particle
	s *collision.ReactiveSpace
}

// Draw redirects to DrawOffsetGen
func (cp *CollisionParticle) Draw(buff draw.Image, xOff, yOff float64) {
	cp.DrawOffsetGen(cp.Particle.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (cp *CollisionParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*CollisionGenerator)
	cp.Particle.DrawOffsetGen(gen.Generator, buff, xOff, yOff)
}

// Cycle updates the collision particles variables once per rotation
func (cp *CollisionParticle) Cycle(generator Generator) {
	gen := generator.(*CollisionGenerator)
	pos := cp.Particle.GetPos()
	cp.s.Space.Location = collision.NewRect(pos.X(), pos.Y(), cp.s.GetW(), cp.s.GetH())

	hitFlag := <-cp.s.CallOnHits()
	if gen.Fragile && hitFlag {
		cp.Particle.GetBaseParticle().Life = 0
	}
}

// GetDims returns the dimensions of the space of the particle
func (cp *CollisionParticle) GetDims() (int, int) {
	return int(cp.s.GetW()), int(cp.s.GetH())
}
