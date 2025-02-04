//go:build darwin && amd64
// +build darwin,amd64

package mtldriver

import (
	"image"

	"dmitri.shuralyov.com/gpu/mtl"
	"github.com/diakovliev/oak/v4/shiny/driver/internal/drawer"
	"github.com/diakovliev/oak/v4/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

// BGRA on AMD64 is RGBA8UNorm.
// TODO: make it configurable
var platformPixelFormat = mtl.PixelFormatRGBA8UNorm

func (w *Window) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	draw.Draw(w.bgra, sr.Sub(sr.Min).Add(dp), src.RGBA(), sr.Min, draw.Src)
}

func (w *Window) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	draw.NearestNeighbor.Transform(w.bgra, src2dst, src.(*textureImpl).rgba, sr, op, nil)
}

func (w *Window) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, op)
}

type BGRA = image.RGBA

var NewBGRA = image.NewRGBA
