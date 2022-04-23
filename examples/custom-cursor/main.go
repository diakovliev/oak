package main

import (
	"fmt"
	"image/color"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func main() {
	oak.AddScene("customcursor", scene.Scene{
		Start: func(ctx *scene.Context) {
			err := ctx.Window.HideCursor()
			if err != nil {
				fmt.Println(err)
			}

			box := render.NewSequence(15,
				render.NewColorBox(2, 2, color.RGBA{255, 255, 0, 255}),
				render.NewColorBox(3, 3, color.RGBA{255, 235, 0, 255}),
				render.NewColorBox(4, 4, color.RGBA{255, 215, 0, 255}),
				render.NewColorBox(5, 5, color.RGBA{255, 195, 0, 255}),
				render.NewColorBox(6, 6, color.RGBA{255, 175, 0, 255}),
				render.NewColorBox(5, 5, color.RGBA{255, 155, 0, 255}),
				render.NewColorBox(4, 4, color.RGBA{255, 135, 0, 255}),
				render.NewColorBox(3, 3, color.RGBA{255, 115, 0, 255}),
				render.NewColorBox(2, 2, color.RGBA{255, 95, 0, 255}),
				render.NewColorBox(1, 1, color.RGBA{255, 75, 0, 255}),
				render.EmptyRenderable(),
				render.EmptyRenderable(),
				render.EmptyRenderable(),
				render.EmptyRenderable(),
			)
			ctx.DrawStack.Draw(box)

			event.GlobalBind(ctx,
				mouse.Drag, func(mouseEvent *mouse.Event) event.Response {
					box.SetPos(mouseEvent.X(), mouseEvent.Y())
					return 0
				})
		},
	})
	oak.Init("customcursor")
}
