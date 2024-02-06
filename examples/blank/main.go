package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/render"
	"github.com/diakovliev/oak/v4/scene"
)

// This example is a blank, default scene with a pprof server. Useful for
// benchmarks and as a base to copy a starting point from.

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	oak.AddScene("blank", scene.Scene{
		Start: func(ctx *scene.Context) {
			ctx.DrawStack.Draw(render.NewDrawFPS(0, nil, 10, 10))
			ctx.DrawStack.Draw(render.NewLogicFPS(0, nil, 10, 20))
		},
	})
	oak.Init("blank")
}
