package main

import (
	"context"
	"fmt"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/diakovliev/oak/v4"
	"github.com/diakovliev/oak/v4/alg/floatgeom"
	"github.com/diakovliev/oak/v4/audio"
	"github.com/diakovliev/oak/v4/audio/pcm"
	"github.com/diakovliev/oak/v4/audio/synth"
	"github.com/diakovliev/oak/v4/dlog"
	"github.com/diakovliev/oak/v4/entities"
	"github.com/diakovliev/oak/v4/event"
	"github.com/diakovliev/oak/v4/key"
	"github.com/diakovliev/oak/v4/mouse"
	"github.com/diakovliev/oak/v4/render"
	"github.com/diakovliev/oak/v4/scene"
)

const (
	whiteKeyWidth  = 26
	whiteKeyHeight = 200
	blackKeyWidth  = 13
	blackKeyHeight = 140

	whiteBlackOverlap = 5

	labelWhiteKey = 0
	labelBlackKey = 1
)

type keyColor int

const keyColorWhite keyColor = 0
const keyColorBlack keyColor = 1

func (kc keyColor) Width() float64 {
	if kc == keyColorBlack {
		return blackKeyWidth
	}
	return whiteKeyWidth
}

func (kc keyColor) Height() float64 {
	if kc == keyColorBlack {
		return blackKeyHeight
	}
	return whiteKeyHeight
}

func (kc keyColor) Color() color.RGBA {
	if kc == keyColorBlack {
		return color.RGBA{60, 60, 60, 255}
	}
	return color.RGBA{255, 255, 255, 255}
}

