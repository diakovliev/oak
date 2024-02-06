package particle

import (
	"image/color"
	"testing"

	"github.com/diakovliev/oak/v4/alg/span"
	"github.com/diakovliev/oak/v4/event"
	"github.com/diakovliev/oak/v4/physics"
	"github.com/diakovliev/oak/v4/render"
	"github.com/diakovliev/oak/v4/shape"
)

func TestSource(t *testing.T) {
	g := NewGradientGenerator(
		Rotation(span.NewConstant(1.0)),
		Color(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255},
			color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Color2(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255},
			color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Size(span.NewConstant(5)),
		EndSize(span.NewConstant(10)),
		Shape(shape.Heart),
		Progress(render.HorizontalProgress),
		And(
			NewPerFrame(span.NewConstant(200.0)),
		),
		Pos(20, 20),
		LifeSpan(span.NewConstant(10.0)),
		Limit(2047),
		Angle(span.NewConstant(0.0)),
		Speed(span.NewConstant(0.0)),
		Spread(10, 10),
		Duration(span.NewConstant(10)),
		Gravity(10, 10),
		SpeedDecay(1, 1),
		End(func(_ Particle) {}),
		Layer(func(_ physics.Vector) int { return 0 }),
	)

	src := g.Generate(0)

	ended := false

	src.EndFunc = func() {
		ended = true
	}

	for i := 0; i < 1000; i++ {
		rotateParticles(src, event.EnterPayload{})
	}
	for clearParticles(src, event.EnterPayload{}) != event.ResponseUnbindThisBinding {
	}

	if !ended {
		t.Fatalf("source did not stop after duration was exceeded")
	}

	src.Pause()
	if !src.IsPaused() {
		t.Fatalf("Pause did not pause source")
	}
	src.UnPause()
	if src.IsPaused() {
		t.Fatalf("Unpause did not unpause source")
	}
	x, y := src.Generator.GetPos()
	src.ShiftX(10)
	src.ShiftY(10)
	x2, y2 := src.Generator.GetPos()
	if x2 != x+10 {
		t.Fatalf("x post shift expected %v, got %v", x+10, x2)
	}
	if y2 != y+10 {
		t.Fatalf("y post shift expected %v, got %v", y+10, y2)
	}
	src.SetPos(-20, -30)
	x2, y2 = src.Generator.GetPos()
	if x2 != -20.0 {
		t.Fatalf("setpos did not set x, expected %v got %v", -20, x2)
	}
	if y2 != -30.0 {
		t.Fatalf("setpos did not set y, expected %v got %v", -30, y2)
	}

	var src2 *Source
	src2.Stop()
}
