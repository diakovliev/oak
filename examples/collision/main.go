package main

import (
	"image/color"
	"time"

	"github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/alg/floatgeom"
	"github.com/diakovliev/oak/v4/collision"
	"github.com/diakovliev/oak/v4/entities"
	"github.com/diakovliev/oak/v4/event"
	"github.com/diakovliev/oak/v4/key"
	"github.com/diakovliev/oak/v4/render"
	"github.com/diakovliev/oak/v4/scene"
	"github.com/diakovliev/oak/v4/shake"
)

const (
	_                   = iota
	RED collision.Label = iota
	GREEN
	BLUE
	TEAL
)

// if true, shake the screen on certain collisions
var demoShake bool = true

func main() {
	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {
		act := &AttachCollisionTest{}
		act.CallerID = ctx.Register(act)
		act.Entity = entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(50, 50, 50, 50)),
			entities.WithColor(color.RGBA{0, 0, 0, 255}),
			entities.WithDrawLayers([]int{0, 1}),
			entities.WithParent(act),
		)

		event.Bind(ctx, event.Enter, act, func(act *AttachCollisionTest, ev event.EnterPayload) event.Response {
			if act.ShouldUpdate {
				act.ShouldUpdate = false
				act.Renderable.Undraw()
				act.Renderable = act.nextR
				render.Draw(act.Renderable, 0, 1)
			}
			if oak.IsDown(key.A) {
				act.ShiftX(-3)
			} else if oak.IsDown(key.D) {
				act.ShiftX(3)
			}
			if oak.IsDown(key.W) {
				act.ShiftY(-3)
			} else if oak.IsDown(key.S) {
				act.ShiftY(3)
			}
			return 0
		})

		collision.PhaseCollision(act.Space, ctx.CollisionTree)

		commonOpts := entities.And(
			entities.WithDrawLayers([]int{0, 0}),
			entities.WithDimensions(floatgeom.Point2{320, 240}),
		)

		upLeft := entities.New(ctx, commonOpts,
			entities.WithColor(color.RGBA{100, 0, 0, 100}),
			entities.WithLabel(RED),
		)

		upRight := entities.New(ctx, commonOpts,
			entities.WithPosition(floatgeom.Point2{320, 0}),
			entities.WithColor(color.RGBA{0, 100, 0, 100}),
			entities.WithLabel(GREEN),
		)
		_ = upRight

		botLeft := entities.New(ctx, commonOpts,
			entities.WithPosition(floatgeom.Point2{0, 240}),
			entities.WithColor(color.RGBA{0, 0, 100, 100}),
			entities.WithLabel(BLUE),
		)

		botRight := entities.New(ctx, commonOpts,
			entities.WithPosition(floatgeom.Point2{320, 240}),
			entities.WithColor(color.RGBA{0, 100, 100, 100}),
			entities.WithLabel(TEAL),
		)

		event.Bind(ctx, collision.Start, act, func(act *AttachCollisionTest, l collision.Label) event.Response {
			switch l {
			case RED:
				act.r += 125
				act.UpdateR()
			case GREEN:
				act.g += 125
				act.UpdateR()
				if demoShake {
					shake.DefaultShaker.Shake(upLeft, time.Second)
					shake.DefaultShaker.Shake(botLeft, time.Second)
					shake.DefaultShaker.Shake(botRight, time.Second)
				}
			case BLUE:
				act.b += 125
				act.UpdateR()
				if demoShake {
					shake.DefaultShaker.Shake(act, time.Second*2)
				}
			case TEAL:
				act.b += 125
				act.g += 125
				act.UpdateR()
				if demoShake {
					shake.DefaultShaker.ShakeScreen(ctx, time.Second)
				}
			}
			return 0
		})
		event.Bind(ctx, collision.Stop, act, func(act *AttachCollisionTest, l collision.Label) event.Response {
			switch l {
			case RED:
				act.r -= 125
				act.UpdateR()
			case GREEN:
				act.g -= 125
				act.UpdateR()
			case BLUE:
				act.b -= 125
				act.UpdateR()
			case TEAL:
				act.b -= 125
				act.g -= 125
				act.UpdateR()
			}
			return 0
		})

	}})
	render.SetDrawStack(
		render.NewDynamicHeap(),
	)
	oak.Init("demo")
}

type AttachCollisionTest struct {
	*entities.Entity
	event.CallerID
	r, g, b      int
	ShouldUpdate bool
	nextR        render.Renderable
}

func (act *AttachCollisionTest) CID() event.CallerID {
	return act.CallerID.CID()
}

// UpdateR with the rgb set on the act.
func (act *AttachCollisionTest) UpdateR() {
	act.nextR = render.NewColorBox(50, 50, color.RGBA{uint8(act.r), uint8(act.g), uint8(act.b), 255})
	act.nextR.SetPos(act.X(), act.Y())
	act.nextR.SetLayer(1)
	act.ShouldUpdate = true
}