func newKey(ctx *scene.Context, note synth.Pitch, c keyColor, k key.Code) *entities.Entity {
	w := c.Width()
	h := c.Height()
	clr := c.Color()
	downClr := clr
	downClr.R -= 60
	downClr.B -= 60
	downClr.G -= 60
	sw := render.NewSwitch("up", map[string]render.Modifiable{
		"up": render.NewCompositeM(
			render.NewColorBox(int(w), int(h), clr),
			render.NewLine(0, 0, 0, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(0, h, w, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, h, w, 0, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, 0, 0, 0, color.RGBA{0, 0, 0, 255}),
		).ToSprite(),
		"down": render.NewCompositeM(
			render.NewColorBox(int(w), int(h), downClr),
			render.NewLine(0, 0, 0, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(0, h, w, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, h, w, 0, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, 0, 0, 0, color.RGBA{0, 0, 0, 255}),
		).ToSprite(),
	})
	s := entities.New(ctx,
		entities.WithUseMouseTree(true),
		entities.WithDimensions(floatgeom.Point2{w, h}),
		entities.WithRenderable(sw),
	)
	if c == keyColorBlack {
		s.Space.SetZLayer(1)
		s.Space.Label = labelBlackKey
	} else {
		s.Space.SetZLayer(2)
		s.Space.Label = labelWhiteKey
	}
	event.GlobalBind(ctx, key.Down(k), func(ev key.Event) event.Response {
		// TODO: add helper function for this?
		if ev.Modifiers&key.ModShift == key.ModShift {
			return 0
		}
		playPitch(ctx, note)
		sw.Set("down")
		return 0
	})
	event.GlobalBind(ctx, key.Up(k), func(ev key.Event) event.Response {
		if ev.Modifiers&key.ModShift == key.ModShift {
			return 0
		}
		releasePitch(note)
		sw.Set("up")
		return 0
	})
	event.Bind(ctx, mouse.PressOn, s, func(_ *entities.Entity, me *mouse.Event) event.Response {
		playPitch(ctx, note)
		me.StopPropagation = true
		sw.Set("down")
		return 0
	})
	event.Bind(ctx, mouse.Release, s, func(_ *entities.Entity, me *mouse.Event) event.Response {
		releasePitch(note)
		sw.Set("up")
		return 0
	})
	return s
}

type keyDef struct {
	color keyColor
	pitch synth.Pitch
	x     float64
}

var keycharOrder = []key.Code{
	key.Z, key.S, key.X, key.D, key.C,
	key.V, key.G, key.B, key.H, key.N, key.J, key.M,
	key.Comma, key.L, key.FullStop, key.Semicolon, key.Slash,
	key.Q, key.Num2, key.W, key.Num3, key.E, key.Num4, key.R,
	key.T, key.Num6, key.Y, key.Num7, key.U,
	key.I, key.Num9, key.O, key.Num0, key.P, key.HyphenMinus, key.LeftSquareBracket,
}

var playLock sync.Mutex
var cancelFuncs = map[synth.Pitch]func(){}

var makeSynth func(ctx context.Context, pitch synth.Pitch)

func playPitch(ctx *scene.Context, pitch synth.Pitch) {
	playLock.Lock()
	defer playLock.Unlock()
	if cancel, ok := cancelFuncs[pitch]; ok {
		cancel()
	}

	gctx, cancel := context.WithCancel(ctx)
	go func() {
		makeSynth(gctx, pitch)
	}()
	cancelFuncs[pitch] = cancel
}

func releasePitch(pitch synth.Pitch) {
	playLock.Lock()
	defer playLock.Unlock()
	if cancel, ok := cancelFuncs[pitch]; ok {
		cancel()
		delete(cancelFuncs, pitch)
	}
}

type pitchText struct {
	pitch *synth.Pitch
}

func (pt *pitchText) String() string {
	if pt.pitch == nil {
		return ""
	}
	return pt.pitch.String() + " - " + strconv.Itoa(int(*pt.pitch))
}

type f64Text struct {
	f64 *float64
}

func (ft *f64Text) String() string {
	if ft.f64 == nil {
		return ""
	}
	return fmt.Sprint(*ft.f64)
}

func main() {
	err := audio.InitDefault()
	if err != nil {
		fmt.Println("init failed:", err)
		os.Exit(1)
	}

	oak.AddScene("piano", scene.Scene{
		Start: func(ctx *scene.Context) {
			var src = new(synth.Source)
			*src = synth.Int16
			src.Format = pcm.Format{
				SampleRate: 80000,
				Channels:   2,
				Bits:       32,
			}
			pt := &pitchText{}
			ft := &f64Text{}
			playWithMonitor := func(gctx context.Context, r pcm.Reader) {
				speaker, err := audio.NewWriter(r.PCMFormat())
				if err != nil {
					fmt.Println("new writer failed:", err)
					return
				}
				monitor := newPCMMonitor(ctx, speaker)
				monitor.SetPos(0, 0)
				render.Draw(monitor)

				pitchDetector := synth.NewPitchDetector(r)
				pt.pitch = &pitchDetector.DetectedPitches[0]
				ft.f64 = &pitchDetector.DetectedRawPitches[0]

				audio.Play(gctx, pitchDetector, func(po *audio.PlayOptions) {
					po.Destination = monitor
				})
				speaker.Close()
				monitor.Undraw()
			}
			makeSynth = func(gctx context.Context, pitch synth.Pitch) {
				toPlay := audio.LoopReader(src.Sin(synth.AtPitch(pitch)))
				fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
				playWithMonitor(gctx, fadeIn)
			}
			render.Draw(render.NewStringerText(pt, 10, 10))
			render.Draw(render.NewStringerText(ft, 10, 20))

			pitch := synth.C3
			kc := keyColorWhite
			x := 20.0
			y := 200.0
			i := 0
			for i < len(keycharOrder) && x+kc.Width() < float64(ctx.Window.Bounds().X()-10) {
				ky := newKey(ctx, pitch, kc, keycharOrder[i])
				ky.SetPos(floatgeom.Point2{x, y})
				layer := 0
				if kc == keyColorBlack {
					layer = 1
				}
				render.Draw(ky.Renderable, layer)
				x += kc.Width()
				pitch = pitch.Up(synth.HalfStep)
				if pitch.IsAccidental() {
					x -= whiteBlackOverlap
					kc = keyColorBlack
				} else if kc != keyColorWhite {
					x -= whiteBlackOverlap
					kc = keyColorWhite
				}
				i++
			}
			// Consider: Adding volume control
			codeKinds := map[key.Code]func(ctx context.Context, pitch synth.Pitch){
				key.S: func(gctx context.Context, pitch synth.Pitch) {
					toPlay := audio.LoopReader(src.Sin(synth.AtPitch(pitch)))
					fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
					playWithMonitor(gctx, fadeIn)
				},
				key.W: func(gctx context.Context, pitch synth.Pitch) {
					toPlay := audio.LoopReader(src.Saw(synth.AtPitch(pitch)))
					fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
					playWithMonitor(gctx, fadeIn)
				},
				key.Q: func(gctx context.Context, pitch synth.Pitch) {
					// demonstrate adding waveforms to play in unison
					unison := 4
					for i := 0; i < unison; i++ {
						go playWithMonitor(gctx, audio.FadeIn(100*time.Millisecond, audio.LoopReader(src.Saw(synth.AtPitch(pitch)))))
						go playWithMonitor(gctx, audio.FadeIn(100*time.Millisecond, audio.LoopReader(src.Saw(synth.AtPitch(pitch), synth.Detune(.04)))))
						go playWithMonitor(gctx, audio.FadeIn(100*time.Millisecond, audio.LoopReader(src.Saw(synth.AtPitch(pitch), synth.Detune(-.05)))))
					}
					playWithMonitor(gctx, audio.FadeIn(100*time.Millisecond, audio.LoopReader(src.Saw(synth.AtPitch(pitch)))))
				},
				key.T: func(gctx context.Context, pitch synth.Pitch) {
					toPlay := audio.LoopReader(src.Triangle(synth.AtPitch(pitch)))
					fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
					playWithMonitor(gctx, fadeIn)
				},
				key.P: func(gctx context.Context, pitch synth.Pitch) {
					toPlay := audio.LoopReader(src.Pulse(2)(synth.AtPitch(pitch)))
					fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
					playWithMonitor(gctx, fadeIn)
				},
				key.N: func(gctx context.Context, pitch synth.Pitch) {
					toPlay := audio.LoopReader(src.Noise(synth.AtPitch(pitch)))
					fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
					playWithMonitor(gctx, fadeIn)
				},
				key.X: func(gctx context.Context, pitch synth.Pitch) {
					// demonstrate combining multiple wave forms in place
					toPlay := src.MultiWave([]synth.Waveform{
						synth.Source.SinWave,
						synth.Source.TriangleWave,
						synth.PulseWave(2),
					}, synth.AtPitch(pitch))
					fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
					playWithMonitor(gctx, fadeIn)

				},
			}
			for kc, synfn := range codeKinds {
				synfn := synfn
				kc := kc
				event.GlobalBind(ctx, key.Down(kc), func(ev key.Event) event.Response {
					if ev.Modifiers&key.ModShift == key.ModShift {
						makeSynth = synfn
					}
					return 0
				})
			}

			help1 := render.NewText("Shift+([S]in/[T]ri/[P]ulse/sa[W]) to change wave style", 10, 500)
			help2 := render.NewText("Keyboard / mouse to play", 10, 520)
			render.Draw(help1)
			render.Draw(help2)

			event.GlobalBind(ctx, mouse.ScrollDown, func(_ *mouse.Event) event.Response {
				mag := globalMagnification - 0.05
				if mag < 1 {
					mag = 1
				}
				globalMagnification = mag
				return 0
			})
			event.GlobalBind(ctx, mouse.ScrollUp, func(_ *mouse.Event) event.Response {
				globalMagnification += 0.05
				return 0
			})
			event.GlobalBind(ctx, key.Down(key.Keypad0), func(_ key.Event) event.Response {
				// TODO: synth all sound like pulse waves at 8 bit
				src.Bits = 8
				return 0
			})
			event.GlobalBind(ctx, key.Down(key.Keypad1), func(_ key.Event) event.Response {
				src.Bits = 16
				return 0
			})
			event.GlobalBind(ctx, key.Down(key.Keypad2), func(_ key.Event) event.Response {
				src.Bits = 32
				return 0
			})
		},
	})
	oak.Init("piano", func(c oak.Config) (oak.Config, error) {
		c.Screen.Height = 600
		c.Title = "Piano Example"
		c.Debug.Level = dlog.INFO.String()
		return c, nil
	})
}

type pcmMonitor struct {
	event.CallerID
	render.LayeredPoint
	pcm.Writer
	pcm.Format
	written []byte
	at      int
}

var globalMagnification float64 = 1

func newPCMMonitor(ctx *scene.Context, w pcm.Writer) *pcmMonitor {
	fmt := w.PCMFormat()
	pm := &pcmMonitor{
		Writer:       w,
		Format:       w.PCMFormat(),
		LayeredPoint: render.NewLayeredPoint(0, 0, 0),
		written:      make([]byte, int(float64(fmt.BytesPerSecond())*audio.WriterBufferLengthInSeconds)),
	}
	return pm
}

func (pm *pcmMonitor) CID() event.CallerID {
	return pm.CallerID
}

func (pm *pcmMonitor) PCMFormat() pcm.Format {
	return pm.Format
}

func (pm *pcmMonitor) WritePCM(b []byte) (n int, err error) {
	copy(pm.written[pm.at:], b)
	if len(b) > len(pm.written[pm.at:]) {
		copy(pm.written[0:], b[len(pm.written[pm.at:]):])
	}
	pm.at += len(b)
	pm.at %= len(pm.written)
	return pm.Writer.WritePCM(b)
}

func (pm *pcmMonitor) Draw(buf draw.Image, xOff, yOff float64) {
	const width = 640
	const height = 200.0
	xJump := len(pm.written) / width
	xJump = int(float64(xJump) / globalMagnification)
	c := color.RGBA{255, 255, 255, 255}
	for x := 0.0; x < width; x++ {
		wIndex := int(x) * xJump

		var val int16
		switch pm.Format.Bits {
		case 8:
			val8 := pm.written[wIndex]
			val = int16(val8) << 8
		case 16:
			wIndex -= wIndex % 2
			val = int16(pm.written[wIndex+1])<<8 +
				int16(pm.written[wIndex])
		case 32:
			wIndex = wIndex - wIndex%4
			val32 := int32(pm.written[wIndex+3])<<24 +
				int32(pm.written[wIndex+2])<<16 +
				int32(pm.written[wIndex+1])<<8 +
				int32(pm.written[wIndex])
			val = int16(val32 / int32(math.Pow(2, 16)))
		}

		// -32768 -> 200
		// 0 -> 100
		// 32768 -> 0
		var y float64
		if val < 0 {
			y = height/2 + float64(val)*float64(height/2/-32768.0)
		} else {
			y = height/2 + -(float64(val) * float64(height/2/32768.0))
		}
		buf.Set(int(x+xOff+pm.X()), int(y+yOff+pm.Y()), c)
	}
}
