package main

import (
	"fmt"
	"time"

	"github.com/diakovliev/oak/v4/debugtools/inputviz"
	"github.com/diakovliev/oak/v4/render"

	"github.com/diakovliev/oak/v4/alg/floatgeom"
	"github.com/diakovliev/oak/v4/event"

	oak "github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/joystick"
	"github.com/diakovliev/oak/v4/scene"
)

func main() {
	oak.AddScene("viz", scene.Scene{Start: func(ctx *scene.Context) {
		joystick.Init()
		latestInput := new(string)
		*latestInput = "Latest Input: Keyboard+Mouse"
		ctx.DrawStack.Draw(render.NewStrPtrText(latestInput, 10, 460), 4)
		ctx.DrawStack.Draw(render.NewText("Space to Vibrate", 10, 440), 4)

		event.GlobalBind(ctx, oak.InputChange, func(input oak.InputType) event.Response {

			switch input {
			case oak.InputJoystick:
				*latestInput = "Latest Input: Joystick"
			case oak.InputKeyboard:
				*latestInput = "Latest Input: Keyboard"
			case oak.InputMouse:
				*latestInput = "Latest Input: Mouse"
			}
			return 0
		})
		go func() {
			rBounds := ctx.Window.Bounds().DivConst(2)
			jCh, cancel := joystick.WaitForJoysticks(1 * time.Second)
			defer cancel()
			for joy := range jCh {
				fmt.Println("new joystick", joy.ID())
				var x, y float64
				switch joy.ID() {
				case 0:
					// 0,0
				case 1:
					x = float64(rBounds.X())
				case 2:
					y = float64(rBounds.Y())
				case 3:
					x = float64(rBounds.X())
					y = float64(rBounds.Y())
				}
				jrend := inputviz.Joystick{
					Rect:          floatgeom.NewRect2WH(x, y, float64(rBounds.X()), float64(rBounds.Y())),
					StickDeadzone: 4000,
					BaseLayer:     -1,
				}
				err := jrend.RenderAndListen(ctx, joy, 1)
				if err != nil {
					fmt.Println("renderer:", err)
				}
			}
		}()
	}})
	oak.Init("viz", func(c oak.Config) (oak.Config, error) {
		c.TrackInputChanges = true
		return c, nil
	})
}
