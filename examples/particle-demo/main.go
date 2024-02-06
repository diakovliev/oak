package main

import (
	"errors"
	"image/color"
	"log"
	"strconv"

	oak "github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/alg"
	"github.com/diakovliev/oak/v4/alg/span"
	"github.com/diakovliev/oak/v4/debugstream"
	"github.com/diakovliev/oak/v4/event"
	"github.com/diakovliev/oak/v4/mouse"
	"github.com/diakovliev/oak/v4/oakerr"
	"github.com/diakovliev/oak/v4/physics"
	"github.com/diakovliev/oak/v4/render"
	pt "github.com/diakovliev/oak/v4/render/particle"
	"github.com/diakovliev/oak/v4/scene"
	"github.com/diakovliev/oak/v4/shape"
)

var (
	startColor     color.Color
	startColorRand color.Color
	endColor       color.Color
	endColorRand   color.Color
	src            *pt.Source
)

func parseShape(args []string) shape.Shape {
	if len(args) > 0 {
		switch args[0] {
		case "heart":
			return shape.Heart
		case "square":
			return shape.Square
		case "circle":
			return shape.Circle
		case "diamond":
			return shape.Diamond
		case "checkered":
			return shape.Checkered
		case "or":
			return shape.JustIn(shape.OrIn(parseShape(args[1:2]).In, parseShape(args[2:]).In))
		case "and":
			return shape.JustIn(shape.AndIn(parseShape(args[1:2]).In, parseShape(args[2:]).In))
		case "not":
			return shape.JustIn(shape.NotIn(parseShape(args[1:]).In))
		}
	}
	return nil
}

func main() {

	debugstream.AddCommand(debugstream.Command{Name: "followMouse", Operation: func(args []string) string {
		event.GlobalBind(event.DefaultBus, event.Enter, func(ev event.EnterPayload) event.Response {
			// It'd be interesting to attach to the mouse position
			src.SetPos(float64(mouse.LastEvent.X()), float64(mouse.LastEvent.Y()))
			return 0
		})
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "shape", Operation: func(args []string) string {
		if len(args) > 0 {
			sh := parseShape(args)
			if sh != nil {
				src.Generator.(pt.Shapeable).SetShape(sh)
			}
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "size", Operation: func(args []string) string {
		f1, f2, two, err := parseInts(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.(pt.Sizeable).SetSize(span.NewConstant(f1))
		} else {
			src.Generator.(pt.Sizeable).SetSize(span.NewLinear(f1, f2))
		}

		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "endsize", Operation: func(args []string) string {
		f1, f2, two, err := parseInts(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.(pt.Sizeable).SetEndSize(span.NewConstant(f1))
		} else {
			src.Generator.(pt.Sizeable).SetEndSize(span.NewLinear(f1, f2))
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "count", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.GetBaseGenerator().NewPerFrame = span.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().NewPerFrame = span.NewLinear(npf, npf2)
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "life", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.GetBaseGenerator().LifeSpan = span.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().LifeSpan = span.NewLinear(npf, npf2)
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "rotation", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.GetBaseGenerator().Rotation = span.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Rotation = span.NewLinear(npf, npf2)
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "angle", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.GetBaseGenerator().Angle = span.NewConstant(npf * alg.DegToRad)
		} else {
			src.Generator.GetBaseGenerator().Angle = span.NewLinear(npf*alg.DegToRad, npf2*alg.DegToRad)
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "speed", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			src.Generator.GetBaseGenerator().Speed = span.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Speed = span.NewLinear(npf, npf2)
		}
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "spread", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "speeds"}.Error()
		}
		src.Generator.GetBaseGenerator().Spread.SetPos(npf, npf2)
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "gravity", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "speeds"}.Error()
		}
		src.Generator.GetBaseGenerator().Gravity.SetPos(npf, npf2)
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "speeddecay", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "speeds"}.Error()
		}
		src.Generator.GetBaseGenerator().SpeedDecay.SetPos(npf, npf2)
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "pos", Operation: func(args []string) string {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "positions"}.Error()
		}
		src.Generator.SetPos(npf, npf2)

		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "startcolor", Operation: func(args []string) string {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}.Error()
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		startColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "startrand", Operation: func(args []string) string {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}.Error()
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		startColorRand = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "endcolor", Operation: func(args []string) string {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}.Error()
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		endColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetEndColor(endColor, endColorRand)
		return ""
	}})

	debugstream.AddCommand(debugstream.Command{Name: "endrand", Operation: func(args []string) string {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}.Error()
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}.Error()
		}
		endColorRand = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetEndColor(endColor, endColorRand)
		return ""
	}})

	oak.AddScene("demo", scene.Scene{Start: func(*scene.Context) {
		render.Draw(render.NewDrawFPS(0, nil, 10, 10))
		x := 320.0
		y := 240.0
		newPf := span.NewLinear(1.0, 2.0)
		life := span.NewLinear(100.0, 120.0)
		angle := span.NewLinear(0.0, 360.0)
		speed := span.NewLinear(1.0, 5.0)
		size := span.NewConstant(1)
		layerFn := func(v physics.Vector) int {
			return 1
		}
		startColor = color.RGBA{255, 255, 255, 255}
		startColorRand = color.RGBA{0, 0, 0, 0}
		endColor = color.RGBA{255, 255, 255, 255}
		endColorRand = color.RGBA{0, 0, 0, 0}
		shape := shape.Square

		src = pt.NewColorGenerator(
			pt.Pos(x, y),
			pt.Duration(pt.Inf),
			pt.LifeSpan(life),
			pt.Angle(angle),
			pt.Speed(speed),
			pt.Layer(layerFn),
			pt.Shape(shape),
			pt.Size(size),
			pt.Color(startColor, startColorRand, endColor, endColorRand),
			pt.NewPerFrame(newPf)).Generate(0)
	}})

	render.SetDrawStack(
		render.NewCompositeR(),
	)

	err := oak.Init("demo", func(c oak.Config) (oak.Config, error) {
		c.Debug.Level = "VERBOSE"
		c.DrawFrameRate = 1200
		c.FrameRate = 60
		c.Title = "Particle Demo"
		c.EnableDebugConsole = true
		return c, nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func parseRGBA(args []string) (r, g, b, a int, err error) {
	if len(args) < 4 {
		return
	}
	r, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	g, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	b, err = strconv.Atoi(args[2])
	if err != nil {
		return
	}
	a, err = strconv.Atoi(args[3])
	return
}

func parseFloats(args []string) (f1, f2 float64, two bool, err error) {
	if len(args) < 1 {
		err = errors.New("no args")
		return
	}
	f1, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		return
	}
	if len(args) < 2 {
		return
	}
	f2, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return
	}
	two = true
	return
}

func parseInts(args []string) (i1, i2 int, two bool, err error) {
	if len(args) < 1 {
		err = errors.New("No args")
		return
	}
	i1, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if len(args) < 2 {
		return
	}
	i2, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	two = true
	return
}
