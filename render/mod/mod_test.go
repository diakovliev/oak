package mod

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/diakovliev/oak/v4/shape"
	"github.com/disintegration/gift"
)

func TestComposedModifications(t *testing.T) {
	modList := []Mod{
		Zoom(2.0, 2.0, 2.0),
		CutFromLeft(2, 2),
	}
	base := setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})
	base = modList[0](base)
	base = modList[1](base)
	chained := setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})

	if reflect.DeepEqual(base, chained) {
		t.Fatalf("expected base and chained rgbas to differ")
	}
	mCombined := And(modList[0], modList[1])
	chained = mCombined(chained)
	if !reflect.DeepEqual(base, chained) {
		t.Fatalf("expected base and chained rgbas to equal after modifications")
	}
}

func TestSafeCompose(t *testing.T) {
	modList := []Mod{
		nil,
		Zoom(2.0, 2.0, 2.0),
		CutFromLeft(2, 2),
		nil,
	}
	base := setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})
	base = modList[1](base)
	base = modList[2](base)
	chained := setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})

	if reflect.DeepEqual(base, chained) {
		t.Fatalf("expected base and chained rgbas to differ")
	}
	mCombined := SafeAnd(modList...)
	chained = mCombined(chained)
	if !reflect.DeepEqual(base, chained) {
		t.Fatalf("expected base and chained rgbas to equal after modifications")
	}

	base = setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})
	modList = []Mod{
		nil,
		nil,
	}
	mCombined = SafeAnd(modList...)
	if !reflect.DeepEqual(base, mCombined(base)) {
		t.Fatalf("expected base and nil modified base to equal")
	}

}

func TestAllModifications(t *testing.T) {
	in := setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})
	type filterCase struct {
		Filter
		*image.RGBA
	}
	filterList := []filterCase{{
		ConformToPalette(color.Palette{color.RGBA{64, 0, 0, 128}}),
		setAll(newrgba(3, 3), color.RGBA{64, 0, 0, 128}),
	}, {
		Fade(10),
		setAll(newrgba(3, 3), color.RGBA{245, 0, 0, 245}),
	}, {
		Fade(500),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}),
	}, {
		FillMask(*setAll(newrgba(3, 3), color.RGBA{0, 255, 0, 255})),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		AndFilter(Fade(500), FillMask(*setAll(newrgba(3, 3), color.RGBA{0, 255, 0, 255}))),
		setAll(newrgba(3, 3), color.RGBA{0, 255, 0, 255}),
	}, {
		AndFilter(Fade(500), ApplyColor(color.RGBA{255, 255, 255, 255})),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}),
	}, {
		ApplyColor(color.RGBA{255, 255, 255, 255}),
		setAll(newrgba(3, 3), color.RGBA{127, 127, 127, 255}),
	}, {
		ApplyMask(*setAll(newrgba(3, 3), color.RGBA{255, 255, 255, 255})),
		setAll(newrgba(3, 3), color.RGBA{127, 127, 127, 255}),
	}, {
		AndFilter(Fade(500), ApplyMask(*setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}))),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}),
	}, {
		Brighten(-100),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 255}),
	}, {
		Saturate(-100),
		setAll(newrgba(3, 3), color.RGBA{128, 128, 128, 255}),
	}, {
		ColorBalance(0, 0, 0),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		InPlace(Scale(2, 2)),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		StripOuterAlpha(setOne(setOne(setOne(setOne(setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
			color.RGBA{0, 1, 0, 122}, 0, 0),
			color.RGBA{0, 1, 0, 122}, 1, 0),
			color.RGBA{0, 1, 0, 122}, 0, 1),
			color.RGBA{0, 1, 0, 240}, 2, 2), 200),
		setOne(setOne(setOne(setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
			color.RGBA{0, 1, 0, 0}, 0, 0),
			color.RGBA{0, 1, 0, 0}, 1, 0),
			color.RGBA{0, 1, 0, 0}, 0, 1),
	},
	}
	for _, f := range filterList {
		in2 := copyrgba(in)
		f.Filter(in2)
		if !reflect.DeepEqual(in2, f.RGBA) {
			t.Fatalf("filtered did not match expected")
		}
	}
	type modCase struct {
		Mod
		*image.RGBA
	}
	modList := []modCase{{
		TrimColor(color.RGBA{255, 0, 0, 255}),
		newrgba(0, 0),
	}, {
		TrimColor(color.RGBA{0, 0, 0, 0}),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		Zoom(1.0, 1.0, 1.0),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		Cut(1, 1),
		setAll(newrgba(1, 1), color.RGBA{255, 0, 0, 255}),
	}, {
		CutRel(.66, .66),
		setAll(newrgba(2, 2), color.RGBA{255, 0, 0, 255}),
	}, {
		CutFromLeft(2, 2),
		setAll(newrgba(2, 2), color.RGBA{255, 0, 0, 255}),
	}, {
		CutRound(.5, .5),
		setOne(setOne(setOne(setOne(setOne(
			setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
			color.RGBA{0, 0, 0, 0}, 0, 0),
			color.RGBA{0, 0, 0, 0}, 1, 0),
			color.RGBA{0, 0, 0, 0}, 2, 0),
			color.RGBA{0, 0, 0, 0}, 0, 1),
			color.RGBA{0, 0, 0, 0}, 0, 2),
	}, {
		Crop(image.Rect(0, 0, 1, 1)),
		setAll(newrgba(1, 1), color.RGBA{255, 0, 0, 255}),
	}, {
		CropToSize(1, 1, gift.TopLeftAnchor),
		setAll(newrgba(1, 1), color.RGBA{255, 0, 0, 255}),
	}, {
		CutShape(shape.Checkered),
		checker(setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})),
	}}

	for _, m := range modList {
		if !reflect.DeepEqual(m.Mod(in), m.RGBA) {
			t.Fatalf("modified did not match expected")
		}
	}
}

// test utils

func copyrgba(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	rgba2 := newrgba(bounds.Max.X, bounds.Max.Y)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rgba2.Set(x, y, rgba.At(x, y))
		}
	}
	return rgba2
}

func newrgba(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

func setAll(rgba *image.RGBA, c color.Color) *image.RGBA {
	bounds := rgba.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rgba.Set(x, y, c)
		}
	}
	return rgba
}

func setOne(rgba *image.RGBA, c color.Color, x, y int) *image.RGBA {
	rgba.Set(x, y, c)
	return rgba
}

func checker(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if (x+y)%2 == 1 {
				rgba.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	return rgba
}
