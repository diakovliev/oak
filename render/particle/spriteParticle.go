package particle

import (
	"image/draw"

	"github.com/diakovliev/oak/v4/render"
	"github.com/diakovliev/oak/v4/render/mod"
)

// A SpriteParticle is a particle that has an amount of sprite rotation
type SpriteParticle struct {
	*baseParticle
	rotation float32
}

// Draw redirects to DrawOffsetGen
func (sp *SpriteParticle) Draw(buff draw.Image, xOff, yOff float64) {
	sp.DrawOffsetGen(sp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (sp *SpriteParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	sp.rotation += sp.rotation
	gen := generator.(*SpriteGenerator)
	rgba := gen.Base.Copy().Modify(mod.Rotate(sp.rotation)).GetRGBA()
	render.DrawImage(buff, rgba, int(sp.X()+xOff), int(sp.Y()+yOff))
}
