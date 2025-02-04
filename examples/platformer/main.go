package main

import (
	"image/color"
	"math"

	"github.com/diakovliev/oak/v4/alg/floatgeom"

	"github.com/diakovliev/oak/v4/collision"

	"github.com/diakovliev/oak/v4/event"
	"github.com/diakovliev/oak/v4/key"

	oak "github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/entities"
	"github.com/diakovliev/oak/v4/scene"
)

const (
	// The only collision label we need for this demo is 'ground',
	// indicating something we shouldn't be able to fall or walk through
	Ground collision.Label = 1
)

func main() {
	oak.AddScene("platformer", scene.Scene{Start: func(ctx *scene.Context) {

		char := entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(100, 100, 16, 32)),
			entities.WithColor(color.RGBA{255, 0, 0, 255}),
			entities.WithSpeed(floatgeom.Point2{3, 7}),
		)

		const fallSpeed = .2

		event.Bind(ctx, event.Enter, char, func(c *entities.Entity, ev event.EnterPayload) event.Response {

			// Move left and right with A and D
			if oak.IsDown(key.A) {
				char.Delta[0] = -char.Speed.X()
			} else if oak.IsDown(key.D) {
				char.Delta[0] = char.Speed.X()
			} else {
				char.Delta[0] = (0)
			}
			oldX, oldY := char.X(), char.Y()
			char.ShiftDelta()
			aboveGround := false

			hit := collision.HitLabel(char.Space, Ground)

			// If we've moved in y value this frame and in the last frame,
			// we were below what we're trying to hit, we are still falling
			if hit != nil && !(oldY != char.Y() && oldY+char.H() > hit.Y()) {
				// Correct our y if we started falling into the ground
				char.SetY(hit.Y() - char.H())
				// Stop falling
				char.Delta[1] = 0
				// Jump with Space when on the ground
				if oak.IsDown(key.Spacebar) {
					char.Delta[1] -= char.Speed.Y()
				}
				aboveGround = true
			} else {
				//Restart when is below ground
				if char.Y() > 500 {
					char.Delta[1] = 0
					char.SetY(100)
					char.SetX(100)

				}

				// Fall if there's no ground
				char.Delta[1] += fallSpeed
			}

			if hit != nil {
				// If we walked into a piece of ground, move back
				xover, yover := char.Space.Overlap(hit)
				// We, perhaps unintuitively, need to check the Y overlap, not
				// the x overlap
				// if the y overlap exceeds a superficial value, that suggests
				// we're in a state like
				//
				// G = Ground, C = Character
				//
				// GG C
				// GG C
				//
				// moving to the left
				if math.Abs(yover) > 1 {
					// We add a buffer so this doesn't retrigger immediately
					xbump := 1.0
					if xover > 0 {
						xbump = -1
					}
					char.SetX(oldX + xbump)
					if char.Delta.Y() < 0 {
						char.Delta[1] = 0
					}
				}

				// If we're below what we hit and we have significant xoverlap, by contrast,
				// then we're about to jump from below into the ground, and we
				// should stop the character.
				if !aboveGround && math.Abs(xover) > 1 {
					// We add a buffer so this doesn't retrigger immediately
					char.SetY(oldY + 1)
					char.Delta[1] = fallSpeed
				}

			}

			return 0
		})

		platforms := []floatgeom.Rect2{
			floatgeom.NewRect2WH(0, 400, 300, 20),
			floatgeom.NewRect2WH(100, 250, 30, 20),
			floatgeom.NewRect2WH(340, 300, 100, 20),
		}

		for _, p := range platforms {
			entities.New(ctx,
				entities.WithRect(p),
				entities.WithColor(color.RGBA{0, 0, 255, 255}),
				entities.WithLabel(Ground),
			)
		}

	}})
	oak.Init("platformer")
}
