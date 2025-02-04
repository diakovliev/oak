package main

import (
	"fmt"
	"image/color"

	"github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/event"
	"github.com/diakovliev/oak/v4/mouse"
	"github.com/diakovliev/oak/v4/render"
	"github.com/diakovliev/oak/v4/scene"
)

func main() {
	c1 := oak.NewWindow()
	c1.DrawStack = render.NewDrawStack(render.NewDynamicHeap())

	// Two windows cannot share the same logic handler
	c1.SetLogicHandler(event.NewBus(nil))
	c1.FirstSceneInput = color.RGBA{255, 0, 0, 255}
	c1.AddScene("scene1", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 1")
			cb := render.NewColorBox(50, 50, ctx.SceneInput.(color.RGBA))
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
			dFPS := render.NewDrawFPS(0.1, nil, 600, 10)
			ctx.DrawStack.Draw(dFPS, 1)
			event.GlobalBind(ctx, mouse.Press, func(me *mouse.Event) event.Response {
				cb.SetPos(me.X(), me.Y())
				return 0
			})
		},
	})
	go func() {
		c1.Init("scene1", func(c oak.Config) (oak.Config, error) {
			c.Debug.Level = "VERBOSE"
			c.DrawFrameRate = 1200
			c.FrameRate = 60
			c.EnableDebugConsole = true
			return c, nil
		})
		fmt.Println("scene 1 exited")
	}()

	c2 := oak.NewWindow()
	c2.DrawStack = render.NewDrawStack(render.NewDynamicHeap())
	c2.SetLogicHandler(event.NewBus(nil))
	c2.FirstSceneInput = color.RGBA{0, 255, 0, 255}
	c2.AddScene("scene2", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 2")
			cb := render.NewColorBox(50, 50, ctx.SceneInput.(color.RGBA))
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
			dFPS := render.NewDrawFPS(0.1, nil, 600, 10)
			ctx.DrawStack.Draw(dFPS, 1)
			event.GlobalBind(ctx, mouse.Press, func(me *mouse.Event) event.Response {
				cb.SetPos(me.X(), me.Y())
				return 0
			})
		},
	})
	c2.Init("scene2", func(c oak.Config) (oak.Config, error) {
		c.Debug.Level = "VERBOSE"
		c.DrawFrameRate = 1200
		c.FrameRate = 60
		c.EnableDebugConsole = true
		return c, nil
	})
	fmt.Println("scene 2 exited")
}
